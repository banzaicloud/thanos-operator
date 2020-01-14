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
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func (b *BucetWeb) deployment() (runtime.Object, reconciler.DesiredState, error) {
	const app = "bucketweb"
	name := app + "-deployment"

	if b.objectStore.Spec.BucketWeb.Enabled {
		bucketWeb := b.objectStore.Spec.BucketWeb.DeepCopy()

		var deployment = &appsv1.Deployment{
			ObjectMeta: b.objectMeta(name, &bucketWeb.BaseObject),
			Spec: appsv1.DeploymentSpec{
				Replicas: utils.IntPointer(1),
				Selector: &metav1.LabelSelector{
					MatchLabels: map[string]string{"app": app},
				},
				Template: corev1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Name:        app,
						Labels:      map[string]string{"app": app},
						Annotations: b.objectStore.Annotations,
					},
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{
							{
								Name:  app,
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
										ContainerPort: GetPort(bucketWeb.HTTPAddress),
										Protocol:      corev1.ProtocolTCP,
									},
								},
								Resources:       bucketWeb.Resources,
								ImagePullPolicy: corev1.PullPolicy(bucketWeb.Image.PullPolicy),
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
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: b.namespace,
		},
	}, reconciler.StateAbsent, nil
}
