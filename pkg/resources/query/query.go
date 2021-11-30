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

package query

import (
	"fmt"
	"regexp"
	"sort"

	"github.com/banzaicloud/operator-tools/pkg/utils"
	"github.com/banzaicloud/thanos-operator/pkg/resources"
	"github.com/banzaicloud/thanos-operator/pkg/resources/rule"
	"github.com/banzaicloud/thanos-operator/pkg/resources/store"
	"github.com/banzaicloud/thanos-operator/pkg/sdk/api/v1alpha1"
	"github.com/imdario/mergo"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const serverCertMountPath = "/etc/tls/server"
const clientCertMountPath = "/etc/tls/client"
const clientCACertMountPath = "/etc/tls/client-ca"
const serverCACertMountPath = "/etc/tls/server-ca"

func New(reconciler *resources.ThanosComponentReconciler) *Query {
	return &Query{
		ThanosComponentReconciler: reconciler,
	}
}

type Query struct {
	*resources.ThanosComponentReconciler
}

func (q *Query) Reconcile() (*reconcile.Result, error) {
	if q.Thanos.Spec.Query != nil {
		err := mergo.Merge(q.Thanos.Spec.Query, v1alpha1.DefaultQuery)
		if err != nil {
			return nil, err
		}
	}
	return q.ReconcileResources(
		[]resources.Resource{
			q.deployment,
			q.service,
			q.serviceMonitor,
			q.ingressHTTP,
			q.ingressGRPC,
			q.grafanaDatasource,
		})
}

func (q *Query) getLabels() resources.Labels {
	return utils.MergeLabels(
		resources.Labels{
			resources.NameLabel: v1alpha1.QueryName,
		},
		q.GetCommonLabels(),
	)
}

func (q *Query) getName(suffix ...string) string {
	name := v1alpha1.QueryName
	if len(suffix) > 0 {
		name = name + "-" + suffix[0]
	}
	return q.QualifiedName(name)
}

func (q *Query) GetGRPCService() string {
	return fmt.Sprintf("_grpc._tcp.%s.%s.svc.%s", q.getName(), q.Thanos.Namespace, q.Thanos.GetClusterDomain())
}

func (q *Query) GetHTTPService() string {
	return fmt.Sprintf("%s.%s.svc.%s", q.getName(), q.Thanos.Namespace, q.Thanos.GetClusterDomain())
}

func (q *Query) GetHTTPServiceURL() string {
	return fmt.Sprintf("http://%s:%d", q.GetHTTPService(), resources.GetPort(q.Thanos.Spec.Query.HttpAddress))
}

func (q *Query) getMeta(name string) metav1.ObjectMeta {
	meta := q.GetObjectMeta(name)
	meta.Labels = q.getLabels()
	return meta
}

func (q *Query) getStoreEndpoints() []string {
	var endpoints []string
	// Discover all query instance
	if q.Thanos.Spec.QueryDiscovery {
		for _, t := range q.ThanosList {
			if t.Spec.Query != nil {
				reconciler := resources.NewThanosComponentReconciler(t.DeepCopy(), nil, nil, nil)
				svc := (&Query{reconciler}).GetGRPCService()
				endpoints = append(endpoints, fmt.Sprintf("--store=dnssrvnoa+%s", svc))
			}
		}
	}
	// Discover local StoreGateway
	if q.Thanos.Spec.StoreGateway != nil {
		for _, u := range store.New(q.ThanosComponentReconciler).GetServiceURLS() {
			endpoints = append(endpoints, fmt.Sprintf("--store=dnssrvnoa+%s.%s", u, q.Thanos.GetClusterDomain()))
		}
	}
	// Discover local Rule
	if q.Thanos.Spec.Rule != nil {
		for _, u := range rule.New(q.ThanosComponentReconciler).GetServiceURLS() {
			endpoints = append(endpoints, fmt.Sprintf("--store=dnssrvnoa+%s.%s", u, q.Thanos.GetClusterDomain()))
		}
	}
	// Discover StoreEndpoint aka Sidecars
	for _, endpoint := range q.StoreEndpoints {
		if url := endpoint.GetServiceURL(); url != "" {
			r, _ := regexp.Compile(`.*\.svc$`)
			switch r.MatchString(url) {
			case true:
				endpoints = append(endpoints, fmt.Sprintf("--store=%s.%s", url, q.Thanos.GetClusterDomain()))
			default:
				endpoints = append(endpoints, fmt.Sprintf("--store=%s", url))
			}
		}
	}
	// Discover static StoreAPI endpoints provided as stores parameter
	for _, endpoint := range q.Thanos.Spec.Query.Stores {
		endpoints = append(endpoints, fmt.Sprintf("--store=%s", endpoint))
	}
	return endpoints
}

func (q *Query) setArgs(originArgs []string) []string {
	query := q.Thanos.Spec.Query.DeepCopy()
	// Get args from the tags
	args := resources.GetArgs(query)

	if query.GRPCClientCertificate != "" {
		args = append(args, "--grpc-client-tls-secure")
		args = append(args, fmt.Sprintf("--grpc-client-tls-cert=%s/%s", clientCertMountPath, "tls.crt"))
		args = append(args, fmt.Sprintf("--grpc-client-tls-key=%s/%s", clientCertMountPath, "tls.key"))
		if query.GRPCClientCA != "" {
			args = append(args, fmt.Sprintf("--grpc-client-tls-ca=%s/%s", clientCACertMountPath, "ca.crt"))
		} else {
			// for backward compatibility only
			args = append(args, fmt.Sprintf("--grpc-client-tls-ca=%s/%s", clientCertMountPath, "ca.crt"))
		}
		if query.GRPCClientServerName != "" {
			args = append(args, fmt.Sprintf("--grpc-client-server-name=%s", query.GRPCClientServerName))
		}
	}
	if query.GRPCServerCertificate != "" {
		args = append(args, fmt.Sprintf("--grpc-server-tls-cert=%s/%s", serverCertMountPath, "tls.crt"))
		args = append(args, fmt.Sprintf("--grpc-server-tls-key=%s/%s", serverCertMountPath, "tls.key"))
		if query.GRPCServerCA != "" {
			args = append(args, fmt.Sprintf("--grpc-server-tls-client-ca=%s/%s", serverCACertMountPath, "ca.crt"))
		} else {
			// for backward compatibility only
			args = append(args, fmt.Sprintf("--grpc-server-tls-client-ca=%s/%s", serverCertMountPath, "ca.crt"))
		}
	}
	// Handle special args
	if query.QueryReplicaLabels != nil {
		for _, l := range query.QueryReplicaLabels {
			args = append(args, fmt.Sprintf("--query.replica-label=%s", l))
		}
	}
	if query.SelectorLabels != nil {
		for k, v := range query.SelectorLabels {
			args = append(args, fmt.Sprintf("--selector-label=%s=%s", k, v))
		}
	}
	// Add discovery args
	args = append(args, q.getStoreEndpoints()...)

	// Sort generated args to prevent accidental diffs
	sort.Strings(args)
	// Concat original and computed args
	finalArgs := append(originArgs, args...)
	return finalArgs
}
