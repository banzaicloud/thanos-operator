// Copyright © 2020 Banzai Cloud
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package controllers

import (
	"context"

	"github.com/go-logr/logr"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	monitoringv1alpha1 "github.com/banzaicloud/thanos-operator/pkg/sdk/api/v1alpha1"
)

// ObjectStoreReconciler reconciles a ObjectStore object
type ObjectStoreReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=monitoring.banzaicloud.io,resources=objectstores,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=monitoring.banzaicloud.io,resources=objectstores/status,verbs=get;update;patch

func (r *ObjectStoreReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	result := ctrl.Result{}
	ctx := context.Background()
	log := r.Log.WithValues("objectstore", req.NamespacedName)

	store := &monitoringv1alpha1.ObjectStore{}
	err := r.Client.Get(ctx, req.NamespacedName, store)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return result, nil
		}
		return result, err
	}
	store, err = store.SetDefaults()
	if err != nil {
		return result, err
	}

	log.V(1).Info("got", "spec", store.Spec)

	log.V(1).Info("reconcile finished", "result", result)

	return result, nil
}

func (r *ObjectStoreReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&monitoringv1alpha1.ObjectStore{}).
		Complete(r)
}
