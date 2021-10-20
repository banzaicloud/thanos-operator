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

package rule

import (
	"emperror.dev/errors"
	"github.com/banzaicloud/operator-tools/pkg/merge"
	"github.com/banzaicloud/operator-tools/pkg/reconciler"
	corev1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func (r *ruleInstance) ingressHTTP() (runtime.Object, reconciler.DesiredState, error) {
	if r.Thanos.Spec.Rule != nil &&
		r.Thanos.Spec.Rule.HTTPIngress != nil {
		rule := r.Thanos.Spec.Rule.DeepCopy()
		ruleIngress := rule.HTTPIngress
		pathType := netv1.PathTypeImplementationSpecific
		ingress := &netv1.Ingress{
			ObjectMeta: rule.MetaOverrides.Merge(r.getMeta("http")),
			Spec: netv1.IngressSpec{
				Rules: []netv1.IngressRule{
					{
						Host: ruleIngress.Host,
						IngressRuleValue: netv1.IngressRuleValue{
							HTTP: &netv1.HTTPIngressRuleValue{
								Paths: []netv1.HTTPIngressPath{
									{
										Path:     ruleIngress.Path,
										PathType: &pathType,
										Backend: netv1.IngressBackend{
											Service: &netv1.IngressServiceBackend{
												Name: r.getName(),
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
		if ruleIngress.Certificate != "" {
			ingress.Spec.TLS = []netv1.IngressTLS{
				{
					Hosts:      []string{ruleIngress.Host},
					SecretName: ruleIngress.Certificate,
				},
			}
		}
		if ruleIngress.IngressOverrides != nil {
			err := merge.Merge(ingress, ruleIngress.IngressOverrides)
			if err != nil {
				return ingress, reconciler.StatePresent, errors.WrapIf(err, "unable to merge overrides to base object")
			}
		}
		return ingress, reconciler.StatePresent, nil
	}
	delete := &corev1.Service{
		ObjectMeta: r.getMeta("http"),
	}
	return delete, reconciler.StateAbsent, nil
}

func (r *ruleInstance) ingressGRPC() (runtime.Object, reconciler.DesiredState, error) {
	if r.Thanos.Spec.Rule != nil &&
		r.Thanos.Spec.Rule.GRPCIngress != nil {
		ruleIngress := r.Thanos.Spec.Rule.GRPCIngress
		ingress := &netv1.Ingress{
			ObjectMeta: r.getMeta("grpc"),
			Spec: netv1.IngressSpec{
				Rules: []netv1.IngressRule{
					{
						Host: ruleIngress.Host,
						IngressRuleValue: netv1.IngressRuleValue{
							HTTP: &netv1.HTTPIngressRuleValue{
								Paths: []netv1.HTTPIngressPath{
									{
										Path: ruleIngress.Path,
										Backend: netv1.IngressBackend{
											Service: &netv1.IngressServiceBackend{
												Name: r.getName(),
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
		if ruleIngress.Certificate != "" {
			ingress.Spec.TLS = []netv1.IngressTLS{
				{
					Hosts:      []string{ruleIngress.Host},
					SecretName: ruleIngress.Certificate,
				},
			}
		}
		if ruleIngress.IngressOverrides != nil {
			err := merge.Merge(ingress, ruleIngress.IngressOverrides)
			if err != nil {
				return ingress, reconciler.StatePresent, errors.WrapIf(err, "unable to merge overrides to base object")
			}
		}
		return ingress, reconciler.StatePresent, nil
	}
	delete := &corev1.Service{
		ObjectMeta: r.getMeta("grpc"),
	}
	return delete, reconciler.StateAbsent, nil
}
