package rule

import (
	"github.com/banzaicloud/operator-tools/pkg/reconciler"
	"github.com/banzaicloud/thanos-operator/pkg/resources"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func (r *Rule) service() (runtime.Object, reconciler.DesiredState, error) {
	if r.Thanos.Spec.Rule != nil {
		query := r.Thanos.Spec.Rule.DeepCopy()
		storeService := &corev1.Service{
			ObjectMeta: r.getMeta(Name),
			Spec: corev1.ServiceSpec{
				Ports: []corev1.ServicePort{
					{
						Name:     "grpc",
						Protocol: corev1.ProtocolTCP,
						Port:     resources.GetPort(query.GRPCAddress),
						TargetPort: intstr.IntOrString{
							Type:   intstr.String,
							StrVal: "grpc",
						},
					},
					{
						Name:     "http",
						Protocol: corev1.ProtocolTCP,
						Port:     resources.GetPort(query.HttpAddress),
						TargetPort: intstr.IntOrString{
							Type:   intstr.String,
							StrVal: "http",
						},
					},
				},
				Selector:  r.getLabels(Name),
				Type:      corev1.ServiceTypeClusterIP,
				ClusterIP: corev1.ClusterIPNone,
			},
		}
		return storeService, reconciler.StatePresent, nil

	}
	delete := &corev1.Service{
		ObjectMeta: r.getMeta(Name),
	}
	return delete, reconciler.StateAbsent, nil
}
