// Copyright 2020 Banzai Cloud
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package rule

import (
	"github.com/banzaicloud/operator-tools/pkg/prometheus"
	"github.com/banzaicloud/operator-tools/pkg/reconciler"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func (r *ruleInstance) serviceMonitor() (runtime.Object, reconciler.DesiredState, error) {
	if r.Thanos.Spec.Rule != nil && r.Thanos.Spec.Rule.Metrics.ServiceMonitor {
		metrics := r.Thanos.Spec.Rule.Metrics
		rule := r.Thanos.Spec.Rule.DeepCopy()
		meta := r.getMeta()
		meta.Name = "thanos-rule-" + meta.Name
		serviceMonitor := &prometheus.ServiceMonitor{
			ObjectMeta: rule.MetaOverrides.Merge(meta),
			Spec: prometheus.ServiceMonitorSpec{
				Endpoints: []prometheus.Endpoint{
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
				NamespaceSelector: prometheus.NamespaceSelector{
					MatchNames: []string{
						r.Thanos.Namespace,
					},
				},
				SampleLimit: 0,
			},
		}
		return serviceMonitor, reconciler.StatePresent, nil
	}
	delete := &prometheus.ServiceMonitor{
		ObjectMeta: r.getMeta(),
	}
	return delete, reconciler.StateAbsent, nil
}
