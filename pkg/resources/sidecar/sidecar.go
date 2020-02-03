package sidecar

import (
	"fmt"

	"github.com/banzaicloud/operator-tools/pkg/utils"
	"github.com/banzaicloud/thanos-operator/pkg/resources"
	"github.com/banzaicloud/thanos-operator/pkg/sdk/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func New(storeEndpoints *v1alpha1.StoreEndpointList, reconciler *resources.StoreEndpointComponentReconciler) *Sidecar {
	return &Sidecar{
		StoreEndpoints:                   storeEndpoints.Items,
		StoreEndpointComponentReconciler: reconciler,
	}
}

type Sidecar struct {
	StoreEndpoints []v1alpha1.StoreEndpoint
	*resources.StoreEndpointComponentReconciler
}

func getMeta(endpoint *v1alpha1.StoreEndpoint) metav1.ObjectMeta {
	meta := metav1.ObjectMeta{
		Name:      fmt.Sprintf("%s-%s", endpoint.Name, v1alpha1.SidecarName),
		Namespace: endpoint.Namespace,
	}
	meta.OwnerReferences = []metav1.OwnerReference{
		{
			APIVersion: endpoint.APIVersion,
			Kind:       endpoint.Kind,
			Name:       endpoint.Name,
			UID:        endpoint.UID,
			Controller: utils.BoolPointer(true),
		},
	}
	meta.Labels = endpoint.Labels
	meta.Annotations = endpoint.Annotations
	return meta
}

func (s *Sidecar) serviceFactory() []resources.Resource {
	var serviceList []resources.Resource

	for _, endpoint := range s.StoreEndpoints {
		serviceList = append(serviceList, (&endpointService{&endpoint}).sidecarService)
	}

	return serviceList
}

func (s *Sidecar) Reconcile() (*reconcile.Result, error) {
	return s.ReconcileResources(s.serviceFactory())
}
