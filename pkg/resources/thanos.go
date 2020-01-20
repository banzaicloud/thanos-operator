package resources

import (
	"fmt"

	"emperror.dev/errors"
	"github.com/banzaicloud/logging-operator/pkg/sdk/util"
	"github.com/banzaicloud/operator-tools/pkg/reconciler"
	"github.com/banzaicloud/thanos-operator/pkg/sdk/api/v1alpha1"
	"github.com/imdario/mergo"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const (
	nameLabel      = "app.kubernetes.io/name"
	instanceLabel  = "app.kubernetes.io/instance"
	versionLabel   = "app.kubernetes.io/version"
	componentLabel = "app.kubernetes.io/component"
	managedByLabel = "app.kubernetes.io/managed-by"
)

type ThanosComponentReconciler struct {
	Thanos      *v1alpha1.Thanos
	ObjectSores []v1alpha1.ObjectStore
	*reconciler.GenericResourceReconciler
}

func (t *ThanosComponentReconciler) getCommonLabels() Labels {
	return Labels{
		managedByLabel: t.Thanos.Name,
	}
}

func (t *ThanosComponentReconciler) qualifiedName(name string) string {
	return fmt.Sprintf("%s-%s", t.Thanos.Name, name)
}

func (t *ThanosComponentReconciler) getNameMeta(name string) metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Name:      t.qualifiedName(name),
		Namespace: t.Thanos.Namespace,
	}
}

func (t *ThanosComponentReconciler) getObjectMeta(name string) metav1.ObjectMeta {
	meta := t.getNameMeta(name)
	meta.OwnerReferences = []metav1.OwnerReference{
		{
			APIVersion: t.Thanos.APIVersion,
			Kind:       t.Thanos.Kind,
			Name:       t.Thanos.Name,
			UID:        t.Thanos.UID,
			Controller: util.BoolPointer(true),
		},
	}
	return meta
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

func (t *ThanosComponentReconciler) setDefaults() error {
	if t.Thanos.Spec.Query != nil {
		err := mergo.Merge(t.Thanos.Spec.Query, v1alpha1.DefaultQuery)
		if err != nil {
			return err
		}
	}
	if t.Thanos.Spec.StoreGateway != nil {
		err := mergo.Merge(t.Thanos.Spec.StoreGateway, v1alpha1.DefaultStoreGateway)
		if err != nil {
			return err
		}
	}
	return nil
}

func NewThanosComponentReconciler(thanos *v1alpha1.Thanos, objectStores *v1alpha1.ObjectStoreList, genericReconciler *reconciler.GenericResourceReconciler) (*ThanosComponentReconciler, error) {
	reconciler := &ThanosComponentReconciler{
		Thanos:                    thanos,
		ObjectSores:               objectStores.Items,
		GenericResourceReconciler: genericReconciler,
	}
	err := reconciler.setDefaults()
	return reconciler, err
}
