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

package query_frontend

import (
	"fmt"
	"sort"

	"github.com/banzaicloud/operator-tools/pkg/utils"
	"github.com/banzaicloud/thanos-operator/pkg/resources"
	"github.com/banzaicloud/thanos-operator/pkg/resources/query"
	"github.com/banzaicloud/thanos-operator/pkg/sdk/api/v1alpha1"
	"github.com/imdario/mergo"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func New(reconciler *resources.ThanosComponentReconciler) *QueryFrontend {
	return &QueryFrontend{
		ThanosComponentReconciler: reconciler,
	}
}

type QueryFrontend struct {
	*resources.ThanosComponentReconciler
}

func (q *QueryFrontend) Reconcile() (*reconcile.Result, error) {
	if q.Thanos.Spec.QueryFrontend != nil {
		err := mergo.Merge(q.Thanos.Spec.QueryFrontend, v1alpha1.DefaultQueryFrontend)
		if err != nil {
			return nil, err
		}
	}
	return q.ReconcileResources(
		[]resources.Resource{
			q.deployment,
			q.service,
			//q.serviceMonitor,
			q.ingressHTTP,
		})
}

func (q *QueryFrontend) getLabels() resources.Labels {
	return utils.MergeLabels(
		resources.Labels{
			resources.NameLabel: v1alpha1.QueryFrontendName,
		},
		q.GetCommonLabels(),
	)
}

func (q *QueryFrontend) getName(suffix ...string) string {
	name := v1alpha1.QueryFrontendName
	if len(suffix) > 0 {
		name = name + "-" + suffix[0]
	}
	return q.QualifiedName(name)
}

func (q *QueryFrontend) GetHTTPService() string {
	return fmt.Sprintf("%s.%s.svc.%s", q.getName(), q.Thanos.Namespace, q.Thanos.GetClusterDomain())
}

func (q *QueryFrontend) GetHTTPServiceURL() string {
	return fmt.Sprintf("http://%s:%d", q.GetHTTPService(), resources.GetPort(q.Thanos.Spec.QueryFrontend.HttpAddress))
}

func (q *QueryFrontend) getMeta(name string) metav1.ObjectMeta {
	meta := q.GetObjectMeta(name)
	meta.Labels = q.getLabels()
	return meta
}

func (q *QueryFrontend) getQueryEndpoint() string {
	queryURL := ""
	if q.Thanos.Spec.QueryFrontend.QueryFrontendDownstreamURL != "" {
		queryURL = q.Thanos.Spec.QueryFrontend.QueryFrontendDownstreamURL
	} else if q.Thanos.Spec.Query != nil {
		queryURL = query.New(q.ThanosComponentReconciler).GetHTTPServiceURL()
	}
	return fmt.Sprintf("--query-frontend.downstream-url=%s", queryURL)
}

func (q *QueryFrontend) setArgs(originArgs []string) []string {
	queryFrontend := q.Thanos.Spec.QueryFrontend.DeepCopy()
	// Get args from the tags
	args := resources.GetArgs(queryFrontend)
	// Add discovery args
	args = append(args, q.getQueryEndpoint())

	// Sort generated args to prevent accidental diffs
	sort.Strings(args)
	// Concat original and computed args
	finalArgs := append(originArgs, args...)
	return finalArgs
}
