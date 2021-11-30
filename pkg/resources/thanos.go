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
	"strconv"
	"strings"

	"github.com/banzaicloud/operator-tools/pkg/reconciler"
	"github.com/banzaicloud/operator-tools/pkg/utils"
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

func GetProbe(port int32, path string) *corev1.Probe {
	return &corev1.Probe{
		Handler: corev1.Handler{
			HTTPGet: &corev1.HTTPGetAction{
				Path: path,
				Port: intstr.FromInt(int(port)),
			},
		},
		InitialDelaySeconds: 5,
		TimeoutSeconds:      5,
		PeriodSeconds:       30,
		SuccessThreshold:    1,
		FailureThreshold:    2,
	}
}

type ThanosComponentReconciler struct {
	Thanos         *v1alpha1.Thanos
	ThanosList     []v1alpha1.Thanos
	StoreEndpoints []v1alpha1.StoreEndpoint
	*reconciler.GenericResourceReconciler
}

func (t *ThanosComponentReconciler) GetCommonLabels() Labels {
	return Labels{
		ManagedByLabel: t.Thanos.Name,
	}
}

func (t *ThanosComponentReconciler) QualifiedName(name string) string {
	return fmt.Sprintf("%s-%s", t.Thanos.Name, name)
}

func (t *ThanosComponentReconciler) GetNameMeta(name string) metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Name:      name,
		Namespace: t.Thanos.Namespace,
	}
}

func (t *ThanosComponentReconciler) GetObjectMeta(name string) metav1.ObjectMeta {
	meta := t.GetNameMeta(name)
	meta.OwnerReferences = []metav1.OwnerReference{
		{
			APIVersion: t.Thanos.APIVersion,
			Kind:       t.Thanos.Kind,
			Name:       t.Thanos.Name,
			UID:        t.Thanos.UID,
			Controller: utils.BoolPointer(true),
		},
	}
	return meta
}

func (t *ThanosComponentReconciler) ReconcileResources(resourceList []Resource) (*reconcile.Result, error) {
	return Dispatch(t.GenericResourceReconciler, resourceList)
}

func NewThanosComponentReconciler(thanos *v1alpha1.Thanos, thanosList []v1alpha1.Thanos, storeEndpoints []v1alpha1.StoreEndpoint, genericReconciler *reconciler.GenericResourceReconciler) *ThanosComponentReconciler {
	return &ThanosComponentReconciler{
		Thanos:                    thanos,
		ThanosList:                thanosList,
		StoreEndpoints:            storeEndpoints,
		GenericResourceReconciler: genericReconciler,
	}
}
