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
	"github.com/banzaicloud/thanos-operator/pkg/sdk/api/v1alpha1"
	netv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func (r receiverInstance) ingressGRPC() (runtime.Object, reconciler.DesiredState, error) {
	ingress := &netv1.Ingress{
		ObjectMeta: r.receiverGroup.MetaOverrides.Merge(r.getMeta(r.receiverGroup.Name + "-grpc")),
	}
	desiredState := reconciler.StateAbsent
	if ingressDef := r.receiverGroup.GRPCIngress; ingressDef != nil {
		configurator := ingressConfigurator{
			Ingress:         ingressDef,
			ServiceName:     r.getName(r.receiverGroup.Name),
			ServicePortName: "grpc",
		}
		if err := configurator.configureIngress(ingress); err != nil {
			return nil, nil, err
		}

		desiredState = reconciler.StatePresent
	}
	return ingress, desiredState, nil
}

func (r receiverInstance) ingressHTTP() (runtime.Object, reconciler.DesiredState, error) {
	ingress := &netv1.Ingress{
		ObjectMeta: r.receiverGroup.MetaOverrides.Merge(r.getMeta(r.receiverGroup.Name + "-remote-write")),
	}
	desiredState := reconciler.StateAbsent
	if ingressDef := r.receiverGroup.HTTPIngress; ingressDef != nil {
		configurator := ingressConfigurator{
			Ingress:         ingressDef,
			ServiceName:     r.getName(r.receiverGroup.Name),
			ServicePortName: "remote-write",
		}
		if err := configurator.configureIngress(ingress); err != nil {
			return nil, nil, err
		}

		desiredState = reconciler.StatePresent
	}
	return ingress, desiredState, nil
}

type ingressConfigurator struct {
	*v1alpha1.Ingress
	ServiceName     string
	ServicePortName string
}

func (i ingressConfigurator) configureIngress(ingress *netv1.Ingress) error {
	if i.Ingress == nil {
		return nil
	}

	pathType := netv1.PathTypeImplementationSpecific
	ingress.Spec = netv1.IngressSpec{
		Rules: []netv1.IngressRule{
			{
				Host: i.Host,
				IngressRuleValue: netv1.IngressRuleValue{
					HTTP: &netv1.HTTPIngressRuleValue{
						Paths: []netv1.HTTPIngressPath{
							{
								Path:     i.Path,
								PathType: &pathType,
								Backend: netv1.IngressBackend{
									Service: &netv1.IngressServiceBackend{
										Name: i.ServiceName,
										Port: netv1.ServiceBackendPort{
											Name: i.ServicePortName,
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
	if i.Certificate != "" {
		ingress.Spec.TLS = []netv1.IngressTLS{
			{
				Hosts:      []string{i.Host},
				SecretName: i.Certificate,
			},
		}
	}
	if i.IngressOverrides != nil {
		if err := merge.Merge(ingress, i.IngressOverrides); err != nil {
			return errors.WrapIf(err, "unable to merge overrides to base object")
		}
	}

	return nil
}
