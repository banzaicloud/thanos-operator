package store

import (
	"github.com/banzaicloud/thanos-operator/pkg/resources"
	"github.com/banzaicloud/thanos-operator/pkg/sdk/api/v1alpha1"
	"github.com/imdario/mergo"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const Name = "store"

func New(thanos *v1alpha1.Thanos, objectStores *v1alpha1.ObjectStoreList, reconciler *resources.ThanosComponentReconciler) *Store {
	return &Store{
		Thanos:                    thanos,
		ObjectSores:               objectStores.Items,
		ThanosComponentReconciler: reconciler,
	}
}

func (s *Store) Reconcile() (*reconcile.Result, error) {
	if s.Thanos.Spec.StoreGateway != nil {
		err := mergo.Merge(s.Thanos.Spec.StoreGateway, v1alpha1.DefaultStoreGateway)
		if err != nil {
			return nil, err
		}
	}
	return s.ReconcileResources(
		[]resources.Resource{
			s.deployment,
			s.service,
		})
}

type Store struct {
	Thanos      *v1alpha1.Thanos
	ObjectSores []v1alpha1.ObjectStore
	*resources.ThanosComponentReconciler
}

func (s *Store) getLabels() resources.Labels {
	labels := resources.Labels{
		resources.NameLabel: Name,
	}.Merge(
		s.GetCommonLabels(),
	)
	if s.Thanos.Spec.Query != nil {
		labels.Merge(s.Thanos.Spec.Query.Labels)
	}
	return labels
}

func (s *Store) getMeta(name string) metav1.ObjectMeta {
	meta := s.GetObjectMeta(name)
	meta.Labels = s.getLabels()
	if s.Thanos.Spec.StoreGateway != nil {
		meta.Annotations = s.Thanos.Spec.StoreGateway.Annotations
	}
	return meta
}

func (s *Store) setArgs(args []string) []string {
	store := s.Thanos.Spec.StoreGateway.DeepCopy()
	args = append(args, resources.GetArgs(store)...)
	return args
}
