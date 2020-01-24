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

	"github.com/banzaicloud/operator-tools/pkg/reconciler"
	"github.com/banzaicloud/thanos-operator/pkg/resources"
	"github.com/banzaicloud/thanos-operator/pkg/resources/query"
	"github.com/banzaicloud/thanos-operator/pkg/resources/rule"
	"github.com/banzaicloud/thanos-operator/pkg/resources/store"
	"github.com/banzaicloud/thanos-operator/pkg/sdk/api/v1alpha1"
	"github.com/go-logr/logr"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// ThanosReconciler reconciles a Thanos object
type ThanosReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=monitoring.banzaicloud.io,resources=thanos,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=monitoring.banzaicloud.io,resources=thanos/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=extensions;apps,resources=deployments;statefulsets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="",resources=services,verbs=get;list;watch;create;update;patch;delete

func (r *ThanosReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	_ = context.Background()
	_ = r.Log.WithValues("thanos", req.NamespacedName)

	thanos := &v1alpha1.Thanos{}
	err := r.Client.Get(context.TODO(), req.NamespacedName, thanos)
	if err != nil {
		// Object not found, return.  Created objects are automatically garbage collected.
		// For additional cleanup logic use finalizers.
		if apierrors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}
	// Collect ObjectStores TODO better way to handle this
	storeEndpoints := &v1alpha1.StoreEndpointList{}
	err = r.Client.List(context.TODO(), storeEndpoints)
	if err != nil {
		// Object not found, return.  Created objects are automatically garbage collected.
		// For additional cleanup logic use finalizers.
		if apierrors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	// Create resource factory
	// Create reconciler for objects
	thanosComponentReconciler := resources.NewThanosComponentReconciler(
		thanos,
		storeEndpoints,
		reconciler.NewReconciler(r.Client, r.Log, reconciler.ReconcilerOpts{}))
	reconcilers := make([]resources.ComponentReconciler, 0)

	// Query
	reconcilers = append(reconcilers, query.New(thanos, thanosComponentReconciler).Reconcile)
	// Store
	reconcilers = append(reconcilers, store.New(thanos, thanosComponentReconciler).Reconcile)
	// Rule
	reconcilers = append(reconcilers, rule.New(thanos, thanosComponentReconciler).Reconcile)

	return resources.RunReconcilers(reconcilers)
}

func (r *ThanosReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.Thanos{}).
		Complete(r)
}
