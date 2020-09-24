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

package store

import (
	"github.com/banzaicloud/operator-tools/pkg/prometheus"
	"github.com/banzaicloud/operator-tools/pkg/reconciler"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func (s *storeInstance) serviceMonitor() (runtime.Object, reconciler.DesiredState, error) {
	if s.Thanos.Spec.StoreGateway != nil && s.Thanos.Spec.StoreGateway.Metrics.ServiceMonitor {
		metrics := s.Thanos.Spec.StoreGateway.Metrics
		meta := s.getMeta()
		meta.Name = "thanos-store-" + meta.Name
		serviceMonitor := &prometheus.ServiceMonitor{
			ObjectMeta: s.StoreEndpoint.Spec.MetaOverrides.Merge(meta),
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
					MatchLabels: s.getLabels(),
				},
				NamespaceSelector: prometheus.NamespaceSelector{
					MatchNames: []string{
						s.Thanos.Namespace,
					},
				},
				SampleLimit: 0,
			},
		}
		return serviceMonitor, reconciler.StatePresent, nil
	}
	delete := &prometheus.ServiceMonitor{
		ObjectMeta: s.getMeta(),
	}
	return delete, reconciler.StateAbsent, nil
}
