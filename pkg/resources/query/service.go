package query

import (
	"github.com/banzaicloud/operator-tools/pkg/reconciler"
	"github.com/banzaicloud/thanos-operator/pkg/resources"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func (q *Query) service() (runtime.Object, reconciler.DesiredState, error) {
	if q.Thanos.Spec.Query != nil {
		query := q.Thanos.Spec.Query.DeepCopy()
		queryService := &corev1.Service{
			ObjectMeta: q.getMeta(q.getName()),
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
				Selector: q.getLabels(),
				Type:     corev1.ServiceTypeClusterIP,
			},
		}
		return queryService, reconciler.StatePresent, nil

	}
	delete := &corev1.Service{
		ObjectMeta: q.getMeta(q.getName()),
	}
	return delete, reconciler.StateAbsent, nil
}
