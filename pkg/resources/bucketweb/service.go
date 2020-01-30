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
	"github.com/banzaicloud/thanos-operator/pkg/resources"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func (b *BucketWeb) service() (runtime.Object, reconciler.DesiredState, error) {
	if b.ObjectStore.Spec.BucketWeb != nil {
		bucketWeb := b.ObjectStore.Spec.BucketWeb.DeepCopy()
		return &corev1.Service{
			ObjectMeta: b.getMeta(),
			Spec: corev1.ServiceSpec{
				Ports: []corev1.ServicePort{
					{
						Protocol: corev1.ProtocolTCP,
						Name:     "http",
						Port:     resources.GetPort(bucketWeb.HTTPAddress),
						TargetPort: intstr.IntOrString{
							Type:   intstr.Int,
							IntVal: resources.GetPort(bucketWeb.HTTPAddress),
						},
					},
				},
				Selector: b.getLabels(),
			},
		}, reconciler.StatePresent, nil
	}

	return &corev1.Service{
		ObjectMeta: b.getMeta(),
	}, reconciler.StateAbsent, nil
}
