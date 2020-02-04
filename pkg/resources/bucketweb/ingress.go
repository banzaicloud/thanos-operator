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
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func (b *BucketWeb) ingress() (runtime.Object, reconciler.DesiredState, error) {
	if b.ObjectStore.Spec.BucketWeb != nil {
		bucketWeb := b.ObjectStore.Spec.BucketWeb.DeepCopy()
		var ingress = &extensionsv1beta1.Ingress{
			ObjectMeta: b.getMeta(),
			Spec: extensionsv1beta1.IngressSpec{
				//Certificate: []extensionsv1beta1.IngressTLS{
				//	{
				//		Hosts: []string{"/"},
				//	},
				//},
				Rules: []extensionsv1beta1.IngressRule{
					{
						Host: "/",
						IngressRuleValue: extensionsv1beta1.IngressRuleValue{
							HTTP: &extensionsv1beta1.HTTPIngressRuleValue{
								Paths: []extensionsv1beta1.HTTPIngressPath{
									{
										Path: "/",
										Backend: extensionsv1beta1.IngressBackend{
											ServiceName: b.getName(),
											ServicePort: intstr.IntOrString{IntVal: resources.GetPort(bucketWeb.HTTPAddress)},
										},
									},
								},
							},
						},
					},
				},
			},
		}

		return ingress, reconciler.StateAbsent, nil
	}

	return &extensionsv1beta1.Ingress{
		ObjectMeta: b.getMeta(),
	}, reconciler.StateAbsent, nil
}
