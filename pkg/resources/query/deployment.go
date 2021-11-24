// Copyright 2020 Banzai Cloud
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package query

import (
	"fmt"

	"emperror.dev/errors"
	"github.com/banzaicloud/operator-tools/pkg/merge"
	"github.com/banzaicloud/operator-tools/pkg/reconciler"
	"github.com/banzaicloud/operator-tools/pkg/utils"
	"github.com/banzaicloud/thanos-operator/pkg/resources"
	"github.com/banzaicloud/thanos-operator/pkg/sdk/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func (q *Query) deployment() (runtime.Object, reconciler.DesiredState, error) {
	if q.Thanos.Spec.Query != nil {
		query := q.Thanos.Spec.Query.DeepCopy()

		deployment := &appsv1.Deployment{
			ObjectMeta: query.MetaOverrides.Merge(q.getMeta(q.getName())),
			Spec: appsv1.DeploymentSpec{
				Replicas: utils.IntPointer(1),
				Selector: &metav1.LabelSelector{
					MatchLabels: q.getLabels(),
				},
				Template: corev1.PodTemplateSpec{
					ObjectMeta: q.getMeta(q.getName()),
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{
							{
								Name:  "query",
								Image: fmt.Sprintf("%s:%s", v1alpha1.ThanosImageRepository, v1alpha1.ThanosImageTag),
								Args: []string{
									"query",
								},
								Ports: []corev1.ContainerPort{
									{
										Name:          "http",
										ContainerPort: resources.GetPort(query.HttpAddress),
										Protocol:      corev1.ProtocolTCP,
									},
									{
										Name:          "grpc",
										ContainerPort: resources.GetPort(query.GRPCAddress),
										Protocol:      corev1.ProtocolTCP,
									},
								},
								ImagePullPolicy: corev1.PullIfNotPresent,
								LivenessProbe:   resources.GetProbe(resources.GetPort(query.HttpAddress), resources.HealthCheckPath),
								ReadinessProbe:  resources.GetProbe(resources.GetPort(query.HttpAddress), resources.ReadyCheckPath),
								VolumeMounts:    q.getVolumeMounts(),
							},
						},
						Volumes: q.getVolumes(),
					},
				},
			},
		}

		// Set up args
		deployment.Spec.Template.Spec.Containers[0].Args = q.setArgs(deployment.Spec.Template.Spec.Containers[0].Args)

		if query.DeploymentOverrides != nil {
			err := merge.Merge(deployment, query.DeploymentOverrides)
			if err != nil {
				return deployment, reconciler.StatePresent, errors.WrapIf(err, "unable to merge overrides to base object")
			}
		}

		return deployment, reconciler.StatePresent, nil
	}
	delete := &appsv1.Deployment{
		ObjectMeta: q.getMeta(q.getName()),
	}
	return delete, reconciler.StateAbsent, nil
}

func (q *Query) getVolumeMounts() []corev1.VolumeMount {
	volumeMounts := make([]corev1.VolumeMount, 0)
	if q.Thanos.Spec.Query.GRPCClientCertificate != "" {
		volumeMounts = append(volumeMounts, corev1.VolumeMount{
			Name:      "client-certificate",
			ReadOnly:  true,
			MountPath: clientCertMountPath,
		})
	}
	if q.Thanos.Spec.Query.GRPCClientCA != "" {
		volumeMounts = append(volumeMounts, corev1.VolumeMount{
			Name:      "client-ca",
			ReadOnly:  true,
			MountPath: clientCACertMountPath,
		})
	}
	if q.Thanos.Spec.Query.GRPCServerCertificate != "" {
		volumeMounts = append(volumeMounts, corev1.VolumeMount{
			Name:      "server-certificate",
			ReadOnly:  true,
			MountPath: serverCertMountPath,
		})
	}
	if q.Thanos.Spec.Query.GRPCServerCA != "" {
		volumeMounts = append(volumeMounts, corev1.VolumeMount{
			Name:      "server-ca",
			ReadOnly:  true,
			MountPath: serverCACertMountPath,
		})
	}
	return volumeMounts
}

func (q *Query) getVolumes() []corev1.Volume {
	volumes := make([]corev1.Volume, 0)
	if q.Thanos.Spec.Query.GRPCClientCertificate != "" {
		volumes = append(volumes, corev1.Volume{
			Name: "client-certificate",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName: q.Thanos.Spec.Query.GRPCClientCertificate,
				},
			},
		})
	}
	if q.Thanos.Spec.Query.GRPCClientCA != "" {
		volumes = append(volumes, corev1.Volume{
			Name: "client-ca",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName: q.Thanos.Spec.Query.GRPCClientCA,
				},
			},
		})
	}
	if q.Thanos.Spec.Query.GRPCServerCertificate != "" {
		volumes = append(volumes, corev1.Volume{
			Name: "server-certificate",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName: q.Thanos.Spec.Query.GRPCServerCertificate,
				},
			},
		})
	}
	if q.Thanos.Spec.Query.GRPCServerCA != "" {
		volumes = append(volumes, corev1.Volume{
			Name: "server-ca",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName: q.Thanos.Spec.Query.GRPCServerCA,
				},
			},
		})
	}
	return volumes
}
