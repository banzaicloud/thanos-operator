package rule

import (
	"github.com/banzaicloud/operator-tools/pkg/reconciler"
	"github.com/banzaicloud/thanos-operator/pkg/sdk/api/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func (r *ruleInstance) serviceMonitor() (runtime.Object, reconciler.DesiredState, error) {
	if r.Thanos.Spec.Rule != nil && r.Thanos.Spec.Rule.Metrics.ServiceMonitor {
		metrics := r.Thanos.Spec.Rule.Metrics
		serviceMonitor := &v1alpha1.ServiceMonitor{
			ObjectMeta: r.getMeta(),
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
					MatchLabels: r.getLabels(),
				},
				NamespaceSelector: v1alpha1.NamespaceSelector{
					MatchNames: []string{
						r.Thanos.Namespace,
					},
				},
				SampleLimit: 0,
			},
		}
		return serviceMonitor, reconciler.StatePresent, nil
	}
	delete := &v1alpha1.ServiceMonitor{
		ObjectMeta: r.getMeta(),
	}
	return delete, reconciler.StateAbsent, nil
}
