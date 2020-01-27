package query

import (
	"github.com/banzaicloud/operator-tools/pkg/reconciler"
	"github.com/banzaicloud/thanos-operator/pkg/resources"
	"github.com/banzaicloud/thanos-operator/pkg/sdk/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func (q *Query) service() (runtime.Object, reconciler.DesiredState, error) {
	if q.Thanos.Spec.Query != nil {
		query := q.Thanos.Spec.Query.DeepCopy()
		queryService := &corev1.Service{
			ObjectMeta: q.getMeta(v1alpha1.QueryName),
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
		ObjectMeta: q.getMeta(v1alpha1.QueryName),
	}
	return delete, reconciler.StateAbsent, nil
}
