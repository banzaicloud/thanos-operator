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

package receiver

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

func (r receiverInstance) statefulset() (runtime.Object, reconciler.DesiredState, error) {
	statefulSet := &appsv1.StatefulSet{
		ObjectMeta: r.receiverGroup.MetaOverrides.Merge(r.getMeta(r.receiverGroup.Name)),
		Spec: appsv1.StatefulSetSpec{
			Replicas: utils.IntPointer(r.receiverGroup.Replicas),
			Selector: &metav1.LabelSelector{
				MatchLabels: r.getLabels(),
			},
			ServiceName: r.getName(r.receiverGroup.Name),
			Template: corev1.PodTemplateSpec{
				ObjectMeta: r.getMeta(),
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "receive",
							Image: fmt.Sprintf("%s:%s", v1alpha1.ThanosImageRepository, v1alpha1.ThanosImageTag),
							Args: []string{
								"receive",
								fmt.Sprintf("--objstore.config-file=/etc/config/%s", r.receiverGroup.Config.MountFrom.SecretKeyRef.Key),
								fmt.Sprintf("--receive.local-endpoint=$(NAME).%s:10907", r.getName(r.receiverGroup.Name)),
								"--receive.hashrings-file=/etc/hashring/hashring.json",
								fmt.Sprintf("--receive.replication-factor=%d", 1),
								"--label=receive_replica=\"$(NAME)\"",
								"--log.level=debug",
							},
							Env: []corev1.EnvVar{
								{
									Name: "NAME",
									ValueFrom: &corev1.EnvVarSource{
										FieldRef: &corev1.ObjectFieldSelector{
											FieldPath: "metadata.name",
										},
									},
								},
							},
							Ports: []corev1.ContainerPort{
								{
									Name:          "http",
									ContainerPort: resources.GetPort(r.receiverGroup.HTTPAddress),
									Protocol:      corev1.ProtocolTCP,
								},
								{
									Name:          "grpc",
									ContainerPort: resources.GetPort(r.receiverGroup.GRPCAddress),
									Protocol:      corev1.ProtocolTCP,
								},
								{
									Name:          "remote-write",
									ContainerPort: resources.GetPort(r.receiverGroup.RemoteWriteAddress),
									Protocol:      corev1.ProtocolTCP,
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "objectstore-secret",
									ReadOnly:  true,
									MountPath: "/etc/config/",
								},
								{
									Name:      "hashring-config",
									ReadOnly:  true,
									MountPath: "/etc/hashring/",
								},
							},
							LivenessProbe:   resources.GetProbe(resources.GetPort(r.receiverGroup.HTTPAddress), resources.HealthCheckPath),
							ReadinessProbe:  resources.GetProbe(resources.GetPort(r.receiverGroup.HTTPAddress), resources.ReadyCheckPath),
							ImagePullPolicy: corev1.PullIfNotPresent,
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "objectstore-secret",
							VolumeSource: corev1.VolumeSource{
								Secret: &corev1.SecretVolumeSource{
									SecretName: r.receiverGroup.Config.MountFrom.SecretKeyRef.Name,
								},
							},
						},
						{
							Name: "hashring-config",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: r.getName("hashring-config"),
									},
								},
							},
						},
					},
				},
			},
		},
	}

	statefulSet.Spec.Template.Spec.Containers[0].Args = r.setArgs(statefulSet.Spec.Template.Spec.Containers[0].Args)

	if r.receiverGroup.DataVolume != nil {
		statefulSet.Spec.Template.Spec.Containers[0].VolumeMounts = append(statefulSet.Spec.Template.Spec.Containers[0].VolumeMounts, corev1.VolumeMount{
			Name:      "data-volume",
			MountPath: r.receiverGroup.TSDBPath,
		})

		if r.receiverGroup.DataVolume.PersistentVolumeClaim != nil {
			statefulSet.Spec.VolumeClaimTemplates = []corev1.PersistentVolumeClaim{
				{
					ObjectMeta: r.getVolumeMeta("data-volume"),
					Spec:       r.receiverGroup.DataVolume.PersistentVolumeClaim.PersistentVolumeClaimSpec,
					Status: corev1.PersistentVolumeClaimStatus{
						Phase: corev1.ClaimPending,
					},
				},
			}
		} else {
			volume, err := r.receiverGroup.DataVolume.GetVolume("data-volume")
			if err != nil {
				return statefulSet, reconciler.StateAbsent, err
			}

			statefulSet.Spec.Template.Spec.Volumes = append(statefulSet.Spec.Template.Spec.Volumes, volume)
		}
	}

	if r.receiverGroup.StatefulSetOverrides != nil {
		if err := merge.Merge(statefulSet, r.receiverGroup.StatefulSetOverrides); err != nil {
			return statefulSet, reconciler.StatePresent, errors.WrapIf(err, "unable to merge overrides to base object")
		}
	}

	return statefulSet, reconciler.StatePresent, nil
}
