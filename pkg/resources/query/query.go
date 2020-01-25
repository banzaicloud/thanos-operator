package query

import (
	"fmt"

	"github.com/banzaicloud/thanos-operator/pkg/resources"
	"github.com/banzaicloud/thanos-operator/pkg/resources/rule"
	"github.com/banzaicloud/thanos-operator/pkg/resources/store"
	"github.com/banzaicloud/thanos-operator/pkg/sdk/api/v1alpha1"
	"github.com/imdario/mergo"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func New(reconciler *resources.ThanosComponentReconciler) *Query {
	return &Query{
		ThanosComponentReconciler: reconciler,
	}
}

type Query struct {
	StoreEndpoints []v1alpha1.StoreEndpoint
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
		})
}

func (q *Query) getLabels() resources.Labels {
	labels := resources.Labels{
		resources.NameLabel: v1alpha1.QueryName,
	}.Merge(
		q.GetCommonLabels(),
	)
	if q.Thanos.Spec.Query != nil {
		labels.Merge(q.Thanos.Spec.Query.Labels)
	}
	return labels
}

func (q *Query) getMeta(name string, params ...string) metav1.ObjectMeta {
	namespace := ""
	if len(params) > 0 {
		namespace = params[0]
	}
	meta := q.GetObjectMeta(name, namespace)
	meta.Labels = q.getLabels()
	if q.Thanos.Spec.Query != nil {
		meta.Annotations = q.Thanos.Spec.Query.Annotations
	}
	return meta
}

func (q *Query) getStoreEndpoints() []string {
	var endpoints []string
	// Discover local StoreGateway
	if q.Thanos.Spec.StoreGateway != nil {
		for _, u := range store.New(q.ThanosComponentReconciler).GetServiceURLS() {
			endpoints = append(endpoints, fmt.Sprintf("--store=dnssrvnoa+%s", u))
		}
	}
	// Discover local Rule
	if q.Thanos.Spec.Rule != nil {
		for _, u := range rule.New(q.ThanosComponentReconciler).GetServiceURLS() {
			endpoints = append(endpoints, fmt.Sprintf("--store=dnssrvnoa+%s", u))
		}
	}
	// Discover StoreEndpoint aka Sidecars
	for _, endpoint := range q.StoreEndpoints {
		if url := endpoint.GetServiceURL(); url != "" {
			endpoints = append(endpoints, fmt.Sprintf("--store=dnssrvnoa+%s", url))
		}
	}
	return endpoints
}

func (q *Query) setArgs(args []string) []string {
	query := q.Thanos.Spec.Query.DeepCopy()
	args = append(args, resources.GetArgs(query)...)
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
	for _, s := range q.getStoreEndpoints() {
		args = append(args, s)
	}

	return args
}
