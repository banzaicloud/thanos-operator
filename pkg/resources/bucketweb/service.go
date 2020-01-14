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
	"github.com/banzaicloud/operator-tools/pkg/reconciler"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func (b *BucetWeb) service() (runtime.Object, reconciler.DesiredState, error) {
	const app = "bucketweb"
	name := app + "-service"

	if b.objectStore.Spec.BucketWeb.Enabled {
		bucketWeb := b.objectStore.Spec.BucketWeb.DeepCopy()

		return &corev1.Service{
			ObjectMeta: b.objectMeta(name, &bucketWeb.BaseObject),
			Spec: corev1.ServiceSpec{
				Ports: []corev1.ServicePort{
					{
						Protocol:   corev1.ProtocolTCP,
						Name:       "http",
						Port:       GetPort(bucketWeb.HTTPAddress),
						TargetPort: intstr.IntOrString{IntVal: GetPort(bucketWeb.HTTPAddress)},
					},
				},
				Selector: b.labels(),
			},
		}, reconciler.StatePresent, nil
	}

	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: b.namespace,
		},
	}, reconciler.StateAbsent, nil
}
