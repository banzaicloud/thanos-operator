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
		q.Thanos.Spec.Query.Ingress != nil {
		queryIngress := q.Thanos.Spec.Query.Ingress
		ingress := &v1beta1.Ingress{
			ObjectMeta: q.getMeta(q.getName()),
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
		ObjectMeta: q.getMeta(q.getName()),
	}
	return delete, reconciler.StateAbsent, nil
}
