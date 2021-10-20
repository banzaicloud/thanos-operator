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
	"emperror.dev/errors"
	"github.com/banzaicloud/operator-tools/pkg/merge"
	"github.com/banzaicloud/operator-tools/pkg/reconciler"
	netv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func (b *BucketWeb) ingress() (runtime.Object, reconciler.DesiredState, error) {
	meta := b.getMeta(b.getName())
	if b.ObjectStore.Spec.BucketWeb != nil &&
		b.ObjectStore.Spec.BucketWeb.HTTPIngress != nil {
		bucketWeb := b.ObjectStore.Spec.BucketWeb
		bucketWebIngress := bucketWeb.HTTPIngress
		pathType := netv1.PathTypeImplementationSpecific
		ingress := &netv1.Ingress{
			ObjectMeta: bucketWeb.MetaOverrides.Merge(meta),
			Spec: netv1.IngressSpec{
				Rules: []netv1.IngressRule{
					{
						Host: bucketWebIngress.Host,
						IngressRuleValue: netv1.IngressRuleValue{
							HTTP: &netv1.HTTPIngressRuleValue{
								Paths: []netv1.HTTPIngressPath{
									{
										Path:     bucketWebIngress.Path,
										PathType: &pathType,
										Backend: netv1.IngressBackend{
											Service: &netv1.IngressServiceBackend{
												Name: b.getName(),
												Port: netv1.ServiceBackendPort{
													Name: "http",
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
		if bucketWebIngress.Certificate != "" {
			ingress.Spec.TLS = []netv1.IngressTLS{
				{
					Hosts:      []string{bucketWebIngress.Host},
					SecretName: bucketWebIngress.Certificate,
				},
			}
		}

		if bucketWeb.HTTPIngress.IngressOverrides != nil {
			err := merge.Merge(ingress, bucketWeb.HTTPIngress.IngressOverrides)
			if err != nil {
				return ingress, reconciler.StatePresent, errors.WrapIf(err, "unable to merge overrides to ingress base object")
			}
		}

		return ingress, reconciler.StatePresent, nil
	}
	delete := &netv1.Ingress{
		ObjectMeta: meta,
	}
	return delete, reconciler.StateAbsent, nil
}
