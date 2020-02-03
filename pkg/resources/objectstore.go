package resources

import (
	"fmt"

	"emperror.dev/errors"
	"github.com/banzaicloud/operator-tools/pkg/reconciler"
	"github.com/banzaicloud/thanos-operator/pkg/sdk/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type ObjectStoreReconciler struct {
	ObjectStore *v1alpha1.ObjectStore
	*reconciler.GenericResourceReconciler
}

func (t *ObjectStoreReconciler) GetCheck(port int32, path string) *corev1.Probe {
	return &corev1.Probe{
		Handler: corev1.Handler{
			HTTPGet: &corev1.HTTPGetAction{
				Path: path,
				Port: intstr.IntOrString{
					Type:   intstr.Int,
					IntVal: port,
				},
			},
		},
		InitialDelaySeconds: 5,
		TimeoutSeconds:      5,
		PeriodSeconds:       30,
		SuccessThreshold:    1,
		FailureThreshold:    2,
	}
}

func (t *ObjectStoreReconciler) GetCommonLabels() Labels {
	return Labels{
		ManagedByLabel: t.ObjectStore.Name,
	}
}

func (t *ObjectStoreReconciler) QualifiedName(name string) string {
	return fmt.Sprintf("%s-%s", t.ObjectStore.Name, name)
}

func (t *ObjectStoreReconciler) GetNameMeta(name string) metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Name:      name,
		Namespace: t.ObjectStore.Namespace,
	}
}

func (t *ObjectStoreReconciler) GetObjectMeta(name string) metav1.ObjectMeta {
	meta := t.GetNameMeta(name)
	meta.OwnerReferences = []metav1.OwnerReference{
		{
			APIVersion: t.ObjectStore.APIVersion,
			Kind:       t.ObjectStore.Kind,
			Name:       t.ObjectStore.Name,
			UID:        t.ObjectStore.UID,
			Controller: util.BoolPointer(true),
		},
	}
	return meta
}

func (t *ObjectStoreReconciler) ReconcileResources(resourceList []Resource) (*reconcile.Result, error) {
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

func NewObjectStoreReconciler(objectStore *v1alpha1.ObjectStore, genericReconciler *reconciler.GenericResourceReconciler) *ObjectStoreReconciler {
	return &ObjectStoreReconciler{
		ObjectStore:               objectStore,
		GenericResourceReconciler: genericReconciler,
	}
}
