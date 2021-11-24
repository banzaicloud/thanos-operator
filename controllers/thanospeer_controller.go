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
	"strings"

	"github.com/banzaicloud/operator-tools/pkg/reconciler"
	"github.com/banzaicloud/thanos-operator/pkg/resources"
	"github.com/banzaicloud/thanos-operator/pkg/resources/thanospeer"
	"github.com/go-logr/logr"
	v1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	monitoringv1alpha1 "github.com/banzaicloud/thanos-operator/pkg/sdk/api/v1alpha1"
)

// ThanosPeerReconciler reconciles a ThanosPeer object
type ThanosPeerReconciler struct {
	Client client.Client
	Log    logr.Logger
}

// +kubebuilder:rbac:groups=monitoring.banzaicloud.io,resources=thanospeers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=monitoring.banzaicloud.io,resources=thanospeers/status,verbs=get;update;patch
// +kubebuilder:rbac:groups="",resources=secrets,verbs=list;get;watch
// +kubebuilder:rbac:groups=integreatly.org,resources=grafanadatasources,verbs=get;list;watch;create;update;patch;delete

func (r *ThanosPeerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	result := ctrl.Result{}
	log := r.Log.WithValues("thanospeers", req.NamespacedName)

	peer := &monitoringv1alpha1.ThanosPeer{}
	err := r.Client.Get(ctx, req.NamespacedName, peer)
	if err != nil {
		if apierrors.IsNotFound(err) {
			log.V(0).Info("notfound reconcile")
			return result, nil
		}
		return result, err
	}

	rec := thanospeer.NewReconciler(log, r.Client, reconciler.NewReconcilerWith(r.Client), peer)

	reconcilers := []resources.ComponentReconciler{
		rec.Reconcile,
	}

	return resources.RunReconcilers(reconcilers)
}

func (r *ThanosPeerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&monitoringv1alpha1.ThanosPeer{}).
		Owns(&monitoringv1alpha1.Thanos{}).
		Watches(&source.Kind{Type: &v1.Secret{}}, handler.EnqueueRequestsFromMapFunc(
			func(object client.Object) []reconcile.Request {
				secret := object.(*v1.Secret)
				if secret.Labels != nil {
					for _, i := range []string{monitoringv1alpha1.PeerCertSecretLabel, monitoringv1alpha1.PeerCASecretLabel} {
						if peer := secret.Labels[i]; peer != "" {
							return []reconcile.Request{
								{
									NamespacedName: types.NamespacedName{
										Name:      peer,
										Namespace: secret.Namespace,
									},
								},
							}
						}
					}
				}
				return nil
			})).
		Watches(&source.Kind{Type: &v1.Service{}}, handler.EnqueueRequestsFromMapFunc(
			func(object client.Object) []reconcile.Request {
				service := object.(*v1.Service)
				const peerNameSuffix = "-" + monitoringv1alpha1.PeerName
				if service != nil {
					managedBy, name := service.Labels[resources.ManagedByLabel], service.Labels[resources.NameLabel]
					if name == monitoringv1alpha1.QueryName && strings.HasSuffix(managedBy, peerNameSuffix) {
						return []reconcile.Request{
							{
								NamespacedName: types.NamespacedName{
									Name:      strings.TrimSuffix(managedBy, peerNameSuffix),
									Namespace: service.Namespace,
								},
							},
						}
					}
				}
				return nil
			})).
		Complete(r)
}
