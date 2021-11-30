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
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"

	"github.com/banzaicloud/thanos-operator/pkg/resources/receiver"
)

func NewReceiverReconciler(client client.Client, logger logr.Logger) *ReceiverReconciler {
	return &ReceiverReconciler{
		NativeReconciler: reconciler.NewNativeReconciler(
			monitoringv1alpha1.ReceiverName,
			reconciler.NewGenericReconciler(client, logger, reconciler.ReconcilerOpts{}),
			client,
			receiver.NewComponent(client.Scheme()),
			func(o runtime.Object) (parent reconciler.ResourceOwner, config interface{}) {
				receiver := o.(*monitoringv1alpha1.Receiver)
				return receiver, nil
			},
			reconciler.NativeReconcilerWithScheme(client.Scheme()),
		),
	}
}

// ReceiverReconciler reconciles a Receiver object
type ReceiverReconciler struct {
	*reconciler.NativeReconciler
}

// +kubebuilder:rbac:groups=monitoring.banzaicloud.io,resources=receivers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=monitoring.banzaicloud.io,resources=receivers/status,verbs=get;update;patch

func (r *ReceiverReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	result := ctrl.Result{}

	receiver := &monitoringv1alpha1.Receiver{}
	if err := r.Client.Get(ctx, req.NamespacedName, receiver); err != nil {
		return result, client.IgnoreNotFound(err)
	}

	res, err := r.NativeReconciler.Reconcile(receiver)
	if res != nil {
		result = *res
	}
	return result, err
}

func (r *ReceiverReconciler) SetupWithManager(mgr ctrl.Manager) (controller.Controller, error) {
	b := ctrl.NewControllerManagedBy(mgr)
	r.NativeReconciler.RegisterWatches(b)
	return b.Build(r)
}
