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
	"github.com/banzaicloud/logging-operator/pkg/sdk/model/secret"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ThanosSpec defines the desired state of Thanos
type ThanosSpec struct {
	Remote          *Remote          `json:"remote,omitempty"`
	ThanosDiscovery *ThanosDiscovery `json:"thanosDiscovery,omitempty"`
	Local           *Local           `json:"local,omitempty"`
	StoreGateway    *StoreGateway    `json:"storeGateway,omitempty"`
	Rule            *Rule            `json:"rule,omitempty"`
	ObjectStore     string           `json:"object_store,omitempty"`
}

type TLS struct {
	Managed     ManagedTLS    `json:"managedTLS,omitempty"`
	Certificate secret.Secret `json:"certificate,omitempty"`
}

// TODO how the runtime generated certificate will work
type ManagedTLS struct {
}

type Remote struct {
	URLs []string `json:"urls"`
	TLS  TLS      `json:"tls"`
}

type Local struct {
	URLs []string `json:"urls"`
	TLS  TLS      `json:"tls"`
}

type ThanosDiscovery struct {
	LabelSelector metav1.LabelSelector `json:"labelSelector"`
}

type TimeRange struct {
	// Start of time range limit to serve. Thanos Store will serve only metrics, which happened
	// later than this value. Option can be a constant time in RFC3339 format or time duration
	// relative to current time, such as -1d or 2h45m. Valid duration units are ms, s, m, h, d, w, y.
	MinTime string `json:"minTime,omitempty"`
	// 	End of time range limit to serve. Thanos Store
	//	will serve only blocks, which happened eariler
	//	than this value. Option can be a constant time
	//	in RFC3339 format or time duration relative to
	//	current time, such as -1d or 2h45m. Valid
	//	duration units are ms, s, m, h, d, w, y.
	MaxTime string `json:"maxTime,omitempty"`
}

type StoreGateway struct {
	BaseObject `json:",inline"`
	LogLevel   string `json:"logLevel,omitempty"`
	LogFormat  string `json:"logFormat,omitempty"`
	// Listen host:port for HTTP endpoints.
	HttpAddress string `json:"httpAddress"`
	// Time to wait after an interrupt received for HTTP Server.
	HttpGracePeriod string `json:"http_grace_period"`
	// Listen ip:port address for gRPC endpoints
	GRPCAddress string `json:"grpcAddress"`
	// Time to wait after an interrupt received for GRPC Server.
	GRPCGracePeriod string `json:"grpcGracePeriod"`
	// Maximum size of items held in the in-memory index cache.
	IndexCacheSize string `json:"indexCacheSize"`
	// Maximum size of concurrently allocatable bytes for chunks.
	ChunkPoolSize string `json:"chunkPoolSize,omitempty"`
	// Maximum amount of samples returned via a single Series call. 0 means no limit. NOTE: For
	// efficiency we take 120 as the number of samples in chunk (it cannot be bigger than that), so
	// the actual number of samples might be lower, even though the maximum could be hit.
	StoreGRPCSeriesSampleLimit string `json:"storeGRPCSeriesSampleLimit,omitempty"`
	// Maximum number of concurrent Series calls.
	StoreGRPCSeriesMaxConcurrency int `json:"storeGRPCSeriesMaxConcurrency,omitempty"`
	// Repeat interval for syncing the blocks between local and remote view.
	SyncBlockDuration string `json:"syncBlockDuration,omitempty"`
	// Number of goroutines to use when constructing index-cache.json blocks from object storage.
	BlockSyncConcurrency int `json:"blockSyncConcurrency,omitempty"`
	// TimeRanges is a list of TimeRange to partition Store Gateway
	TimeRanges []TimeRange `json:"timeRanges,omitempty"`
}

type Rule struct {
	BaseObject `json:",inline"`
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
