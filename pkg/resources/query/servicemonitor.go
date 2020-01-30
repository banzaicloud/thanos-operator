package query

import (
	"github.com/banzaicloud/operator-tools/pkg/reconciler"
	"github.com/banzaicloud/thanos-operator/pkg/sdk/api/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func (q *Query) serviceMonitor() (runtime.Object, reconciler.DesiredState, error) {
	if q.Thanos.Spec.Query != nil && q.Thanos.Spec.Query.Metrics.ServiceMonitor {
		metrics := q.Thanos.Spec.Query.Metrics
		serviceMonitor := &v1alpha1.ServiceMonitor{
			ObjectMeta: q.getMeta(q.getName()),
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
					MatchLabels: q.getLabels(),
				},
				NamespaceSelector: v1alpha1.NamespaceSelector{
					MatchNames: []string{
						q.Thanos.Namespace,
					},
				},
				SampleLimit: 0,
			},
		}
		return serviceMonitor, reconciler.StatePresent, nil
	}
	delete := &v1alpha1.ServiceMonitor{
		ObjectMeta: q.getMeta(q.getName()),
	}
	return delete, reconciler.StateAbsent, nil
}
