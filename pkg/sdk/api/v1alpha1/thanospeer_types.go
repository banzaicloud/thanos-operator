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

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

type ThanosPeer struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ThanosPeerSpec   `json:"spec,omitempty"`
	Status ThanosPeerStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

type ThanosPeerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Thanos `json:"items"`
}

type ThanosPeerSpec struct {
	// Host (or IP) and port of the remote Thanos endpoint
	EndpointAddress string `json:"endpointAddress"`

	// Optional alias for the remote endpoint in case we have to access it through a different name.
	// This is typically needed if the remote endpoint has a certificate created for a predefined hostname.
	// The controller should create an externalName service for this backed buy the actual peer endpoint host
	// or a k8s service with a manually crafted k8s endpoint if EndpointAddress doesn't have a host but only an IP.
	PeerEndpointAlias string `json:"peerEndpointAlias,omitempty"`

	// CA certificate to verify the server cert
	CABundle string `json:"caBundle,omitempty"`

	// Custom replica labels if the default doesn't apply
	ReplicaLabels []string `json:"replicaLabels,omitempty"`
}

type ThanosPeerStatus struct {

}

func init() {
	SchemeBuilder.Register(&ThanosPeer{}, &ThanosPeerList{})
}
