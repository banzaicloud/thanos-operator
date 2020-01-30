package compactor

import (
	"github.com/banzaicloud/operator-tools/pkg/reconciler"
	"github.com/banzaicloud/thanos-operator/pkg/sdk/api/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func (c *Compactor) serviceMonitor() (runtime.Object, reconciler.DesiredState, error) {
	if c.ObjectStore.Spec.Compactor != nil && c.ObjectStore.Spec.Compactor.Metrics.ServiceMonitor {
		metrics := c.ObjectStore.Spec.Compactor.Metrics
		serviceMonitor := &v1alpha1.ServiceMonitor{
			ObjectMeta: c.getMeta(),
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
					MatchLabels: c.getLabels(),
				},
				NamespaceSelector: v1alpha1.NamespaceSelector{
					MatchNames: []string{
						c.ObjectStore.Namespace,
					},
				},
				SampleLimit: 0,
			},
		}
		return serviceMonitor, reconciler.StatePresent, nil
	}
	delete := &v1alpha1.ServiceMonitor{
		ObjectMeta: c.getMeta(),
	}
	return delete, reconciler.StateAbsent, nil
}
