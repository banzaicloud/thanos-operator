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
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func (b *BucetWeb) ingress() (runtime.Object, reconciler.DesiredState, error) {
	const app = "bucketweb"
	name := app + "-ingress"

	if b.objectStore.Spec.BucketWeb.Enabled {
		bucketWeb := b.objectStore.Spec.BucketWeb.DeepCopy()

		var ingress = &extensionsv1beta1.Ingress{
			ObjectMeta: b.objectMeta(name, &bucketWeb.BaseObject),
			Spec: extensionsv1beta1.IngressSpec{
				//TLS: []extensionsv1beta1.IngressTLS{
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
											ServiceName: name,
											ServicePort: intstr.IntOrString{IntVal: GetPort(bucketWeb.HTTPAddress)},
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
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: b.namespace,
		},
	}, reconciler.StateAbsent, nil
}
