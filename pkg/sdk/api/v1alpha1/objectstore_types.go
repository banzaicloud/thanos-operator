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

package v1alpha1

import (
	"time"

	"github.com/banzaicloud/operator-tools/pkg/secret"
	"github.com/banzaicloud/operator-tools/pkg/volume"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ObjectStoreSpec defines the desired state of ObjectStore
type ObjectStoreSpec struct {
	// Config
	Config    secret.Secret `json:"config"`
	Compactor *Compactor    `json:"compactor,omitempty"`
	BucketWeb *BucketWeb    `json:"bucketWeb,omitempty"`
}

var DefaultCompactor = &Compactor{
	BaseObject: BaseObject{
		Image: ImageSpec{
			Repository: thanosImageRepository,
			Tag:        thanosImageTag,
			PullPolicy: defaultPullPolicy,
		},
	},
	Metrics: &Metrics{
		Interval:       "15s",
		Timeout:        "5s",
		Path:           "/metrics",
		ServiceMonitor: false,
	},
	HTTPAddress:            "0.0.0.0:10902",
	HTTPGracePeriod:        metav1.Duration{Duration: 2 * time.Minute},
	DataDir:                "/data",
	ConsistencyDelay:       metav1.Duration{Duration: 30 * time.Minute},
	RetentionResolutionRaw: metav1.Duration{Duration: 0},
	RetentionResolution5m:  metav1.Duration{Duration: 0},
	RetentionResolution1h:  metav1.Duration{Duration: 0},
	BlockSyncConcurrency:   20,
	CompactConcurrency:     1,
	Wait:                   true,
}

type Compactor struct {
	BaseObject `json:",inline"`
	Metrics    *Metrics `json:"metrics,omitempty"`
	// Listen host:port for HTTP endpoints.
	HTTPAddress string `json:"httpAddress,omitempty"`
	// Time to wait after an interrupt received for HTTP Server.
	HTTPGracePeriod metav1.Duration `json:"httpGracePeriod,omitempty"`
	// Data directory in which to cache blocks and process compactions.
	DataDir string `json:"dataDir,omitempty"`
	// Kubernetes volume abstraction refers to different types of volumes to be mounted to pods: emptyDir, hostPath, pvc.
	DataVolume *volume.KubernetesVolume `json:"dataVolume,omitempty"`
	// Minimum age of fresh (non-compacted) blocks before they are being processed.
	// Malformed blocks older than the maximum of consistency-delay and 48h0m0s will be removed.
	ConsistencyDelay metav1.Duration `json:"consistencyDelay,omitempty"`
	// How long to retain raw samples in bucket. 0d - disables this retention.
	RetentionResolutionRaw metav1.Duration `json:"retentionResolutionRaw,omitempty"`
	// How long to retain samples of resolution 1 (5 minutes) in bucket. 0d - disables this retention.
	RetentionResolution5m metav1.Duration `json:"retentionResolution5m,omitempty"`
	// How long to retain samples of resolution 2 (1 hour) in bucket. 0d - disables this retention.
	RetentionResolution1h metav1.Duration `json:"retentionResolution1h,omitempty"`
	// Do not exit after all compactions have been processed and wait for new work.
	Wait bool `json:"wait,omitempty"`
	// Disables downsampling. This is not recommended as querying long time ranges without non-downsampleddata
	// is not efficient and useful e.g it is not possible to render all samples for a human eye anyway.
	DownsamplingDisable bool `json:"downsamplingDisable,omitempty"`
	// Number of goroutines to use when syncing block metadata from object storage.
	// +kubebuilder:validation:Minimum=1
	BlockSyncConcurrency int `json:"blockSyncConcurrency,omitempty"`
	// Number of goroutines to use when compacting groups.
	// +kubebuilder:validation:Minimum=1
	CompactConcurrency int `json:"compactConcurrency,omitempty"`
}

var DefaultBucketWeb = &BucketWeb{
	BaseObject: BaseObject{
		Image: ImageSpec{
			Repository: thanosImageRepository,
			Tag:        thanosImageTag,
			PullPolicy: defaultPullPolicy,
		},
	},
	Metrics: &Metrics{
		Interval:       "15s",
		Timeout:        "5s",
		Path:           "/metrics",
		ServiceMonitor: false,
	},
	HTTPAddress:     "0.0.0.0:10902",
	HTTPGracePeriod: metav1.Duration{Duration: 2 * time.Minute},
	Refresh:         metav1.Duration{Duration: 30 * time.Minute},
	Timeout:         metav1.Duration{Duration: 5 * time.Minute},
}

type BucketWeb struct {
	BaseObject `json:",inline"`
	Metrics    *Metrics `json:"metrics,omitempty"`
	// Listen host:port for HTTP endpoints.
	HTTPAddress string `json:"httpAddress,omitempty"`
	// Time to wait after an interrupt received for HTTP Server.
	HTTPGracePeriod metav1.Duration `json:"httpGracePeriod,omitempty"`
	// Static prefix for all HTML links and redirect URLs in the bucket web UI interface. Actual endpoints are still served on / or the web.route-prefix. This allows thanos bucket web UI to be served behind a reverse proxy that strips a URL sub-path.
	WebExternalPrefix string `json:"web_external_prefix,omitempty"`
	// Name of HTTP request header used for dynamic prefixing of UI links and redirects. This option is ignored if web.external-prefix argument is set. Security risk: enable this option only if a reverse proxy in front of thanos is resetting the header. The --web.prefix-header=X-Forwarded-Prefix option can be useful, for example, if Thanos UI is served via Traefik reverse proxy with PathPrefixStrip option enabled, which sends the stripped prefix value in X-Forwarded-Prefix header. This allows thanos UI to be served on a sub-path.
	WebPrefixHeader string `json:"web_prefix_header,omitempty"`
	// Refresh interval to download metadata from remote storage.
	Refresh metav1.Duration `json:"refresh,omitempty"`
	// Timeout to download metadata from remote.
	Timeout metav1.Duration `json:"timeout,omitempty"`
	// Prometheus label to use as timeline title.
	Label string `json:"label,omitempty"`
}

// ObjectStoreStatus defines the observed state of ObjectStore
type ObjectStoreStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true

// ObjectStore is the Schema for the objectstores API
type ObjectStore struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ObjectStoreSpec   `json:"spec,omitempty"`
	Status ObjectStoreStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ObjectStoreList contains a list of ObjectStore
type ObjectStoreList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ObjectStore `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ObjectStore{}, &ObjectStoreList{})
}
