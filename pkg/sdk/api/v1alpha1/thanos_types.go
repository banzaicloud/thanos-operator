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
	"github.com/banzaicloud/operator-tools/pkg/typeoverride"
	"github.com/banzaicloud/operator-tools/pkg/volume"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	ThanosImageRepository = "quay.io/thanos/thanos"
	ThanosImageTag        = "v0.22.0"
)

var DefaultQueryFrontend = QueryFrontend{
	Metrics: &Metrics{
		Interval:       "15s",
		Timeout:        "5s",
		Path:           "/metrics",
		ServiceMonitor: false,
	},
	LogLevel:    "info",
	HttpAddress: "0.0.0.0:9090",
}

var DefaultQuery = Query{
	Metrics: &Metrics{
		Interval:       "15s",
		Timeout:        "5s",
		Path:           "/metrics",
		ServiceMonitor: false,
	},
	LogLevel:    "info",
	HttpAddress: "0.0.0.0:10902",
	GRPCAddress: "0.0.0.0:10901",
}

var DefaultStoreGateway = StoreGateway{
	Metrics: &Metrics{
		Interval:       "15s",
		Timeout:        "5s",
		Path:           "/metrics",
		ServiceMonitor: false,
	},
	LogLevel:    "info",
	HttpAddress: "0.0.0.0:10902",
	GRPCAddress: "0.0.0.0:10901",
}

var DefaultRule = Rule{
	DataDir: "/data",
	Metrics: &Metrics{
		Interval:       "15s",
		Timeout:        "5s",
		Path:           "/metrics",
		ServiceMonitor: false,
	},
	LogLevel:    "info",
	HttpAddress: "0.0.0.0:10902",
	GRPCAddress: "0.0.0.0:10901",
}

// ThanosSpec defines the desired state of Thanos
type ThanosSpec struct {
	QueryDiscovery                               bool           `json:"queryDiscovery,omitempty"`
	StoreGateway                                 *StoreGateway  `json:"storeGateway,omitempty"`
	Rule                                         *Rule          `json:"rule,omitempty"`
	Query                                        *Query         `json:"query,omitempty"`
	QueryFrontend                                *QueryFrontend `json:"queryFrontend,omitempty"`
	ClusterDomain                                string         `json:"clusterDomain,omitempty"`
	EnableRecreateWorkloadOnImmutableFieldChange bool           `json:"enableRecreateWorkloadOnImmutableFieldChange,omitempty"`
}

// Metrics defines the service monitor endpoints
type Metrics struct {
	Interval              string `json:"interval,omitempty"`
	Timeout               string `json:"timeout,omitempty"`
	Port                  int32  `json:"port,omitempty"`
	Path                  string `json:"path,omitempty"`
	ServiceMonitor        bool   `json:"serviceMonitor,omitempty"`
	PrometheusAnnotations bool   `json:"prometheusAnnotations,omitempty"`
}

type Ingress struct {
	// See [Ingress override](../overrides/override/#ingressnetworkingv1beta1)
	IngressOverrides *typeoverride.IngressNetworkingV1beta1 `json:"ingressOverrides,omitempty"`
	// Certificate in the ingress namespace
	Certificate string `json:"certificate,omitempty"`
	Host        string `json:"host,omitempty"`
	Path        string `json:"path,omitempty"`
}

type QueryFrontend struct {
	// See [ObjectMeta override](../overrides/override/#objectmeta)
	MetaOverrides *typeoverride.ObjectMeta `json:"metaOverrides,omitempty"`
	// See [Deployment override](../overrides/override/#deployment)
	DeploymentOverrides *typeoverride.Deployment `json:"deploymentOverrides,omitempty"`
	// See [Service override](../overrides/override/#service)
	ServiceOverrides *typeoverride.Service `json:"serviceOverrides,omitempty"`
	Metrics          *Metrics              `json:"metrics,omitempty"`
	HTTPIngress      *Ingress              `json:"HTTPIngress,omitempty"`
	LogLevel         string                `json:"logLevel,omitempty" thanos:"--log.level=%s"`
	LogFormat        string                `json:"logFormat,omitempty" thanos:"--log.format=%s"`
	// Split queries by an interval and execute in parallel, 0 disables it.
	QueryRangeSplit string `json:"queryRangeSplit,omitempty" thanos:"--query-range.split-interval=%s"`
	// Maximum number of retries for a single request; beyond this, the downstream error is returned.
	QueryRangeMaxRetriesPerRequest int `json:"queryRangeMaxRetriesPerRequest,omitempty" thanos:"--query-range.max-retries-per-request=%d"`
	// Limit the query time range (end - start time) in the query-frontend, 0 disables it.
	QueryRangeMaxQueryLength int `json:"queryRangeMaxQueryLength,omitempty" thanos:"--query-range.max-query-length=%d"`
	// Maximum number of queries will be scheduled in parallel by the frontend.
	QueryRangeMaxQueryParallelism int `json:"queryRangeMaxQueryParallelism,omitempty" thanos:"--query-range.max-query-parallelism=%d"`
	// Most recent allowed cacheable result, to prevent	caching very recent results that might still be in flux.
	QueryRangeResponseCacheMaxFreshness string `json:"queryRangeResponseCacheMaxFreshness,omitempty" thanos:"--query-range.response-cache-max-freshness=%s"`
	// Enable partial response for queries if no partial_response param is specified.
	QueryRangePartialResponse *bool `json:"queryRangePartialResponse,omitempty" thanos:"--query-range.partial-response"`
	// Path to YAML file that contains response cache configuration.
	QueryRangeResponseCacheConfigFile string `json:"queryRangeResponseCacheConfigFile,omitempty" thanos:"--query-range.response-cache-config-file=%s"`
	// Alternative to 'query-range.response-cache-config-file' flag (lower priority). Content of YAML file that contains response cache configuration.
	QueryRangeResponseCache string `json:"queryRangeResponseCache,omitempty" thanos:"--query-range.response-cache-config=%s"`
	// Listen host:port for HTTP endpoints.
	HttpAddress string `json:"httpAddress,omitempty" thanos:"--http-address=%s"`
	// Time to wait after an interrupt received for HTTP Server.
	HttpGracePeriod string `json:"http_grace_period,omitempty" thanos:"--http-grace-period=%s"`
	// URL of downstream Prometheus Query compatible API.
	QueryFrontendDownstreamURL string `json:"queryFrontendDownstreamURL,omitempty"`
	// Compress HTTP responses.
	QueryFrontendCompressResponses *bool `json:"queryFrontendCompressResponses,omitempty" thanos:"--query-frontend.compress-responses"`
	// 	Log queries that are slower than the specified duration. Set to 0 to disable. Set to < 0 to enable on all queries.
	QueryFrontendLogQueriesLongerThan int `json:"queryFrontendLogQueriesLongerThan,omitempty" thanos:"--query-frontend.log-queries-longer-than=%d"`
	// 	Request Logging for logging the start and end of requests. LogFinishCall is enabled by default.
	//	LogFinishCall : Logs the finish call of the requests.
	//	LogStartAndFinishCall : Logs the start and finish call of the requests.
	//	NoLogCall : Disable request logging.
	LogRequestDecision string `json:"logRequestDecision,omitempty" thanos:"--log.request.decision=%s"`
}

type Query struct {
	// See [ObjectMeta override](../overrides/override/#objectmeta)
	MetaOverrides typeoverride.ObjectMeta `json:"metaOverrides,omitempty"`
	// See [Deployment override](../overrides/override/#deployment)
	DeploymentOverrides *typeoverride.Deployment `json:"deploymentOverrides,omitempty"`
	// See [Service override](../overrides/override/#service)
	ServiceOverrides *typeoverride.Service `json:"serviceOverrides,omitempty"`
	Metrics          *Metrics              `json:"metrics,omitempty"`
	HTTPIngress      *Ingress              `json:"HTTPIngress,omitempty"`
	GRPCIngress      *Ingress              `json:"GRPCIngress,omitempty"`
	// Cert and key expected under tls.crt, tls.key
	GRPCClientCertificate string `json:"GRPCClientCertificate,omitempty"`
	// CA bundle to verify servers against, expected under ca.crt
	GRPCClientCA string `json:"GRPCClientCA,omitempty"`
	// Server name to verify server certificate against
	GRPCClientServerName string `json:"GRPCClientServerName,omitempty"`
	// Cert and key expected under tls.crt, tls.key
	GRPCServerCertificate string `json:"GRPCServerCertificate,omitempty"`
	// CA bundle to verify clients against, expected under ca.crt
	GRPCServerCA string `json:"GRPCServerCA,omitempty"`
	LogLevel     string `json:"logLevel,omitempty" thanos:"--log.level=%s"`
	LogFormat    string `json:"logFormat,omitempty" thanos:"--log.format=%s"`
	// Listen host:port for HTTP endpoints.
	HttpAddress string `json:"httpAddress,omitempty" thanos:"--http-address=%s"`
	// Time to wait after an interrupt received for HTTP Server.
	HttpGracePeriod string `json:"http_grace_period,omitempty" thanos:"--http-grace-period=%s"`
	// Listen ip:port address for gRPC endpoints
	GRPCAddress string `json:"grpcAddress,omitempty" thanos:"--grpc-address=%s"`
	// Time to wait after an interrupt received for GRPC Server.
	GRPCGracePeriod string `json:"grpcGracePeriod,omitempty" thanos:"--grpc-grace-period=%s"`
	// Prefix for API and UI endpoints. This allows thanos UI to be served on a sub-path. This
	// option is analogous to --web.route-prefix of Promethus.
	WebRoutePrefix string `json:"webRoutePrefix,omitempty" thanos:"--web.route-prefix=%s"`
	// Static prefix for all HTML links and redirect URLs in the UI query web interface. Actual
	// endpoints are still served on / or the web.route-prefix. This allows thanos UI to be
	// served behind a reverse proxy that strips a URL sub-path.
	WebExternalPrefix string `json:"webExternalPrefix,omitempty" thanos:"--web.external-prefix=%s"`
	// Name of HTTP request header used for dynamic prefixing of UI links and redirects. This
	// option is ignored if web.external-prefix argument is set. Security risk: enable this
	// option only if a reverse proxy in front of thanos is resetting the header. The
	// --web.prefix-header=X-Forwarded-Prefix option can be useful, for example, if Thanos UI is
	// served via Traefik reverse proxy with PathPrefixStrip option enabled, which sends the
	// stripped prefix value in X-Forwarded-Prefix header. This allows thanos UI to be served on a
	// sub-path.
	WebPrefixHeader string `json:"webPrefixHeader,omitempty" thanos:"--web.prefix-header=%s"`
	// Maximum time to process query by query node.
	QueryTimeout metav1.Duration `json:"queryTimeout,omitempty" thanos:"--query.timeout=%s"`
	// Maximum number of queries processed concurrently by query node.
	QueryMaxConcurrent int `json:"queryMaxConcurrent,omitempty" thanos:"--query.max-concurrent=%d"`
	// Labels to treat as a replica indicator along which data is deduplicated. Still you will be
	// able to query without deduplication using 'dedup=false' parameter.
	QueryReplicaLabels []string `json:"queryReplicaLabel,omitempty"`
	// Query selector labels that will be exposed in info endpoint (repeated).
	SelectorLabels map[string]string `json:"selectorLabels,omitempty"`
	// Addresses of statically configured store API servers (repeatable). The scheme may be
	// prefixed with 'dns+' or 'dnssrv+' to detect store API servers through respective DNS lookups.
	Stores []string `json:"stores,omitempty"`
	//	Interval between DNS resolutions.
	StoreSDDNSInterval metav1.Duration `json:"storeSDDNSInterval,omitempty" thanos:"--store.sd-dns-interval=%s"`
	//	Timeout before an unhealthy store is cleaned from the store UI page.
	StoreUnhealthyTimeout metav1.Duration `json:"storeUnhealthyTimeout,omitempty" thanos:"--store.unhealthy-timeout=%s"`
	// Enable automatic adjustment (step / 5) to what source of data should be used in store gateways
	// if no max_source_resolution param is specified.
	QueryAutoDownsampling bool `json:"queryAutoDownsampling,omitempty" thanos:"--query.auto-downsampling"`
	// Enable partial response for queries if no partial_response param is specified.
	QueryPartialResponse bool `json:"queryPartialResponse,omitempty" thanos:"--query.partial-response"`
	//	Set default evaluation interval for sub queries.
	QueryDefaultEvaluationInterval metav1.Duration `json:"queryDefaultEvaluationInterval,omitempty" thanos:" --query.default-evaluation-interval=%s"`
	//	If a Store doesn't send any data in this specified duration then a Store will be ignored
	//	and partial data will be returned if it's enabled. 0 disables timeout.
	StoreResponseTimeout metav1.Duration `json:"storeResponseTimeout,omitempty" thanos:"--store.response-timeout=%s"`
	// create Grafana data source
	GrafanaDatasource bool `json:"grafanaDatasource,omitempty"`
}

type ThanosDiscovery struct {
	metav1.LabelSelector `json:",omitempty,inline"`
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
	// See [ObjectMeta override](../overrides/override/#objectmeta)
	MetaOverrides *typeoverride.ObjectMeta `json:"metaOverrides,omitempty"`
	// See [Deployment override](../overrides/override/#deployment)
	DeploymentOverrides *typeoverride.Deployment `json:"deploymentOverrides,omitempty"`
	// See [Service override](../overrides/override/#service)
	ServiceOverrides      *typeoverride.Service `json:"serviceOverride,omitempty"`
	Metrics               *Metrics              `json:"metrics,omitempty"`
	GRPCServerCertificate string                `json:"GRPCServerCertificate,omitempty"`
	LogLevel              string                `json:"logLevel,omitempty" thanos:"--log.level=%s"`
	LogFormat             string                `json:"logFormat,omitempty" thanos:"--log.format=%s"`
	// Listen host:port for HTTP endpoints.
	HttpAddress string `json:"httpAddress,omitempty" thanos:"--http-address=%s"`
	// Time to wait after an interrupt received for HTTP Server.
	HttpGracePeriod string `json:"http_grace_period,omitempty" thanos:"--http-grace-period=%s"`
	// Listen ip:port address for gRPC endpoints
	GRPCAddress string `json:"grpcAddress,omitempty" thanos:"--grpc-address=%s"`
	// Time to wait after an interrupt received for GRPC Server.
	GRPCGracePeriod string `json:"grpcGracePeriod,omitempty" thanos:"--grpc-grace-period=%s"`
	// Maximum size of items held in the in-memory index cache.
	IndexCacheSize string `json:"indexCacheSize,omitempty" thanos:"--index-cache-size=%s"`
	// Path to YAML file that contains index cache configuration. See format details:
	// https://thanos.io/tip/components/store.md/#index-cache
	IndexCacheConfigFile string `json:"indexCacheConfigFile,omitempty" thanos:"index-cache.config-file=%s"`
	// Alternative to 'index-cache.config-file' flag (lower priority). Content of YAML file that contains index cache configuration. See format details:
	// https://thanos.io/tip/components/store.md/#index-cache
	IndexCacheConfig string `json:"indexCacheConfig,omitempty" thanos:"--index-cache.config=%s"`
	// Maximum size of concurrently allocatable bytes for chunks.
	ChunkPoolSize string `json:"chunkPoolSize,omitempty" thanos:"--chunk-pool-size=%s"`
	// Maximum amount of samples returned via a single Series call. 0 means no limit. NOTE: For
	// efficiency we take 120 as the number of samples in chunk (it cannot be bigger than that), so
	// the actual number of samples might be lower, even though the maximum could be hit.
	StoreGRPCSeriesSampleLimit string `json:"storeGRPCSeriesSampleLimit,omitempty" thanos:"--store.grpc.series-sample-limit=%s"`
	// Maximum amount of touched series returned via a single Series call. The Series call fails if this limit is exceeded. 0 means no limit.
	StoreGRPCTouchedSeriesSampleLimit int `json:"storeGRPCTouchedSeriesSampleLimit,omitempty" thanos:"--store.grpc.touched-series-limit=%d"`
	// Maximum number of concurrent Series calls.
	StoreGRPCSeriesMaxConcurrency int `json:"storeGRPCSeriesMaxConcurrency,omitempty" thanos:"--store.grpc.series-max-concurrency=%d"`
	// Repeat interval for syncing the blocks between local and remote view.
	SyncBlockDuration string `json:"syncBlockDuration,omitempty" thanos:"--sync-block-duration=%s"`
	// Number of goroutines to use when constructing index-cache.json blocks from object storage.
	BlockSyncConcurrency int `json:"blockSyncConcurrency,omitempty" thanos:"--block-sync-concurrency=%d"`
	// Number of goroutines to use when fetching block metadata from object storage.
	BlockMetaFetchConcurrency int `json:"blockMetaFetchConcurrency,omitempty" thanos:"--block-meta-fetch-concurrency=%d"`
	// Path to YAML file that contains relabeling configuration that allows selecting blocks. It
	// follows native Prometheus relabel-config syntax. See format details:
	// https://prometheus.io/docs/prometheus/latest/configuration/configuration/#relabel_config
	SelectorRelabelConfigFile string `json:"selectorRelabelConfigFile,omitempty" thanos:"--selector.relabel-config-file=%s"`
	// Alternative to 'selector.relabel-config-file' flag (lower priority). Content of YAML file
	// that contains relabeling configuration that allows selecting blocks. It follows native
	// Prometheus relabel-config syntax. See format details:
	// https://prometheus.io/docs/prometheus/latest/configuration/configuration/#relabel_config
	SelectorRelabelConfig string `json:"selectorRelabelConfig,omitempty" thanos:"--selector.relabel-config=%s"`
	// Minimum age of all blocks before they are being read. Set it to safe value (e.g 30m) if your
	// object storage is eventually consistent. GCS and S3 are (roughly) strongly consistent.
	ConsistencyDelay string `json:"consistencyDelay,omitempty" thanos:"--consistency-delay=%s"`
	// Duration after which the blocks marked for deletion will be filtered out while fetching blocks. The idea of ignore-deletion-marks-delay
	// is to ignore blocks that are marked for deletion with some delay. This ensures store can still serve blocks that are meant to be
	// deleted but do not have a replacement yet. If delete-delay duration is provided to compactor or bucket verify component, it will upload
	// deletion-mark.json file to mark after what duration the block should be deleted rather than deleting the block straight away. If
	// delete-delay is non-zero for compactor or bucket verify component, ignore-deletion-marks-delay should be set to
	// (delete-delay)/2 so that blocks marked for deletion are filtered out while fetching blocks
	// before being deleted from bucket. Default is 24h, half of the default value for --delete-delay on compactor.
	IgnoreDeletionMarksDelay string `json:"ignoreDeletionMarksDelay,omitempty" thanos:"--ignore-deletion-marks-delay=%s"`
	// 	If true, Store Gateway will lazy memory map index-header only once the block is required by a query.
	StoreEnableIndexHeaderLazyReader *bool `json:"storeEnableIndexHeaderLazyReader,omitempty" thanos:"--store.enable-index-header-lazy-reader"`
	// Static prefix for all HTML links and redirect URLs in the bucket web UI interface. Actual endpoints are still served on / or the
	// web.route-prefix. This allows thanos bucket web UI to be served behind a reverse proxy that
	// strips a URL sub-path.
	WebExternalPrefix string `json:"webExternalPrefix,omitempty" thanos:"--web.external-prefix=%s"`
	// Name of HTTP request header used for dynamic prefixing of UI links and redirects. This
	// option is ignored if web.external-prefix argument is set. Security risk: enable this
	// option only if a reverse proxy in front of thanos is resetting the header. The
	// --web.prefix-header=X-Forwarded-Prefix option can be useful, for example, if Thanos UI is
	// served via Traefik reverse proxy with PathPrefixStrip option enabled, which sends the
	// stripped prefix value in X-Forwarded-Prefix header. This allows thanos UI to be served on a sub-path.
	WebPrefixHeader string `json:"webPrefixHeader,omitempty" thanos:"--web.prefix-header=%s"`
	// TimeRanges is a list of TimeRange to partition Store Gateway
	TimeRanges []TimeRange `json:"timeRanges,omitempty"`
}

type Rule struct {
	// See [ObjectMeta override](../overrides/override/#objectmeta)
	MetaOverrides *typeoverride.ObjectMeta `json:"metaOverrides,omitempty"`
	// See [StatefulSet override](../overrides/override/#statefulset)
	StatefulsetOverrides *typeoverride.StatefulSet `json:"statefulsetOverrides,omitempty"`
	// See [Service override](../overrides/override/#service)
	ServiceOverrides *typeoverride.Service `json:"serviceOverrides,omitempty"`
	Metrics          *Metrics              `json:"metrics,omitempty"`
	HTTPIngress      *Ingress              `json:"HTTPIngress,omitempty"`
	GRPCIngress      *Ingress              `json:"GRPCIngress,omitempty"`
	LogLevel         string                `json:"logLevel,omitempty" thanos:"--log.level=%s"`
	LogFormat        string                `json:"logFormat,omitempty" thanos:"--log.format=%s"`
	// Listen host:port for HTTP endpoints.
	HttpAddress string `json:"httpAddress,omitempty" thanos:"--http-address=%s"`
	// Time to wait after an interrupt received for HTTP Server.
	HttpGracePeriod string `json:"http_grace_period,omitempty" thanos:"--http-grace-period=%s"`
	// Data directory.
	DataDir string `json:"dataDir,omitempty"`
	// Kubernetes volume abstraction refers to different types of volumes to be mounted to pods: emptyDir, hostPath, pvc.
	DataVolume *volume.KubernetesVolume `json:"dataVolume,omitempty"`
	// Listen ip:port address for gRPC endpoints
	GRPCAddress string `json:"grpcAddress,omitempty" thanos:"--grpc-address=%s"`
	// Time to wait after an interrupt received for GRPC Server.
	GRPCGracePeriod string `json:"grpcGracePeriod,omitempty" thanos:"--grpc-grace-period=%s"`
	// 	Labels to be applied to all generated metrics
	//(repeated). Similar to external labels for
	//	Prometheus, used to identify ruler and its
	//	blocks as unique source.
	Labels map[string]string `json:"labels,omitempty"`
	// Rules
	Rules string `json:"rules,omitempty"`
	// Minimum amount of time to wait before resending an alert to Alertmanager.
	ResendDelay string `json:"resendDelay,omitempty" thanos:"--resend-delay=%s"`
	// The default evaluation interval to use.
	EvalInterval string `json:"evalInterval,omitempty" thanos:"--eval-interval=%s"`
	// Block duration for TSDB block.
	TSDBBlockDuration string `json:"tsdbBlockDuration,omitempty" thanos:"--tsdb.block-duration=%s"`
	// Block retention time on local disk.
	TSDBRetention string `json:"tsdbRetention,omitempty" thanos:"--tsdb.retention=%s"`
	// Alertmanager replica URLs to push firing alerts. Ruler claims success if push to at
	// least one alertmanager from discovered succeeds. The scheme should not be empty e.g
	// `http` might be used. The scheme may be prefixed with 'dns+' or 'dnssrv+' to detect
	// Alertmanager IPs through respective DNS lookups. The port defaults to 9093 or the SRV
	// record's value. The URL path is used as a prefix for the regular Alertmanager API path.
	AlertmanagersURLs []string `json:"alertmanagersURLs,omitempty"`
	// Timeout for sending alerts to Alertmanager
	AlertmanagersSendTimeout string `json:"alertmanagersSendTimeout,omitempty" thanos:"--alertmanagers.send-timeout=%s"`
	// Interval between DNS resolutions of Alertmanager hosts.
	AlertmanagersSDDNSInterval string `json:"alertmanagersSDDNSInterval,omitempty" thanos:"--alertmanagers.sd-dns-interval=%s"`
	// The external Thanos Query URL that would be set in all alerts 'Source' field
	AlertQueryURL string `json:"alertQueryUrl,omitempty" thanos:"--alert.query-url=%s"`
	// Labels by name to drop before sending to alertmanager. This allows alert to be
	// deduplicated on replica label (repeated). Similar Prometheus alert relabelling
	AlertLabelDrop map[string]string `json:"alertLabelDrop,omitempty"`
	// Prefix for API and UI endpoints. This allows thanos UI to be served on a sub-path. This
	// option is analogous to --web.route-prefix of Promethus.
	WebRoutePrefix string `json:"webRoutePrefix,omitempty" thanos:"--web.route-prefix=%s"`
	// Static prefix for all HTML links and redirect URLs in the UI query web interface. Actual
	// endpoints are still served on / or the web.route-prefix. This allows thanos UI to be
	// served behind a reverse proxy that strips a URL sub-path.
	WebExternalPrefix string `json:"webExternalPrefix,omitempty" thanos:"--web.external-prefix=%s"`
	// Name of HTTP request header used for dynamic prefixing of UI links and redirects. This
	// option is ignored if web.external-prefix argument is set. Security risk: enable this
	// option only if a reverse proxy in front of thanos is resetting the header. The
	// --web.prefix-header=X-Forwarded-Prefix option can be useful, for example, if Thanos UI is
	// served via Traefik reverse proxy with PathPrefixStrip option enabled, which sends the
	// stripped prefix value in X-Forwarded-Prefix header. This allows thanos UI to be served on a
	// sub-path.
	WebPrefixHeader string `json:"webPrefixHeader,omitempty" thanos:"--web.prefix-header=%s"`
	// Addresses of statically configured query API servers (repeatable). The scheme may be
	// prefixed with 'dns+' or 'dnssrv+' to detect query API servers through respective DNS
	// lookups.
	Queries []string `json:"queries,omitempty"`
	// Interval between DNS resolutions.
	QuerySDDNSInterval string `json:"querySddnsInterval,omitempty" thanos:"--query.sd-dns-interval=%s"`
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

func (t *Thanos) GetClusterDomain() string {
	if t.Spec.ClusterDomain != "" {
		return t.Spec.ClusterDomain
	}
	return "cluster.local"
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
