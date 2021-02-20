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

package receiver

import (
	"emperror.dev/errors"
	"github.com/banzaicloud/operator-tools/pkg/merge"
	"github.com/banzaicloud/operator-tools/pkg/reconciler"
	netv1 "k8s.io/api/networking/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func (r *receiverInstance) ingressGRPC() (runtime.Object, reconciler.DesiredState, error) {
	if r.receiverGroup.GRPCIngress != nil {
		endpointIngress := r.receiverGroup.GRPCIngress
		pathType := netv1.PathTypeImplementationSpecific
		ingress := &netv1.Ingress{
			ObjectMeta: r.receiverGroup.MetaOverrides.Merge(r.getMeta(r.receiverGroup.Name + "-grpc")),
			Spec: netv1.IngressSpec{
				Rules: []netv1.IngressRule{
					{
						Host: endpointIngress.Host,
						IngressRuleValue: netv1.IngressRuleValue{
							HTTP: &netv1.HTTPIngressRuleValue{
								Paths: []netv1.HTTPIngressPath{
									{
										Path:     endpointIngress.Path,
										PathType: &pathType,
										Backend: netv1.IngressBackend{
											ServiceName: r.getName(r.receiverGroup.Name),
											ServicePort: intstr.FromString("grpc"),
										},
									},
								},
							},
						},
					},
				},
			},
		}
		if endpointIngress.Certificate != "" {
			ingress.Spec.TLS = []netv1.IngressTLS{
				{
					Hosts:      []string{endpointIngress.Host},
					SecretName: endpointIngress.Certificate,
				},
			}
		}
		return ingress, reconciler.StatePresent, nil
	}
	delete := &netv1.Ingress{
		ObjectMeta: r.getMeta(r.receiverGroup.Name + "-grpc"),
	}
	return delete, reconciler.StateAbsent, nil
}

func (r *receiverInstance) ingressHTTP() (runtime.Object, reconciler.DesiredState, error) {
	if r.receiverGroup.HTTPIngress != nil {
		endpointIngress := r.receiverGroup.HTTPIngress
		pathType := netv1.PathTypeImplementationSpecific
		ingress := &netv1.Ingress{
			ObjectMeta: r.receiverGroup.MetaOverrides.Merge(r.getMeta(r.receiverGroup.Name + "-remote-write")),
			Spec: netv1.IngressSpec{
				Rules: []netv1.IngressRule{
					{
						Host: endpointIngress.Host,
						IngressRuleValue: netv1.IngressRuleValue{
							HTTP: &netv1.HTTPIngressRuleValue{
								Paths: []netv1.HTTPIngressPath{
									{
										Path:     endpointIngress.Path,
										PathType: &pathType,
										Backend: netv1.IngressBackend{
											ServiceName: r.getName(r.receiverGroup.Name),
											ServicePort: intstr.FromString("remote-write"),
										},
									},
								},
							},
						},
					},
				},
			},
		}
		if endpointIngress.Certificate != "" {
			ingress.Spec.TLS = []netv1.IngressTLS{
				{
					Hosts:      []string{endpointIngress.Host},
					SecretName: endpointIngress.Certificate,
				},
			}
		}
		if endpointIngress.IngressOverrides != nil {
			err := merge.Merge(ingress, endpointIngress.IngressOverrides)
			if err != nil {
				return ingress, reconciler.StatePresent, errors.WrapIf(err, "unable to merge overrides to base object")
			}
		}

		return ingress, reconciler.StatePresent, nil
	}
	delete := &netv1.Ingress{
		ObjectMeta: r.getMeta(r.receiverGroup.Name + "-grpc"),
	}
	return delete, reconciler.StateAbsent, nil
}
