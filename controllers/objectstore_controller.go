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
	"github.com/banzaicloud/thanos-operator/pkg/resources/bucketweb"
	"github.com/banzaicloud/thanos-operator/pkg/resources/compactor"
	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"

	monitoringv1alpha1 "github.com/banzaicloud/thanos-operator/pkg/sdk/api/v1alpha1"
)

// ObjectStoreReconciler reconciles a ObjectStore object
type ObjectStoreReconciler struct {
	Client client.Client
	Log    logr.Logger
}

// +kubebuilder:rbac:groups=monitoring.banzaicloud.io,resources=objectstores,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=monitoring.banzaicloud.io,resources=objectstores/status,verbs=get;update;patch

func (r *ObjectStoreReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	result := ctrl.Result{}
	log := r.Log.WithValues("objectstore", req.NamespacedName)

	store := &monitoringv1alpha1.ObjectStore{}
	err := r.Client.Get(ctx, req.NamespacedName, store)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return result, nil
		}
		return result, err
	}
	objectStoreReconciler := resources.NewObjectStoreReconciler(store, reconciler.NewGenericReconciler(r.Client, log, reconciler.ReconcilerOpts{}))

	reconcilers := make([]resources.ComponentReconciler, 0)

	// Bucket Web
	reconcilers = append(reconcilers, bucketweb.New(objectStoreReconciler).Reconcile)
	// Compactor
	reconcilers = append(reconcilers, compactor.New(objectStoreReconciler).Reconcile)

	return resources.RunReconcilers(reconcilers)
}

func (r *ObjectStoreReconciler) SetupWithManager(mgr ctrl.Manager) (controller.Controller, error) {
	return ctrl.NewControllerManagedBy(mgr).
		For(&monitoringv1alpha1.ObjectStore{}).
		Owns(&appsv1.Deployment{}).
		Owns(&corev1.Service{}).
		Owns(&netv1.Ingress{}).
		Owns(&corev1.PersistentVolumeClaim{}).
		Build(r)
}
