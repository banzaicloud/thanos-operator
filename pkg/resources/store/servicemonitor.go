package store

import (
	"github.com/banzaicloud/operator-tools/pkg/reconciler"
	"github.com/banzaicloud/thanos-operator/pkg/sdk/api/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func (s *storeInstance) serviceMonitor() (runtime.Object, reconciler.DesiredState, error) {
	if s.Thanos.Spec.StoreGateway != nil && s.Thanos.Spec.StoreGateway.Metrics.ServiceMonitor {
		metrics := s.Thanos.Spec.StoreGateway.Metrics
		serviceMonitor := &v1alpha1.ServiceMonitor{
			ObjectMeta: s.getMeta(),
			Spec: v1alpha1.ServiceMonitorSpec{
				Endpoints: []v1alpha1.Endpoint{
					{
						Port:          "http",
						Path:          metrics.Path,
						Interval:      metrics.Interval,
						ScrapeTimeout: metrics.Timeout,
					},
				},
				Selector: v1.LabelSelector{
					MatchLabels: s.getLabels(),
				},
				NamespaceSelector: v1alpha1.NamespaceSelector{
					MatchNames: []string{
						s.Thanos.Namespace,
					},
				},
				SampleLimit: 0,
			},
		}
		return serviceMonitor, reconciler.StatePresent, nil
	}
	delete := &v1alpha1.ServiceMonitor{
		ObjectMeta: s.getMeta(),
	}
	return delete, reconciler.StateAbsent, nil
}
