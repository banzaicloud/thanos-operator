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

package query

import (
	"emperror.dev/errors"
	"github.com/banzaicloud/operator-tools/pkg/merge"
	"github.com/banzaicloud/operator-tools/pkg/reconciler"
	netv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func (q *Query) ingressHTTP() (runtime.Object, reconciler.DesiredState, error) {
	if q.Thanos.Spec.Query != nil &&
		q.Thanos.Spec.Query.HTTPIngress != nil {
		queryIngress := q.Thanos.Spec.Query.HTTPIngress
		pathType := netv1.PathTypeImplementationSpecific
		ingress := &netv1.Ingress{
			ObjectMeta: q.Thanos.Spec.Query.MetaOverrides.Merge(q.getMeta(q.getName("http"))),
			Spec: netv1.IngressSpec{
				Rules: []netv1.IngressRule{
					{
						Host: queryIngress.Host,
						IngressRuleValue: netv1.IngressRuleValue{
							HTTP: &netv1.HTTPIngressRuleValue{
								Paths: []netv1.HTTPIngressPath{
									{
										Path:     queryIngress.Path,
										PathType: &pathType,
										Backend: netv1.IngressBackend{
											Service: &netv1.IngressServiceBackend{
												Name: q.getName(),
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
		if queryIngress.Certificate != "" {
			ingress.Spec.TLS = []netv1.IngressTLS{
				{
					Hosts:      []string{queryIngress.Host},
					SecretName: queryIngress.Certificate,
				},
			}
		}

		if q.Thanos.Spec.Query.HTTPIngress.IngressOverrides != nil {
			err := merge.Merge(ingress, q.Thanos.Spec.Query.HTTPIngress.IngressOverrides)
			if err != nil {
				return ingress, reconciler.StatePresent, errors.WrapIf(err, "unable to merge overrides to base object")
			}
		}

		return ingress, reconciler.StatePresent, nil
	}
	delete := &netv1.Ingress{
		ObjectMeta: q.getMeta(q.getName("http")),
	}
	return delete, reconciler.StateAbsent, nil
}

func (q *Query) ingressGRPC() (runtime.Object, reconciler.DesiredState, error) {
	if q.Thanos.Spec.Query != nil &&
		q.Thanos.Spec.Query.GRPCIngress != nil {
		queryIngress := q.Thanos.Spec.Query.GRPCIngress
		pathType := netv1.PathTypeImplementationSpecific
		ingress := &netv1.Ingress{
			ObjectMeta: q.Thanos.Spec.Query.MetaOverrides.Merge(q.getMeta(q.getName("grpc"))),
			Spec: netv1.IngressSpec{
				Rules: []netv1.IngressRule{
					{
						Host: queryIngress.Host,
						IngressRuleValue: netv1.IngressRuleValue{
							HTTP: &netv1.HTTPIngressRuleValue{
								Paths: []netv1.HTTPIngressPath{
									{
										Path:     queryIngress.Path,
										PathType: &pathType,
										Backend: netv1.IngressBackend{
											Service: &netv1.IngressServiceBackend{
												Name: q.getName(),
												Port: netv1.ServiceBackendPort{
													Name: "grpc",
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
		if queryIngress.Certificate != "" {
			ingress.Spec.TLS = []netv1.IngressTLS{
				{
					Hosts:      []string{queryIngress.Host},
					SecretName: queryIngress.Certificate,
				},
			}
		}

		if q.Thanos.Spec.Query.GRPCIngress.IngressOverrides != nil {
			err := merge.Merge(ingress, q.Thanos.Spec.Query.GRPCIngress.IngressOverrides)
			if err != nil {
				return ingress, reconciler.StatePresent, errors.WrapIf(err, "unable to merge overrides to base object")
			}
		}

		return ingress, reconciler.StatePresent, nil
	}
	delete := &netv1.Ingress{
		ObjectMeta: q.getMeta(q.getName("grpc")),
	}
	return delete, reconciler.StateAbsent, nil
}
