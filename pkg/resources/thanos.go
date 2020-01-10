package resources

import (
	"github.com/banzaicloud/operator-tools/pkg/reconciler"
	"github.com/banzaicloud/thanos-operator/pkg/sdk/api/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type ThanosComponentReconciler struct {
	Thanos      *v1alpha1.Thanos
	ObjectSores []v1alpha1.ObjectStore
	*reconciler.GenericResourceReconciler
}

func (t *ThanosComponentReconciler) Reconcile() (*reconcile.Result, error) {
	// Generate objects from resources
	object, state, err := t.queryDeployment()
	// Reconcile objects
	result, err := t.ReconcileResource(object, state)
	if err != nil {
		return &reconcile.Result{}, err
	}
	if result != nil {
		// short circuit if requested explicitly
		return result, err
	}
	return &reconcile.Result{}, err
}

func NewThanosComponentReconciler(thanos *v1alpha1.Thanos, objectStores *v1alpha1.ObjectStoreList, genericReconciler *reconciler.GenericResourceReconciler) *ThanosComponentReconciler {
	return &ThanosComponentReconciler{
		Thanos:                    thanos,
		ObjectSores:               objectStores.Items,
		GenericResourceReconciler: genericReconciler,
	}
}
