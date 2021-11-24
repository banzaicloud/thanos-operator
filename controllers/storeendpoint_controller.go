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
	"github.com/banzaicloud/thanos-operator/pkg/resources/sidecar"
	"github.com/banzaicloud/thanos-operator/pkg/sdk/api/v1alpha1"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// StoreEndpointReconciler reconciles a StoreEndpoint object
type StoreEndpointReconciler struct {
	Client client.Client
	Log    logr.Logger
}

// +kubebuilder:rbac:groups=monitoring.banzaicloud.io,resources=storeendpoints,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=monitoring.banzaicloud.io,resources=storeendpoints/status,verbs=get;update;patch

func (r *StoreEndpointReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	result := ctrl.Result{}
	log := r.Log.WithValues("storeendpoint", req.NamespacedName)

	endpoints := &v1alpha1.StoreEndpointList{}
	err := r.Client.List(ctx, endpoints)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return result, nil
		}
		return result, err
	}
	storeEndpointReconciler := resources.NewStoreEndpointComponentReconciler(endpoints, reconciler.NewGenericReconciler(r.Client, log, reconciler.ReconcilerOpts{}))

	reconcilers := make([]resources.ComponentReconciler, 0)

	// Bucket Web
	reconcilers = append(reconcilers, sidecar.New(endpoints, storeEndpointReconciler).Reconcile)

	return resources.RunReconcilers(reconcilers)
}

func (r *StoreEndpointReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.StoreEndpoint{}).
		Owns(&corev1.Service{}).
		Owns(&netv1.Ingress{}).
		Complete(r)
}
