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
	thanosconfig "github.com/banzaicloud/thanos-operator/controllers/config"
	"github.com/banzaicloud/thanos-operator/pkg/resources"
	"github.com/banzaicloud/thanos-operator/pkg/sdk/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func (r *receiverInstance) statefulset() (runtime.Object, reconciler.DesiredState, error) {
	cfg := thanosconfig.GetControllerConfig()
	if r.receiverGroup != nil {
		receiveGroup := r.receiverGroup

		statefulset := &appsv1.StatefulSet{
			ObjectMeta: receiveGroup.MetaOverrides.Merge(r.getMeta(r.receiverGroup.Name)),
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
								Image: fmt.Sprintf("%s:%s", cfg.GetConfigString("ThanosImage", v1alpha1.ThanosImageRepository), cfg.GetConfigString("ThanosImageTag", v1alpha1.ThanosImageTag)),
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
								WorkingDir: "",
								Ports: []corev1.ContainerPort{
									{
										Name:          "http",
										ContainerPort: resources.GetPort(receiveGroup.HTTPAddress),
										Protocol:      corev1.ProtocolTCP,
									},
									{
										Name:          "grpc",
										ContainerPort: resources.GetPort(receiveGroup.GRPCAddress),
										Protocol:      corev1.ProtocolTCP,
									},
									{
										Name:          "remote-write",
										ContainerPort: resources.GetPort(receiveGroup.RemoteWriteAddress),
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
								LivenessProbe:   r.GetCheck(resources.GetPort(receiveGroup.HTTPAddress), resources.HealthCheckPath),
								ReadinessProbe:  r.GetCheck(resources.GetPort(receiveGroup.HTTPAddress), resources.ReadyCheckPath),
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

		statefulset.Spec.Template.Spec.Containers[0].Args = r.setArgs(statefulset.Spec.Template.Spec.Containers[0].Args)

		if receiveGroup.DataVolume != nil {
			if receiveGroup.DataVolume.PersistentVolumeClaim != nil {
				statefulset.Spec.Template.Spec.Containers[0].VolumeMounts = append(statefulset.Spec.Template.Spec.Containers[0].VolumeMounts, corev1.VolumeMount{
					Name:      "data-volume",
					MountPath: receiveGroup.TSDBPath,
				})
				statefulset.Spec.VolumeClaimTemplates = []corev1.PersistentVolumeClaim{
					{
						ObjectMeta: r.getVolumeMeta("data-volume"),
						Spec:       receiveGroup.DataVolume.PersistentVolumeClaim.PersistentVolumeClaimSpec,
						Status: corev1.PersistentVolumeClaimStatus{
							Phase: corev1.ClaimPending,
						},
					},
				}
			} else {
				statefulset.Spec.Template.Spec.Containers[0].VolumeMounts = append(statefulset.Spec.Template.Spec.Containers[0].VolumeMounts, corev1.VolumeMount{
					Name:      "data-volume",
					MountPath: receiveGroup.TSDBPath,
				})

				volume, err := receiveGroup.DataVolume.GetVolume("data-volume")
				if err != nil {
					return statefulset, reconciler.StateAbsent, err
				}

				statefulset.Spec.Template.Spec.Volumes = append(statefulset.Spec.Template.Spec.Volumes, volume)
			}
		}

		if r.receiverGroup.StatefulSetOverrides != nil {
			if err := merge.Merge(statefulset, r.receiverGroup.StatefulSetOverrides); err != nil {
				return statefulset, reconciler.StatePresent, errors.WrapIf(err, "unable to merge overrides to base object")
			}
		}

		return statefulset, reconciler.StatePresent, nil
	}

	delete := &appsv1.StatefulSet{
		ObjectMeta: r.getMeta(),
	}
	return delete, reconciler.StateAbsent, nil
}
