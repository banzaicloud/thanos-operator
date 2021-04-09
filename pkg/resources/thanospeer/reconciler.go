// Copyright 2021 Banzai Cloud
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

package thanospeer

import (
	"context"
	"fmt"

	"emperror.dev/errors"
	"github.com/banzaicloud/operator-tools/pkg/reconciler"
	"github.com/banzaicloud/operator-tools/pkg/utils"
	"github.com/banzaicloud/thanos-operator/pkg/resources"
	"github.com/banzaicloud/thanos-operator/pkg/sdk/api/v1alpha1"
	"github.com/go-logr/logr"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type Reconciler struct {
	log                logr.Logger
	resourceReconciler reconciler.ResourceReconciler
	client             client.Client
	peer               *v1alpha1.ThanosPeer
}

func NewReconciler(logger logr.Logger, client client.Client, reconciler reconciler.ResourceReconciler, peer *v1alpha1.ThanosPeer) *Reconciler {
	return &Reconciler{
		log:                logger,
		resourceReconciler: reconciler,
		peer:               peer,
		client:             client,
	}
}

func (r Reconciler) Reconcile() (*reconcile.Result, error) {
	var resourceList []resources.Resource

	peerDetails := []interface{}{
		"name", r.peer.Name,
		"namespace", r.peer.Namespace,
	}

	r.log.V(0).Info("started reconcile", peerDetails...)

	var peerCert, peerCA string
	peerCerts := &v1.SecretList{}

	err := r.client.List(context.TODO(), peerCerts, client.MatchingLabels{
		v1alpha1.PeerCertSecretLabel: r.peer.Name,
	}, client.InNamespace(r.peer.Namespace))
	if err != nil {
		return nil, errors.Wrap(err, "failed to list ThanosPeer certificate secrets")
	}

	switch len(peerCerts.Items) {
	case 0:
	case 1:
		peerCert = peerCerts.Items[0].Name
		r.log.V(0).Info("peer cert available", append(peerDetails, "cert", peerCert)...)
	default:
		return nil, errors.NewWithDetails("more than one certs available, expecting only one", peerDetails)
	}

	peerCAs := &v1.SecretList{}

	err = r.client.List(context.TODO(), peerCAs, client.MatchingLabels{
		v1alpha1.PeerCASecretLabel: r.peer.Name,
	}, client.InNamespace(r.peer.Namespace))
	if err != nil {
		return nil, errors.Wrap(err, "failed to list ThanosPeer CA secrets")
	}

	switch len(peerCAs.Items) {
	case 0:
	case 1:
		peerCA = peerCAs.Items[0].Name
		r.log.V(0).Info("peer ca available", append(peerDetails, "ca", peerCA)...)
	default:
		return nil, errors.NewWithDetails("more than one CAs available, expecting only one", peerDetails)
	}

	resourceList = append(resourceList, r.query(peerCert, peerCA))
	return resources.Dispatch(r.resourceReconciler, resourceList)
}

func (r Reconciler) getDescendantMeta() metav1.ObjectMeta {
	meta := metav1.ObjectMeta{
		Name:      r.getDescendantResourceName(),
		Namespace: r.peer.Namespace,
	}
	meta.OwnerReferences = []metav1.OwnerReference{
		{
			APIVersion: r.peer.APIVersion,
			Kind:       r.peer.Kind,
			Name:       r.peer.Name,
			UID:        r.peer.UID,
			Controller: utils.BoolPointer(true),
		},
	}
	meta.Labels = r.getLabels()
	return meta
}

func (r Reconciler) getDescendantResourceName() string {
	name := r.qualifiedName(v1alpha1.PeerName)
	return name
}

func (r Reconciler) qualifiedName(name string) string {
	return fmt.Sprintf("%s-%s", r.peer.Name, name)
}

func (r Reconciler) getLabels() resources.Labels {
	return resources.Labels{
		resources.NameLabel:      v1alpha1.PeerName,
		resources.InstanceLabel:  r.peer.Name,
		resources.ManagedByLabel: resources.ManagedByValue,
	}
}
