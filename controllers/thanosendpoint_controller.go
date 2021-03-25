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

	"github.com/banzaicloud/operator-tools/pkg/reconciler"
	"github.com/banzaicloud/thanos-operator/pkg/resources"
	"github.com/banzaicloud/thanos-operator/pkg/resources/thanosendpoint"
	"github.com/go-logr/logr"
	"k8s.io/api/networking/v1beta1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	monitoringv1alpha1 "github.com/banzaicloud/thanos-operator/pkg/sdk/api/v1alpha1"
)

// ThanosEndpointReconciler reconciles a ThanosEndpoint object
type ThanosEndpointReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=monitoring.banzaicloud.io,resources=thanosendpoints,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=monitoring.banzaicloud.io,resources=thanosendpoints/status,verbs=get;update;patch

func (r *ThanosEndpointReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	result := ctrl.Result{}
	ctx := context.Background()
	log := r.Log.WithValues("thanosendpoints", req.NamespacedName)

	endpoint := &monitoringv1alpha1.ThanosEndpoint{}
	err := r.Client.Get(ctx, req.NamespacedName, endpoint)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return result, nil
		}
		return result, err
	}

	rec := thanosendpoint.NewReconciler(log, r.Client, reconciler.NewReconcilerWith(r), endpoint)

	reconcilers := []resources.ComponentReconciler{
		rec.Reconcile,
	}

	return resources.RunReconcilers(reconcilers)
}

func (r *ThanosEndpointReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&monitoringv1alpha1.ThanosEndpoint{}).
		Owns(&monitoringv1alpha1.Thanos{}).
		Owns(&monitoringv1alpha1.StoreEndpoint{}).
		Watches(&source.Kind{Type: &v1beta1.Ingress{}}, &handler.EnqueueRequestsFromMapFunc{
			ToRequests: handler.ToRequestsFunc(func(object handler.MapObject) []reconcile.Request {
				ing := object.Object.(*v1beta1.Ingress)
				if ing.Labels != nil {
					if mb := ing.Labels["managed-by"]; mb != "" {
						return []reconcile.Request{
							{
								NamespacedName: types.NamespacedName{
									Name: mb,
									Namespace: ing.Namespace,
								},
							},
						}
					}
				}
				return nil
			}),
	}).
		Complete(r)
}
