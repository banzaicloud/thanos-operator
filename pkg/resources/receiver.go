// Copyright 2020 Banzai Cloud
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package resources

import (
	"fmt"

	"emperror.dev/errors"
	"github.com/banzaicloud/operator-tools/pkg/reconciler"
	"github.com/banzaicloud/operator-tools/pkg/utils"
	"github.com/banzaicloud/thanos-operator/pkg/sdk/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type ReceiverReconciler struct {
	*v1alpha1.Receiver
	*reconciler.GenericResourceReconciler
}

func (t *ReceiverReconciler) GetCheck(port int32, path string) *corev1.Probe {
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

func (t *ReceiverReconciler) GetClusterDomain() string {
	if t.Spec.ClusterDomain != "" {
		return t.Spec.ClusterDomain
	}
	return "cluster.local"
}

func (t *ReceiverReconciler) GetCommonLabels() Labels {
	return Labels{
		ManagedByLabel: t.Receiver.Name,
	}
}

func (t *ReceiverReconciler) QualifiedName(name string) string {
	return fmt.Sprintf("%s-%s", t.Receiver.Name, name)
}

func (t *ReceiverReconciler) GetNameMeta(name string, namespaceOverride string) metav1.ObjectMeta {
	namespace := t.Receiver.Namespace
	if namespaceOverride != "" {
		namespace = namespaceOverride
	}
	return metav1.ObjectMeta{
		Name:      name,
		Namespace: namespace,
	}
}

func (t *ReceiverReconciler) GetObjectMeta(name string) metav1.ObjectMeta {
	meta := t.GetNameMeta(name, "")
	meta.OwnerReferences = []metav1.OwnerReference{
		{
			APIVersion: t.Receiver.APIVersion,
			Kind:       t.Receiver.Kind,
			Name:       t.Receiver.Name,
			UID:        t.Receiver.UID,
			Controller: utils.BoolPointer(true),
		},
	}
	return meta
}

func (t *ReceiverReconciler) ReconcileResources(resourceList []Resource) (*reconcile.Result, error) {
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

func NewReceiverReconciler(receiver *v1alpha1.Receiver, genericReconciler *reconciler.GenericResourceReconciler) *ReceiverReconciler {
	return &ReceiverReconciler{
		Receiver:                  receiver,
		GenericResourceReconciler: genericReconciler,
	}
}
