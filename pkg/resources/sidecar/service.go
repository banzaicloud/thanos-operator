package sidecar

import (
	"github.com/banzaicloud/operator-tools/pkg/reconciler"
	"github.com/banzaicloud/thanos-operator/pkg/sdk/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
)

type endpointService struct {
	*v1alpha1.StoreEndpoint
}

func (e *endpointService) sidecarService() (runtime.Object, reconciler.DesiredState, error) {
	if e.Spec.Selector != nil {
		var grpcPort int32 = 10901
		var httpPort int32 = 10902
		labels := map[string]string{
			"app": "prometheus",
		}
		if e.Spec.Selector.GRPCPort != 0 {
			grpcPort = e.Spec.Selector.GRPCPort
		}
		if e.Spec.Selector.HTTPPort != 0 {
			httpPort = e.Spec.Selector.HTTPPort
		}
		if e.Spec.Selector.Labels != nil {
			labels = e.Spec.Selector.Labels
		}
		storeService := &corev1.Service{
			ObjectMeta: getMeta(e.StoreEndpoint),
			Spec: corev1.ServiceSpec{
				Ports: []corev1.ServicePort{
					{
						Name:     "grpc",
						Protocol: corev1.ProtocolTCP,
						Port:     grpcPort,
						TargetPort: intstr.IntOrString{
							Type:   intstr.String,
							StrVal: "grpc",
						},
					},
					{
						Name:     "http",
						Protocol: corev1.ProtocolTCP,
						Port:     httpPort,
						TargetPort: intstr.IntOrString{
							Type:   intstr.String,
							StrVal: "http",
						},
					},
				},
				Selector:  labels,
				Type:      corev1.ServiceTypeClusterIP,
				ClusterIP: corev1.ClusterIPNone,
			},
		}
		return storeService, reconciler.StatePresent, nil
	}
	delete := &corev1.Service{
		ObjectMeta: getMeta(e.StoreEndpoint),
	}
	return delete, reconciler.StateAbsent, nil
}
