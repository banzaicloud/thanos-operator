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

package bucketweb

import (
	"github.com/banzaicloud/operator-tools/pkg/reconciler"
	policyv1beta1 "k8s.io/api/policy/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func (b *BucetWeb) podDistributionBucket() (runtime.Object, reconciler.DesiredState, error) {
	const app = "bucketweb"
	name := app + "-poddisruptionbudget"

	if b.objectStore.Spec.BucketWeb.Enabled {
		bucketWeb := b.objectStore.Spec.BucketWeb.DeepCopy()

		return &policyv1beta1.PodDisruptionBudget{
			ObjectMeta: b.objectMeta(name, &bucketWeb.BaseObject),
			Spec: policyv1beta1.PodDisruptionBudgetSpec{
				MinAvailable: &intstr.IntOrString{IntVal: 1},
				Selector:     &metav1.LabelSelector{},
			},
		}, reconciler.StatePresent, nil
	}

	return &policyv1beta1.PodDisruptionBudget{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: b.namespace,
		},
	}, reconciler.StateAbsent, nil
}
