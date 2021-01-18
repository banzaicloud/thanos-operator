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
	"github.com/banzaicloud/operator-tools/pkg/reconciler"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func (q *Query) ingressHTTP() (runtime.Object, reconciler.DesiredState, error) {
	if q.Thanos.Spec.Query != nil &&
		q.Thanos.Spec.Query.HTTPIngress != nil {
		queryIngress := q.Thanos.Spec.Query.HTTPIngress
		query := q.Thanos.Spec.Query.DeepCopy()
		ingress := &v1beta1.Ingress{
			ObjectMeta: query.DeploymentOverrides.Merge(q.getMeta(q.getName("http"))),
			Spec: v1beta1.IngressSpec{
				Rules: []v1beta1.IngressRule{
					{
						Host: queryIngress.Host,
						IngressRuleValue: v1beta1.IngressRuleValue{
							HTTP: &v1beta1.HTTPIngressRuleValue{
								Paths: []v1beta1.HTTPIngressPath{
									{
										Path: queryIngress.Path,
										Backend: v1beta1.IngressBackend{
											ServiceName: q.getName(),
											ServicePort: intstr.IntOrString{
												Type:   intstr.String,
												StrVal: "http",
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
			ingress.Spec.TLS = []v1beta1.IngressTLS{
				{
					Hosts:      []string{queryIngress.Host},
					SecretName: queryIngress.Certificate,
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

func (q *Query) ingressGRPC() (runtime.Object, reconciler.DesiredState, error) {
	if q.Thanos.Spec.Query != nil &&
		q.Thanos.Spec.Query.GRPCIngress != nil {
		queryIngress := q.Thanos.Spec.Query.GRPCIngress
		query := q.Thanos.Spec.Query.DeepCopy()
		ingress := &v1beta1.Ingress{
			ObjectMeta: query.DeploymentOverrides.Merge(q.getMeta(q.getName("grpc"))),
			Spec: v1beta1.IngressSpec{
				Rules: []v1beta1.IngressRule{
					{
						Host: queryIngress.Host,
						IngressRuleValue: v1beta1.IngressRuleValue{
							HTTP: &v1beta1.HTTPIngressRuleValue{
								Paths: []v1beta1.HTTPIngressPath{
									{
										Path: queryIngress.Path,
										Backend: v1beta1.IngressBackend{
											ServiceName: q.getName(),
											ServicePort: intstr.IntOrString{
												Type:   intstr.String,
												StrVal: "grpc",
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
			ingress.Spec.TLS = []v1beta1.IngressTLS{
				{
					Hosts:      []string{queryIngress.Host},
					SecretName: queryIngress.Certificate,
				},
			}
		}
		return ingress, reconciler.StatePresent, nil
	}
	delete := &corev1.Service{
		ObjectMeta: q.getMeta(q.getName("grpc")),
	}
	return delete, reconciler.StateAbsent, nil
}
