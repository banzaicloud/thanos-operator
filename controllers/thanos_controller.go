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

package controllers

import (
	"context"
	"fmt"

	"github.com/banzaicloud/operator-tools/pkg/reconciler"
	"github.com/banzaicloud/thanos-operator/pkg/resources"
	"github.com/banzaicloud/thanos-operator/pkg/resources/query"
	"github.com/banzaicloud/thanos-operator/pkg/resources/query_frontend"
	"github.com/banzaicloud/thanos-operator/pkg/resources/rule"
	"github.com/banzaicloud/thanos-operator/pkg/resources/store"
	"github.com/banzaicloud/thanos-operator/pkg/sdk/api/v1alpha1"
	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

// ThanosReconciler reconciles a Thanos object
type ThanosReconciler struct {
	Client client.Client
	Log    logr.Logger
}

// +kubebuilder:rbac:groups=monitoring.banzaicloud.io,resources=thanos,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=monitoring.banzaicloud.io,resources=thanos/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=coordination.k8s.io,resources=leases,verbs=*

func (r *ThanosReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("thanos", req.NamespacedName)

	thanos := &v1alpha1.Thanos{}
	err := r.Client.Get(ctx, req.NamespacedName, thanos)
	if err != nil {
		// Object not found, return.  Created objects are automatically garbage collected.
		// For additional cleanup logic use finalizers.
		if apierrors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}
	thanosList := &v1alpha1.ThanosList{}
	err = r.Client.List(context.TODO(), thanosList)
	if err != nil {
		// Object not found, return.  Created objects are automatically garbage collected.
		// For additional cleanup logic use finalizers.
		if apierrors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}
	var thanosInstances []v1alpha1.Thanos
	for _, t := range thanosList.Items {
		if t.Name != thanos.Name {
			thanosInstances = append(thanosInstances, t)
		}
	}
	// Collect StoreEndpoints for matching Thanos CR
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
	var storeEndpointList []v1alpha1.StoreEndpoint
	for _, s := range storeEndpoints.Items {
		if s.Spec.Thanos == thanos.Name {
			storeEndpointList = append(storeEndpointList, s)
		}
	}
	reconcilerOpts := reconciler.ReconcilerOpts{
		EnableRecreateWorkloadOnImmutableFieldChange: thanos.Spec.EnableRecreateWorkloadOnImmutableFieldChange,
		EnableRecreateWorkloadOnImmutableFieldChangeHelp: "Object has to be recreated, but refusing to remove without explicitly being told so. " +
			"Use thanos.spec.enableRecreateWorkloadOnImmutableFieldChange to move on but make sure to understand the consequences. " +
			"As of rule, to avoid data loss, make sure to use a persistent volume for buffers, which is the default, unless explicitly disabled or configured differently.",
	}
	// Create reconciler for objects
	thanosComponentReconciler := resources.NewThanosComponentReconciler(
		thanos,
		thanosInstances,
		storeEndpointList,
		reconciler.NewGenericReconciler(r.Client, log, reconcilerOpts))
	reconcilers := make([]resources.ComponentReconciler, 0)

	// Query
	reconcilers = append(reconcilers, query.New(thanosComponentReconciler).Reconcile)
	// Store
	reconcilers = append(reconcilers, store.New(thanosComponentReconciler).Reconcile)
	// Rule
	reconcilers = append(reconcilers, rule.New(thanosComponentReconciler).Reconcile)
	// QueryFrontend
	reconcilers = append(reconcilers, query_frontend.New(thanosComponentReconciler).Reconcile)

	return resources.RunReconcilers(reconcilers)
}

func (r *ThanosReconciler) SetupWithManager(mgr ctrl.Manager) (controller.Controller, error) {
	requestMapper := handler.EnqueueRequestsFromMapFunc(func(object client.Object) []reconcile.Request {
		if o, ok := object.(*v1alpha1.StoreEndpoint); ok {
			thanos := &v1alpha1.Thanos{}
			if err := mgr.GetClient().Get(context.TODO(), types.NamespacedName{Name: o.Spec.Thanos, Namespace: o.Namespace}, thanos); err != nil {
				r.Log.Error(err, fmt.Sprintf("failed to get thanos resources %q for endpoint %q", o.Spec.Thanos, o.Name))
				return nil
			}
			return []reconcile.Request{
				{
					NamespacedName: types.NamespacedName{
						Namespace: thanos.Namespace,
						Name:      thanos.Name,
					},
				},
			}
		}
		return nil
	})

	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.Thanos{}).
		Watches(&source.Kind{Type: &v1alpha1.StoreEndpoint{}}, requestMapper).
		Owns(&appsv1.Deployment{}).
		Owns(&appsv1.StatefulSet{}).
		Owns(&corev1.Service{}).
		Owns(&netv1.Ingress{}).
		Owns(&appsv1.StatefulSet{}).
		Build(r)
}
