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

// Package v1alpha1 contains API Schema definitions for the monitoring v1alpha1 API group
// +kubebuilder:object:generate=true
// +groupName=monitoring.banzaicloud.io
package v1alpha1

import corev1 "k8s.io/api/core/v1"

type BaseObject struct {
	Annotations  map[string]string           `json:"annotations,omitempty"`
	Labels       map[string]string           `json:"labels,omitempty"`
	Resources    corev1.ResourceRequirements `json:"resources,omitempty"`
	Tolerations  []corev1.Toleration         `json:"tolerations,omitempty"`
	NodeSelector map[string]string           `json:"nodeSelector,omitempty"`
	Image        ImageSpec                   `json:"image,omitempty"`
}

type ImageSpec struct {
	Repository string `json:"repository,omitempty"`
	Tag        string `json:"tag,omitempty"`
	PullPolicy string `json:"pullPolicy,omitempty"`
}
