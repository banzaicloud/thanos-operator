// Copyright 2021 Banzai Cloud
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

package receiver

import (
	prometheus "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	"github.com/banzaicloud/operator-tools/pkg/reconciler"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func (r *receiverInstance) serviceMonitor() (runtime.Object, reconciler.DesiredState, error) {
	if r.receiverGroup.Metrics != nil && r.receiverGroup.Metrics.ServiceMonitor {
		metrics := r.receiverGroup.Metrics
		serviceMonitor := &prometheus.ServiceMonitor{
			ObjectMeta: r.receiverGroup.MetaOverrides.Merge(r.getMeta(r.receiverGroup.Name)),
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
				SampleLimit: 0,
			},
		}
		return serviceMonitor, reconciler.StatePresent, nil
	}
	delete := &prometheus.ServiceMonitor{
		ObjectMeta: r.getMeta(r.receiverGroup.Name),
	}
	return delete, reconciler.StateAbsent, nil
}
