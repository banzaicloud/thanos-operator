package sidecar

import (
	"github.com/banzaicloud/operator-tools/pkg/reconciler"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func (e *endpointService) ingressGRPC() (runtime.Object, reconciler.DesiredState, error) {
	if e.Spec.Ingress != nil {
		endpointIngress := e.Spec.Ingress
		ingress := &v1beta1.Ingress{
			ObjectMeta: e.getMeta(),
			Spec: v1beta1.IngressSpec{
				Rules: []v1beta1.IngressRule{
					{
						Host: endpointIngress.Host,
						IngressRuleValue: v1beta1.IngressRuleValue{
							HTTP: &v1beta1.HTTPIngressRuleValue{
								Paths: []v1beta1.HTTPIngressPath{
									{
										Path: endpointIngress.Path,
										Backend: v1beta1.IngressBackend{
											ServiceName: e.GetName(),
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
		if endpointIngress.Certificate != "" {
			ingress.Spec.TLS = []v1beta1.IngressTLS{
				{
					Hosts:      []string{endpointIngress.Host},
					SecretName: endpointIngress.Certificate,
				},
			}
		}
		return ingress, reconciler.StatePresent, nil
	}
	delete := &corev1.Service{
		ObjectMeta: e.getMeta(),
	}
	return delete, reconciler.StateAbsent, nil
}
