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

package query_frontend

import (
	"github.com/banzaicloud/operator-tools/pkg/reconciler"
	corev1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func (q *QueryFrontend) ingressHTTP() (runtime.Object, reconciler.DesiredState, error) {
	if q.Thanos.Spec.QueryFrontend != nil &&
		q.Thanos.Spec.QueryFrontend.HTTPIngress != nil {
		queryFrontendIngress := q.Thanos.Spec.QueryFrontend.HTTPIngress
		queryFrontend := q.Thanos.Spec.QueryFrontend.DeepCopy()
		pathType := netv1.PathTypeImplementationSpecific
		ingress := &netv1.Ingress{
			ObjectMeta: queryFrontend.MetaOverrides.Merge(q.getMeta(q.getName("http"))),
			Spec: netv1.IngressSpec{
				Rules: []netv1.IngressRule{
					{
						Host: queryFrontendIngress.Host,
						IngressRuleValue: netv1.IngressRuleValue{
							HTTP: &netv1.HTTPIngressRuleValue{
								Paths: []netv1.HTTPIngressPath{
									{
										Path:     queryFrontendIngress.Path,
										PathType: &pathType,
										Backend: netv1.IngressBackend{
											ServiceName: q.getName(),
											ServicePort: intstr.FromString("http"),
										},
									},
								},
							},
						},
					},
				},
			},
		}
		if queryFrontendIngress.Certificate != "" {
			ingress.Spec.TLS = []netv1.IngressTLS{
				{
					Hosts:      []string{queryFrontendIngress.Host},
					SecretName: queryFrontendIngress.Certificate,
				},
			}
		}
		return ingress, reconciler.StatePresent, nil
	}
	delete := &corev1.Service{
		ObjectMeta: q.getMeta(q.getName("http")),
	}
	return delete, reconciler.StateAbsent, nil
}
