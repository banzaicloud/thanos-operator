package store

import (
	"fmt"

	"github.com/banzaicloud/thanos-operator/pkg/resources"
	"github.com/banzaicloud/thanos-operator/pkg/sdk/api/v1alpha1"
	"github.com/imdario/mergo"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func New(reconciler *resources.ThanosComponentReconciler) *Store {
	return &Store{
		ThanosComponentReconciler: reconciler,
	}
}

type storeInstance struct {
	*Store
	*v1alpha1.StoreEndpoint
}

func (s *storeInstance) getName() string {
	return fmt.Sprintf("%s-%s", s.StoreEndpoint.Name, v1alpha1.StoreName)
}

func (s *storeInstance) getSvc() string {
	return fmt.Sprintf("_grpc._tcp.%s.%s.svc.cluster.local", s.getName(), s.StoreEndpoint.Namespace)
}

func (s *Store) resourceFactory() []resources.Resource {
	var resourceList []resources.Resource

	for _, endpoint := range s.StoreEndpoints {
		resourceList = append(resourceList, (&storeInstance{s, &endpoint}).deployment)
		resourceList = append(resourceList, (&storeInstance{s, &endpoint}).service)
	}

	return resourceList
}

func (r *Store) GetServiceURLS() []string {
	var urls []string
	for _, endpoint := range r.StoreEndpoints {
		urls = append(urls, (&storeInstance{r, &endpoint}).getSvc())
	}
	return urls
}

func (s *Store) Reconcile() (*reconcile.Result, error) {
	if s.Thanos.Spec.StoreGateway != nil {
		err := mergo.Merge(s.Thanos.Spec.StoreGateway, v1alpha1.DefaultStoreGateway)
		if err != nil {
			return nil, err
		}
	}
	return s.ReconcileResources(s.resourceFactory())
}

type Store struct {
	*resources.ThanosComponentReconciler
}

func (s *Store) getLabels() resources.Labels {
	labels := resources.Labels{
		resources.NameLabel: v1alpha1.StoreName,
	}.Merge(
		s.GetCommonLabels(),
	)
	if s.Thanos.Spec.Query != nil {
		labels.Merge(s.Thanos.Spec.Query.Labels)
	}
	return labels
}

func (s *Store) getMeta(name string) metav1.ObjectMeta {
	meta := s.GetObjectMeta(name, "")
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
