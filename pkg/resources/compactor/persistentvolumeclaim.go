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

package compactor

import (
	"github.com/banzaicloud/operator-tools/pkg/reconciler"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func (c *Compactor) persistentVolumeClaim() (runtime.Object, reconciler.DesiredState, error) {
	if c.ObjectStore.Spec.Compactor != nil &&
		c.ObjectStore.Spec.Compactor.DataVolume != nil &&
		c.ObjectStore.Spec.Compactor.DataVolume.PersistentVolumeClaim != nil &&
		c.ObjectStore.Spec.Compactor.DataVolume.PersistentVolumeClaim.PersistentVolumeSource.ClaimName == "" {
		compactor := c.ObjectStore.Spec.Compactor.DeepCopy()
		pvc := &corev1.PersistentVolumeClaim{
			ObjectMeta: compactor.DeploymentOverrides.Merge(c.getMeta()),
			Spec:       c.ObjectStore.Spec.Compactor.DataVolume.PersistentVolumeClaim.PersistentVolumeClaimSpec,
		}
		return pvc, reconciler.StatePresent, nil
	}

	pvc := &corev1.PersistentVolumeClaim{
		ObjectMeta: c.getMeta(),
	}
	return pvc, reconciler.StateAbsent, nil
}
