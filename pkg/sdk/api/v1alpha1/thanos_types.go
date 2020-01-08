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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ThanosSpec defines the desired state of Thanos
type ThanosSpec struct {
	Remote          *Remote          `json:"remote,omitempty"`
	ThanosDiscovery *ThanosDiscovery `json:"thanosDiscovery,omitempty"`
	Local           *Local           `json:"local,omitempty"`
	// StoreGateway
	StoreGateway *StoreGateway `json:"storeGateway,omitempty"`
	Rule         *Rule         `json:"rule,omitempty"`
}

type Remote struct {
}

type Local struct {
}

type ThanosDiscovery struct {
}

type StoreGateway struct {
}

type Rule struct {
}

// ThanosStatus defines the observed state of Thanos
type ThanosStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true

// Thanos is the Schema for the thanos API
type Thanos struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ThanosSpec   `json:"spec,omitempty"`
	Status ThanosStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ThanosList contains a list of Thanos
type ThanosList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Thanos `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Thanos{}, &ThanosList{})
}
