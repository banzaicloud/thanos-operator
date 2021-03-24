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

package v1alpha1

import (
	"github.com/banzaicloud/operator-tools/pkg/typeoverride"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Endpoint Address",type="string",JSONPath=".status.endpointAddress"

type ThanosEndpoint struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ThanosEndpointSpec   `json:"spec,omitempty"`
	Status ThanosEndpointStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

type ThanosEndpointList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ThanosEndpoint `json:"items"`
}

type ThanosEndpointSpec struct {
	// The endpoint should use this server certificate (tls.crt, tls.key) in the current namespace
	Certificate string `json:"certificate,omitempty"`

	// Name of the secret that contains the CA certificate in ca.crt to verify client certs in the current namespace
	CABundle string `json:"caBundle,omitempty"`

	// List of statically configured store addresses
	Stores []string `json:"stores,omitempty"`

	// Custom replica labels if the default doesn't apply
	ReplicaLabels []string `json:"replicaLabels,omitempty"`

	// Override metadata for managed resources
	MetaOverrides typeoverride.ObjectMeta `json:"metaOverrides,omitempty"`

	// Override any of the Query parameters
	QueryOverrides *Query `json:"queryOverrides,omitempty"`

	// Override any of the StoreEndpoint parameters
	StoreEndpointOverrides []StoreEndpointSpec `json:"storeEndpointOverrides,omitempty"`
}

type ThanosEndpointStatus struct {
	// Host (or IP) and port of the exposed Thanos endpoint
	EndpointAddress string `json:"endpointAddress,omitempty"`
}

func init() {
	SchemeBuilder.Register(&ThanosEndpoint{}, &ThanosEndpointList{})
}
