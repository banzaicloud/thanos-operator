package resources

import (
	"fmt"
	"strconv"
	"strings"

	"emperror.dev/errors"
	"github.com/banzaicloud/logging-operator/pkg/sdk/util"
	"github.com/banzaicloud/operator-tools/pkg/reconciler"
	"github.com/banzaicloud/thanos-operator/pkg/sdk/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func GetPort(address string) int32 {
	res := strings.Split(address, ":")
	if len(res) > 1 {
		port, err := strconv.Atoi(res[1])
		if err != nil {
			return 0
		}
		return int32(port)
	}
	return 0
}

type ThanosComponentReconciler struct {
	Thanos         *v1alpha1.Thanos
	StoreEndpoints []v1alpha1.StoreEndpoint
	*reconciler.GenericResourceReconciler
}

func (t *ThanosComponentReconciler) GetCheck(port int32, path string) *corev1.Probe {
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

func (t *ThanosComponentReconciler) GetCommonLabels() Labels {
	return Labels{
		ManagedByLabel: t.Thanos.Name,
	}
}

func (t *ThanosComponentReconciler) QualifiedName(name string) string {
	return fmt.Sprintf("%s-%s", t.Thanos.Name, name)
}

func (t *ThanosComponentReconciler) GetNameMeta(name string, namespaceOverride string) metav1.ObjectMeta {
	namespace := t.Thanos.Namespace
	if namespaceOverride != "" {
		namespace = namespaceOverride
	}
	return metav1.ObjectMeta{
		Name:      t.QualifiedName(name),
		Namespace: namespace,
	}
}

func (t *ThanosComponentReconciler) GetObjectMeta(name string, namespaceOverride string) metav1.ObjectMeta {
	meta := t.GetNameMeta(name, namespaceOverride)
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

func (t *ThanosComponentReconciler) ReconcileResources(resourceList []Resource) (*reconcile.Result, error) {
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

func NewThanosComponentReconciler(thanos *v1alpha1.Thanos, storeEndpoints []v1alpha1.StoreEndpoint, genericReconciler *reconciler.GenericResourceReconciler) *ThanosComponentReconciler {
	return &ThanosComponentReconciler{
		Thanos:                    thanos,
		StoreEndpoints:            storeEndpoints,
		GenericResourceReconciler: genericReconciler,
	}
}
