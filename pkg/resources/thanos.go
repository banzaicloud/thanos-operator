package resources

import (
	"emperror.dev/errors"
	"github.com/banzaicloud/operator-tools/pkg/reconciler"
	"github.com/banzaicloud/thanos-operator/pkg/sdk/api/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const (
	nameLabel      = "app.kubernetes.io/name"
	instanceLabel  = "app.kubernetes.io/instance"
	versionLabel   = "app.kubernetes.io/version"
	componentLabel = "app.kubernetes.io/component"
	managedByLabel = "app.kubernetes.io/managed-by"
)

// Resource redeclaration of function with return type kubernetes Object
type Resource func() (runtime.Object, reconciler.DesiredState, error)

type ThanosComponentReconciler struct {
	Thanos      *v1alpha1.Thanos
	ObjectSores []v1alpha1.ObjectStore
	*reconciler.GenericResourceReconciler
}

func (t *ThanosComponentReconciler) Reconcile() (*reconcile.Result, error) {
	resourceList := []Resource{
		t.queryDeployment,
		t.storeDeployment,
		t.storeService,
	}
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

func NewThanosComponentReconciler(thanos *v1alpha1.Thanos, objectStores *v1alpha1.ObjectStoreList, genericReconciler *reconciler.GenericResourceReconciler) *ThanosComponentReconciler {
	return &ThanosComponentReconciler{
		Thanos:                    thanos,
		ObjectSores:               objectStores.Items,
		GenericResourceReconciler: genericReconciler,
	}
}
