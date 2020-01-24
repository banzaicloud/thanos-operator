package resources

import (
	"emperror.dev/errors"
	"github.com/banzaicloud/operator-tools/pkg/reconciler"
	"github.com/banzaicloud/thanos-operator/pkg/sdk/api/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type StoreEndpointComponentReconciler struct {
	StoreEndpoints []v1alpha1.StoreEndpoint
	*reconciler.GenericResourceReconciler
}

func (t *StoreEndpointComponentReconciler) ReconcileResources(resourceList []Resource) (*reconcile.Result, error) {
	// Generate objects from resources
	for _, res := range resourceList {
		o, state, err := res()
		if err != nil {
			return nil, errors.WrapIf(err, "failed to create desired object")
		}
		if o == nil {
			return nil, errors.Errorf("Reconcile error! Resource %#v returns with nil object", res)
		}
		result, err := t.ReconcileResource(o, state)
		if err != nil {
			return nil, errors.WrapIf(err, "failed to reconcile resource")
		}
		if result != nil {
			return result, nil
		}
	}

	return nil, nil
}

func NewStoreEndpointComponentReconciler(storeEndpoints *v1alpha1.StoreEndpointList, genericReconciler *reconciler.GenericResourceReconciler) *StoreEndpointComponentReconciler {
	return &StoreEndpointComponentReconciler{
		StoreEndpoints:            storeEndpoints.Items,
		GenericResourceReconciler: genericReconciler,
	}
}
