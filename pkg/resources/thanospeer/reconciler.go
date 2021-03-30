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
	"fmt"

	"github.com/banzaicloud/operator-tools/pkg/reconciler"
	"github.com/banzaicloud/operator-tools/pkg/utils"
	"github.com/banzaicloud/thanos-operator/pkg/resources"
	"github.com/banzaicloud/thanos-operator/pkg/sdk/api/v1alpha1"
	"github.com/go-logr/logr"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type Reconciler struct {
	log logr.Logger
	resourceReconciler reconciler.ResourceReconciler
	client client.Client
	peer *v1alpha1.ThanosPeer
}

func NewReconciler(logger logr.Logger, client client.Client, reconciler reconciler.ResourceReconciler, peer *v1alpha1.ThanosPeer) *Reconciler {
	return &Reconciler{
		log: logger,
		resourceReconciler: reconciler,
		peer: peer,
		client: client,
	}
}

func (r Reconciler) Reconcile() (*reconcile.Result, error) {
	var resourceList []resources.Resource

	//endpointDetails := []interface{}{
	//	"name", r.peer.Name,
	//	"namespace", r.peer.Namespace,
	//}



	resourceList = append(resourceList, r.query)
	return resources.Dispatch(r.resourceReconciler, resourceList)
}


func (r Reconciler) getDescendantMeta() metav1.ObjectMeta {
	meta := metav1.ObjectMeta{
		Name: r.getDescendantResourceName(),
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
		resources.NameLabel:     v1alpha1.PeerName,
		resources.InstanceLabel:  r.peer.Name,
		resources.ManagedByLabel: resources.ManagedByValue,
	}
}
