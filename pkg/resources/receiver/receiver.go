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

package receiver

import (
	"github.com/banzaicloud/operator-tools/pkg/utils"
	"github.com/banzaicloud/thanos-operator/pkg/resources"
	"github.com/banzaicloud/thanos-operator/pkg/sdk/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func extend(r *v1alpha1.Receiver) receiverExt {
	return receiverExt{
		Receiver: r,
	}
}

type receiverExt struct {
	*v1alpha1.Receiver
}

func (r receiverExt) getLabels() resources.Labels {
	return resources.Labels{
		resources.ManagedByLabel: r.Name,
		resources.NameLabel:      v1alpha1.ReceiverName,
	}
}

func (r receiverExt) getMeta(suffix ...string) metav1.ObjectMeta {
	meta := r.GetObjectMeta(r.getName(suffix...))
	return meta
}

func (r receiverExt) getMetaWithLabels() metav1.ObjectMeta {
	meta := r.getMeta()
	meta.Labels = r.getLabels()
	return meta
}

func (r receiverExt) getName(suffix ...string) string {
	return resources.QualifiedName(append([]string{r.Name, v1alpha1.ReceiverName}, suffix...)...)
}

func (r receiverExt) GetObjectMeta(name string) metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Name:      name,
		Namespace: r.Namespace,
		OwnerReferences: []metav1.OwnerReference{
			{
				APIVersion: r.APIVersion,
				Kind:       r.Kind,
				Name:       r.Name,
				UID:        r.UID,
				Controller: utils.BoolPointer(true),
			},
		},
	}
}

//func (r receiverExt) GetServiceURLS() []string {
//	var urls []string
//	for _, endpoint := range r.StoreEndpoints {
//		urls = append(urls, (&receiverInstance{r, endpoint.DeepCopy()}).getSvc())
//	}
//	return urls
//}

type receiverInstance struct {
	receiverExt
	receiverGroup *v1alpha1.ReceiverGroup
}

func (r *receiverInstance) getVolumeMeta(name string) metav1.ObjectMeta {
	meta := r.GetObjectMeta(name)
	meta.Labels = r.getLabels()
	return meta
}

func (r *receiverInstance) getMeta(suffix ...string) metav1.ObjectMeta {
	meta := r.receiverExt.getMeta(suffix...)
	meta.Labels = r.getLabels()
	return meta
}

func (r *receiverInstance) getLabels() resources.Labels {
	labels := r.receiverExt.getLabels()
	if r.receiverGroup != nil {
		labels["receiverGroup"] = r.receiverGroup.Name
	}
	return labels
}

func (r *receiverInstance) setArgs(args []string) []string {
	args = append(args, resources.GetArgs(r.receiverGroup)...)

	//Label

	// Local-endpoint

	return args
}
