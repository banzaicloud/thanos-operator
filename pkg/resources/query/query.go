package query

import (
	"fmt"

	"github.com/banzaicloud/thanos-operator/pkg/resources"
	"github.com/banzaicloud/thanos-operator/pkg/resources/store"
	"github.com/banzaicloud/thanos-operator/pkg/sdk/api/v1alpha1"
	"github.com/imdario/mergo"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const Name = "query"

func New(thanos *v1alpha1.Thanos, objectStores *v1alpha1.ObjectStoreList, reconciler *resources.ThanosComponentReconciler) *Query {
	return &Query{
		Thanos:                    thanos,
		ObjectSores:               objectStores.Items,
		ThanosComponentReconciler: reconciler,
	}
}

type Query struct {
	Thanos      *v1alpha1.Thanos
	ObjectSores []v1alpha1.ObjectStore
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
		resources.NameLabel: Name,
	}.Merge(
		q.GetCommonLabels(),
	)
	if q.Thanos.Spec.Query != nil {
		labels.Merge(q.Thanos.Spec.Query.Labels)
	}
	return labels
}

func (q *Query) getMeta(name string) metav1.ObjectMeta {
	meta := q.GetObjectMeta(name)
	meta.Labels = q.getLabels()
	if q.Thanos.Spec.Query != nil {
		meta.Annotations = q.Thanos.Spec.Query.Annotations
	}
	return meta
}

func (q *Query) getStoreEndpoints() []string {
	var endpoints []string
	if q.Thanos.Spec.StoreGateway != nil {
		endpoints = append(endpoints, q.QualifiedName(store.Name))
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
	if q.Thanos.Spec.ThanosDiscovery != nil {
		for _, s := range q.getStoreEndpoints() {
			args = append(args, fmt.Sprintf("--store=dnssrvnoa+_grpc._tcp.%s.%s.svc.cluster.local", s, q.Thanos.Namespace))
		}
	}
	return args
}
