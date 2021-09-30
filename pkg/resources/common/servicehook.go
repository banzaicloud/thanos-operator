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

package common

import (
	"emperror.dev/errors"
	"github.com/banzaicloud/operator-tools/pkg/reconciler"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func ServiceUpdateHook(service *corev1.Service) reconciler.DesiredStateHook {
	return func(current runtime.Object) error {
		if s, ok := current.(*corev1.Service); ok {
			service.Spec.ClusterIP = s.Spec.ClusterIP
		} else {
			return errors.Errorf("failed to cast service object %+v", current)
		}
		return nil
	}
}
