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
	monitoringv1alpha1 "github.com/banzaicloud/thanos-operator/pkg/sdk/api/v1alpha1"
	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"

	"github.com/banzaicloud/thanos-operator/pkg/resources"
	"github.com/banzaicloud/thanos-operator/pkg/resources/receiver"
)

// ReceiverReconciler reconciles a Receiver object
type ReceiverReconciler struct {
	Client client.Client
	Log    logr.Logger
}

// +kubebuilder:rbac:groups=monitoring.banzaicloud.io,resources=receivers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=monitoring.banzaicloud.io,resources=receivers/status,verbs=get;update;patch

func (r *ReceiverReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	result := ctrl.Result{}
	log := r.Log.WithValues("receivers", req.NamespacedName)

	rec := &monitoringv1alpha1.Receiver{}
	if err := r.Client.Get(ctx, req.NamespacedName, rec); err != nil {
		return result, client.IgnoreNotFound(err)
	}
	receiverReconciler := receiver.NewReconciler(rec, reconciler.NewReconcilerWith(r.Client, reconciler.WithLog(log)))

	reconcilers := []resources.ComponentReconciler{
		receiverReconciler.Reconcile,
	}

	return resources.RunReconcilers(reconcilers)
}

func (r *ReceiverReconciler) SetupWithManager(mgr ctrl.Manager) (controller.Controller, error) {
	return ctrl.NewControllerManagedBy(mgr).
		For(&monitoringv1alpha1.Receiver{}).
		Owns(&corev1.Service{}).
		Owns(&netv1.Ingress{}).
		Owns(&appsv1.StatefulSet{}).
		Owns(&corev1.ConfigMap{}).
		Build(r)
}
