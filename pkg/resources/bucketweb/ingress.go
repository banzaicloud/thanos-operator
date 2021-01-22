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

package bucketweb

import (
	"github.com/banzaicloud/operator-tools/pkg/reconciler"
	"github.com/banzaicloud/thanos-operator/pkg/resources"
	netv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func (b *BucketWeb) ingress() (runtime.Object, reconciler.DesiredState, error) {
	if b.ObjectStore.Spec.BucketWeb != nil {
		bucketWeb := b.ObjectStore.Spec.BucketWeb.DeepCopy()
		pathType := netv1.PathTypeImplementationSpecific
		ingress := &netv1.Ingress{
			ObjectMeta: bucketWeb.MetaOverrides.Merge(b.getMeta()),
			Spec: netv1.IngressSpec{
				//Certificate: []networkingv1.IngressTLS{
				//	{
				//		Hosts: []string{"/"},
				//	},
				//},
				Rules: []netv1.IngressRule{
					{
						Host: "/",
						IngressRuleValue: netv1.IngressRuleValue{
							HTTP: &netv1.HTTPIngressRuleValue{
								Paths: []netv1.HTTPIngressPath{
									{
										Path:     "/",
										PathType: &pathType,
										Backend: netv1.IngressBackend{
											Service: &netv1.IngressServiceBackend{
												Name: b.getName(),
												Port: netv1.ServiceBackendPort{
													Number: resources.GetPort(bucketWeb.HTTPAddress),
												},
											},
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

	return &netv1.Ingress{
		ObjectMeta: b.getMeta(),
	}, reconciler.StateAbsent, nil
}
