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

package store

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

func (s *storeInstance) deployment() (runtime.Object, reconciler.DesiredState, error) {
	storeGateway := s.Store.Thanos.Spec.StoreGateway
	if storeGateway != nil {
		if s.StoreEndpoint.Spec.Config.MountFrom == nil {
			return nil, nil, fmt.Errorf("missing config for StorageGateway %q", s.StoreEndpoint.Name)
		}

		deployment := &appsv1.Deployment{
			ObjectMeta: storeGateway.MetaOverrides.Merge(s.getMeta()),
			Spec: appsv1.DeploymentSpec{
				Replicas: utils.IntPointer(1),
				Selector: &metav1.LabelSelector{
					MatchLabels: s.getLabels(),
				},
				Template: corev1.PodTemplateSpec{
					ObjectMeta: s.getMeta(),
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{
							{
								Name:  v1alpha1.StoreName,
								Image: fmt.Sprintf("%s:%s", v1alpha1.ThanosImageRepository, v1alpha1.ThanosImageTag),
								Args: []string{
									"store",
									fmt.Sprintf("--objstore.config-file=/etc/config/%s", s.StoreEndpoint.Spec.Config.MountFrom.SecretKeyRef.Key),
								},
								Ports: []corev1.ContainerPort{
									{
										Name:          "http",
										ContainerPort: resources.GetPort(storeGateway.HttpAddress),
										Protocol:      corev1.ProtocolTCP,
									},
									{
										Name:          "grpc",
										ContainerPort: resources.GetPort(storeGateway.GRPCAddress),
										Protocol:      corev1.ProtocolTCP,
									},
								},
								VolumeMounts:    s.getVolumeMounts(),
								ImagePullPolicy: corev1.PullIfNotPresent,
								LivenessProbe:   resources.GetProbe(resources.GetPort(storeGateway.HttpAddress), resources.HealthCheckPath),
								ReadinessProbe:  resources.GetProbe(resources.GetPort(storeGateway.HttpAddress), resources.ReadyCheckPath),
							},
						},
						Volumes: s.getVolumes(),
					},
				},
			},
		}

		// Set up args
		deployment.Spec.Template.Spec.Containers[0].Args = s.setArgs(deployment.Spec.Template.Spec.Containers[0].Args)

		if storeGateway.DeploymentOverrides != nil {
			if err := merge.Merge(deployment, storeGateway.DeploymentOverrides); err != nil {
				return deployment, reconciler.StatePresent, errors.WrapIf(err, "unable to merge overrides to base object")
			}
		}

		return deployment, reconciler.StatePresent, nil
	}
	delete := &appsv1.Deployment{
		ObjectMeta: s.getMeta(),
	}
	return delete, reconciler.StateAbsent, nil
}

func (s *storeInstance) getVolumeMounts() []corev1.VolumeMount {
	volumeMounts := []corev1.VolumeMount{
		{
			Name:      "objectstore-secret",
			ReadOnly:  true,
			MountPath: "/etc/config/",
		},
	}
	if s.Thanos.Spec.StoreGateway.GRPCServerCertificate != "" {
		volumeMounts = append(volumeMounts, corev1.VolumeMount{
			Name:      "server-certificate",
			ReadOnly:  true,
			MountPath: serverCertMountPath,
		})
	}
	return volumeMounts
}

func (s *storeInstance) getVolumes() []corev1.Volume {
	volumes := []corev1.Volume{
		{
			Name: "objectstore-secret",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName: s.StoreEndpoint.Spec.Config.MountFrom.SecretKeyRef.Name,
				},
			},
		},
	}
	if s.Thanos.Spec.StoreGateway.GRPCServerCertificate != "" {
		volumes = append(volumes, corev1.Volume{
			Name: "server-certificate",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName: s.Thanos.Spec.StoreGateway.GRPCServerCertificate,
				},
			},
		})
	}
	return volumes
}
