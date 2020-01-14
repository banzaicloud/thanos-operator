// Copyright © 2020 Banzai Cloud
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package compactor

import (
	"github.com/banzaicloud/operator-tools/pkg/reconciler"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func (c *Compactor) persistentVolumeClaim() (runtime.Object, reconciler.DesiredState, error) {
	const app = "compactor"
	name := app + "-service"

	if c.objectStore.Spec.BucketWeb.Enabled {
		compactor := c.objectStore.Spec.Compactor.DeepCopy()

		return &corev1.PersistentVolumeClaim{
			ObjectMeta: c.objectMeta(name, &compactor.BaseObject),
			Spec:       corev1.PersistentVolumeClaimSpec{
				//Selector: c.labels(),
			},
		}, reconciler.StatePresent, nil
	}

	return &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: c.namespace,
		},
	}, reconciler.StateAbsent, nil
}
