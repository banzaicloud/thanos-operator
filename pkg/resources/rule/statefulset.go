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

package rule

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

func (r *ruleInstance) statefulset() (runtime.Object, reconciler.DesiredState, error) {
	if r.Thanos.Spec.Rule != nil {
		rule := r.Thanos.Spec.Rule
		statefulset := &appsv1.StatefulSet{
			ObjectMeta: rule.MetaOverrides.Merge(r.getMeta()),
			Spec: appsv1.StatefulSetSpec{
				Replicas: utils.IntPointer(1),
				Selector: &metav1.LabelSelector{
					MatchLabels: r.getLabels(),
				},
				Template: corev1.PodTemplateSpec{
					ObjectMeta: r.getMeta(),
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{
							{
								Name:  "rule",
								Image: fmt.Sprintf("%s:%s", v1alpha1.ThanosImageRepository, v1alpha1.ThanosImageTag),
								Args: []string{
									"rule",
									fmt.Sprintf("--objstore.config-file=/etc/config/%s", r.StoreEndpoint.Spec.Config.MountFrom.SecretKeyRef.Key),
								},
								WorkingDir: "",
								Ports: []corev1.ContainerPort{
									{
										Name:          "http",
										ContainerPort: resources.GetPort(rule.HttpAddress),
										Protocol:      corev1.ProtocolTCP,
									},
									{
										Name:          "grpc",
										ContainerPort: resources.GetPort(rule.GRPCAddress),
										Protocol:      corev1.ProtocolTCP,
									},
								},
								VolumeMounts: []corev1.VolumeMount{
									{
										Name:      "objectstore-secret",
										ReadOnly:  true,
										MountPath: "/etc/config/",
									},
								},
								LivenessProbe:   resources.GetProbe(resources.GetPort(rule.HttpAddress), resources.HealthCheckPath),
								ReadinessProbe:  resources.GetProbe(resources.GetPort(rule.HttpAddress), resources.ReadyCheckPath),
								ImagePullPolicy: corev1.PullIfNotPresent,
							},
						},
						Volumes: []corev1.Volume{
							{
								Name: "objectstore-secret",
								VolumeSource: corev1.VolumeSource{
									Secret: &corev1.SecretVolumeSource{
										SecretName: r.StoreEndpoint.Spec.Config.MountFrom.SecretKeyRef.Name,
									},
								},
							},
						},
					},
				},
			},
		}

		statefulset.Spec.Template.Spec.Containers[0].Args = r.setArgs(statefulset.Spec.Template.Spec.Containers[0].Args)

		if r.Thanos.Spec.Rule.DataVolume != nil {
			if r.Thanos.Spec.Rule.DataVolume.PersistentVolumeClaim != nil {
				statefulset.Spec.Template.Spec.Containers[0].VolumeMounts = append(statefulset.Spec.Template.Spec.Containers[0].VolumeMounts, corev1.VolumeMount{
					Name:      "data-volume",
					MountPath: r.Thanos.Spec.Rule.DataDir,
				})
				statefulset.Spec.VolumeClaimTemplates = []corev1.PersistentVolumeClaim{
					{
						ObjectMeta: r.getVolumeMeta("data-volume"),
						Spec:       r.Thanos.Spec.Rule.DataVolume.PersistentVolumeClaim.PersistentVolumeClaimSpec,
						Status: corev1.PersistentVolumeClaimStatus{
							Phase: corev1.ClaimPending,
						},
					},
				}
			} else {
				statefulset.Spec.Template.Spec.Containers[0].VolumeMounts = append(statefulset.Spec.Template.Spec.Containers[0].VolumeMounts, corev1.VolumeMount{
					Name:      "data-volume",
					MountPath: r.Thanos.Spec.Rule.DataDir,
				})

				volume, err := r.Thanos.Spec.Rule.DataVolume.GetVolume("data-volume")
				if err != nil {
					return statefulset, reconciler.StateAbsent, err
				}

				statefulset.Spec.Template.Spec.Volumes = append(statefulset.Spec.Template.Spec.Volumes, volume)
			}
		}

		if rule.StatefulsetOverrides != nil {
			if err := merge.Merge(statefulset, rule.StatefulsetOverrides); err != nil {
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
