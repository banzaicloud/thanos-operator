// Copyright © 2020 Banzai Cloud
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package bucketweb

import (
	"fmt"
	"math"
	"strconv"

	"github.com/banzaicloud/operator-tools/pkg/reconciler"
	"github.com/banzaicloud/operator-tools/pkg/utils"
	"github.com/banzaicloud/thanos-operator/pkg/resources"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func (b *BucketWeb) deployment() (runtime.Object, reconciler.DesiredState, error) {
	if b.ObjectStore.Spec.BucketWeb != nil {
		bucketWeb := b.ObjectStore.Spec.BucketWeb.DeepCopy()
		var deployment = &appsv1.Deployment{
			ObjectMeta: b.getMeta(),
			Spec: appsv1.DeploymentSpec{
				Replicas: utils.IntPointer(1),
				Selector: &metav1.LabelSelector{
					MatchLabels: b.getLabels(),
				},
				Template: corev1.PodTemplateSpec{
					ObjectMeta: b.getMeta(),
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{
							{
								Name:  Name,
								Image: fmt.Sprintf("%s:%s", bucketWeb.Image.Repository, bucketWeb.Image.Tag),
								Args: []string{
									"bucket",
									"web",
									"--log.level=info",
									"--http-address=" + bucketWeb.HTTPAddress,
									// TODO: get secret file path from secret mount
									"--objstore.config-file=/etc/config/object-store.yaml",
									"--refresh=" + strconv.Itoa(int(math.Floor(bucketWeb.Refresh.Duration.Seconds()))) + "s",
									"--timeout=" + strconv.Itoa(int(math.Floor(bucketWeb.Timeout.Duration.Seconds()))) + "s",
								},
								Ports: []corev1.ContainerPort{
									{
										Name:          "http",
										ContainerPort: resources.GetPort(bucketWeb.HTTPAddress),
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
								Resources:       bucketWeb.Resources,
								ImagePullPolicy: corev1.PullPolicy(bucketWeb.Image.PullPolicy),
							},
						},
						Volumes: []corev1.Volume{
							{
								Name: "objectstore-secret",
								VolumeSource: corev1.VolumeSource{
									Secret: &corev1.SecretVolumeSource{
										SecretName: b.ObjectStore.Spec.Config.MountFrom.SecretKeyRef.Name,
									},
								},
							},
						},
					},
				},
			},
		}

		if bucketWeb.Label != "" {
			deployment.Spec.Template.Spec.Containers[0].Args = append(deployment.Spec.Template.Spec.Containers[0].Args, "--label="+bucketWeb.Label)
		}

		return deployment, reconciler.StatePresent, nil
	}

	return &appsv1.Deployment{
		ObjectMeta: b.getMeta(),
	}, reconciler.StateAbsent, nil
}
