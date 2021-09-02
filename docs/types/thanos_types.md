## ThanosSpec

ThanosSpec defines the desired state of Thanos

### queryDiscovery (bool, optional) {#thanosspec-querydiscovery}

Default: -

### storeGateway (*StoreGateway, optional) {#thanosspec-storegateway}

Default: -

### rule (*Rule, optional) {#thanosspec-rule}

Default: -

### query (*Query, optional) {#thanosspec-query}

Default: -

### queryFrontend (*QueryFrontend, optional) {#thanosspec-queryfrontend}

Default: -

### clusterDomain (string, optional) {#thanosspec-clusterdomain}

Default: -

### enableRecreateWorkloadOnImmutableFieldChange (bool, optional) {#thanosspec-enablerecreateworkloadonimmutablefieldchange}

Default: -


## Metrics

Metrics defines the service monitor endpoints

### interval (string, optional) {#metrics-interval}

Default: -

### timeout (string, optional) {#metrics-timeout}

Default: -

### port (int32, optional) {#metrics-port}

Default: -

### path (string, optional) {#metrics-path}

Default: -

### serviceMonitor (bool, optional) {#metrics-servicemonitor}

Default: -

### prometheusAnnotations (bool, optional) {#metrics-prometheusannotations}

Default: -


## Ingress

### ingressOverrides (*typeoverride.IngressNetworkingV1beta1, optional) {#ingress-ingressoverrides}

See [Ingress override](../overrides/override/#ingressnetworkingv1beta1) 

Default: -

### certificate (string, optional) {#ingress-certificate}

Certificate in the ingress namespace 

Default: -

### host (string, optional) {#ingress-host}

Default: -

### path (string, optional) {#ingress-path}

Default: -


## QueryFrontend

### metaOverrides (*typeoverride.ObjectMeta, optional) {#queryfrontend-metaoverrides}

See [ObjectMeta override](../overrides/override/#objectmeta) 

Default: -

### deploymentOverrides (*typeoverride.Deployment, optional) {#queryfrontend-deploymentoverrides}

See [Deployment override](../overrides/override/#deployment) 

Default: -

### serviceOverrides (*typeoverride.Service, optional) {#queryfrontend-serviceoverrides}

See [Service override](../overrides/override/#service) 

Default: -

### metrics (*Metrics, optional) {#queryfrontend-metrics}

Default: -

### HTTPIngress (*Ingress, optional) {#queryfrontend-httpingress}

Default: -

### logLevel (string, optional) {#queryfrontend-loglevel}

Default: -

### logFormat (string, optional) {#queryfrontend-logformat}

Default: -

### queryRangeSplit (string, optional) {#queryfrontend-queryrangesplit}

Split queries by an interval and execute in parallel, 0 disables it. 

Default: -

### queryRangeMaxRetriesPerRequest (int, optional) {#queryfrontend-queryrangemaxretriesperrequest}

Maximum number of retries for a single request; beyond this, the downstream error is returned. 

Default: -

### queryRangeMaxQueryLength (int, optional) {#queryfrontend-queryrangemaxquerylength}

Limit the query time range (end - start time) in the query-frontend, 0 disables it. 

Default: -

### queryRangeMaxQueryParallelism (int, optional) {#queryfrontend-queryrangemaxqueryparallelism}

Maximum number of queries will be scheduled in parallel by the frontend. 

Default: -

### queryRangeResponseCacheMaxFreshness (string, optional) {#queryfrontend-queryrangeresponsecachemaxfreshness}

Most recent allowed cacheable result, to prevent	caching very recent results that might still be in flux. 

Default: -

### queryRangePartialResponse (*bool, optional) {#queryfrontend-queryrangepartialresponse}

Enable partial response for queries if no partial_response param is specified. 

Default: -

### queryRangeResponseCacheConfigFile (string, optional) {#queryfrontend-queryrangeresponsecacheconfigfile}

Path to YAML file that contains response cache configuration. 

Default: -

### queryRangeResponseCache (string, optional) {#queryfrontend-queryrangeresponsecache}

Alternative to 'query-range.response-cache-config-file' flag (lower priority). Content of YAML file that contains response cache configuration. 

Default: -

### httpAddress (string, optional) {#queryfrontend-httpaddress}

Listen host:port for HTTP endpoints. 

Default: -

### http_grace_period (string, optional) {#queryfrontend-http_grace_period}

Time to wait after an interrupt received for HTTP Server. 

Default: -

### queryFrontendDownstreamURL (string, optional) {#queryfrontend-queryfrontenddownstreamurl}

URL of downstream Prometheus Query compatible API. 

Default: -

### queryFrontendCompressResponses (*bool, optional) {#queryfrontend-queryfrontendcompressresponses}

Compress HTTP responses. 

Default: -

### queryFrontendLogQueriesLongerThan (int, optional) {#queryfrontend-queryfrontendlogquerieslongerthan}

Log queries that are slower than the specified duration. Set to 0 to disable. Set to < 0 to enable on all queries. 

Default: -

### logRequestDecision (string, optional) {#queryfrontend-logrequestdecision}

Request Logging for logging the start and end of requests. LogFinishCall is enabled by default. LogFinishCall : Logs the finish call of the requests. LogStartAndFinishCall : Logs the start and finish call of the requests. NoLogCall : Disable request logging. 

Default: -


## Query

### metaOverrides (typeoverride.ObjectMeta, optional) {#query-metaoverrides}

See [ObjectMeta override](../overrides/override/#objectmeta) 

Default: -

### deploymentOverrides (*typeoverride.Deployment, optional) {#query-deploymentoverrides}

See [Deployment override](../overrides/override/#deployment) 

Default: -

### serviceOverrides (*typeoverride.Service, optional) {#query-serviceoverrides}

See [Service override](../overrides/override/#service) 

Default: -

### metrics (*Metrics, optional) {#query-metrics}

Default: -

### HTTPIngress (*Ingress, optional) {#query-httpingress}

Default: -

### GRPCIngress (*Ingress, optional) {#query-grpcingress}

Default: -

### GRPCClientCertificate (string, optional) {#query-grpcclientcertificate}

Cert and key expected under tls.crt, tls.key 

Default: -

### GRPCClientCA (string, optional) {#query-grpcclientca}

CA bundle to verify servers against, expected under ca.crt 

Default: -

### GRPCClientServerName (string, optional) {#query-grpcclientservername}

Server name to verify server certificate against 

Default: -

### GRPCServerCertificate (string, optional) {#query-grpcservercertificate}

Cert and key expected under tls.crt, tls.key 

Default: -

### GRPCServerCA (string, optional) {#query-grpcserverca}

CA bundle to verify clients against, expected under ca.crt 

Default: -

### logLevel (string, optional) {#query-loglevel}

Default: -

### logFormat (string, optional) {#query-logformat}

Default: -

### httpAddress (string, optional) {#query-httpaddress}

Listen host:port for HTTP endpoints. 

Default: -

### http_grace_period (string, optional) {#query-http_grace_period}

Time to wait after an interrupt received for HTTP Server. 

Default: -

### grpcAddress (string, optional) {#query-grpcaddress}

Listen ip:port address for gRPC endpoints 

Default: -

### grpcGracePeriod (string, optional) {#query-grpcgraceperiod}

Time to wait after an interrupt received for GRPC Server. 

Default: -

### webRoutePrefix (string, optional) {#query-webrouteprefix}

Prefix for API and UI endpoints. This allows thanos UI to be served on a sub-path. This option is analogous to --web.route-prefix of Promethus. 

Default: -

### webExternalPrefix (string, optional) {#query-webexternalprefix}

Static prefix for all HTML links and redirect URLs in the UI query web interface. Actual endpoints are still served on / or the web.route-prefix. This allows thanos UI to be served behind a reverse proxy that strips a URL sub-path. 

Default: -

### webPrefixHeader (string, optional) {#query-webprefixheader}

Name of HTTP request header used for dynamic prefixing of UI links and redirects. This option is ignored if web.external-prefix argument is set. Security risk: enable this option only if a reverse proxy in front of thanos is resetting the header. The --web.prefix-header=X-Forwarded-Prefix option can be useful, for example, if Thanos UI is served via Traefik reverse proxy with PathPrefixStrip option enabled, which sends the stripped prefix value in X-Forwarded-Prefix header. This allows thanos UI to be served on a sub-path. 

Default: -

### queryTimeout (metav1.Duration, optional) {#query-querytimeout}

Maximum time to process query by query node. 

Default: -

### queryMaxConcurrent (int, optional) {#query-querymaxconcurrent}

Maximum number of queries processed concurrently by query node. 

Default: -

### queryReplicaLabel ([]string, optional) {#query-queryreplicalabel}

Labels to treat as a replica indicator along which data is deduplicated. Still you will be able to query without deduplication using 'dedup=false' parameter. 

Default: -

### selectorLabels (map[string]string, optional) {#query-selectorlabels}

Query selector labels that will be exposed in info endpoint (repeated). 

Default: -

### stores ([]string, optional) {#query-stores}

Addresses of statically configured store API servers (repeatable). The scheme may be prefixed with 'dns+' or 'dnssrv+' to detect store API servers through respective DNS lookups. 

Default: -

### storeSDDNSInterval (metav1.Duration, optional) {#query-storesddnsinterval}

Interval between DNS resolutions. 

Default: -

### storeUnhealthyTimeout (metav1.Duration, optional) {#query-storeunhealthytimeout}

Timeout before an unhealthy store is cleaned from the store UI page. 

Default: -

### queryAutoDownsampling (bool, optional) {#query-queryautodownsampling}

Enable automatic adjustment (step / 5) to what source of data should be used in store gateways if no max_source_resolution param is specified. 

Default: -

### queryPartialResponse (bool, optional) {#query-querypartialresponse}

Enable partial response for queries if no partial_response param is specified. 

Default: -

### queryDefaultEvaluationInterval (metav1.Duration, optional) {#query-querydefaultevaluationinterval}

Set default evaluation interval for sub queries. 

Default: -

### storeResponseTimeout (metav1.Duration, optional) {#query-storeresponsetimeout}

If a Store doesn't send any data in this specified duration then a Store will be ignored and partial data will be returned if it's enabled. 0 disables timeout. 

Default: -

### grafanaDatasource (bool, optional) {#query-grafanadatasource}

create Grafana data source 

Default: -


## ThanosDiscovery

###  (metav1.LabelSelector, optional) {#thanosdiscovery-}

Default: -


## TimeRange

### minTime (string, optional) {#timerange-mintime}

Start of time range limit to serve. Thanos Store will serve only metrics, which happened later than this value. Option can be a constant time in RFC3339 format or time duration relative to current time, such as -1d or 2h45m. Valid duration units are ms, s, m, h, d, w, y. 

Default: -

### maxTime (string, optional) {#timerange-maxtime}

End of time range limit to serve. Thanos Store will serve only blocks, which happened eariler than this value. Option can be a constant time in RFC3339 format or time duration relative to current time, such as -1d or 2h45m. Valid duration units are ms, s, m, h, d, w, y. 

Default: -


## StoreGateway

### metaOverrides (*typeoverride.ObjectMeta, optional) {#storegateway-metaoverrides}

See [ObjectMeta override](../overrides/override/#objectmeta) 

Default: -

### deploymentOverrides (*typeoverride.Deployment, optional) {#storegateway-deploymentoverrides}

See [Deployment override](../overrides/override/#deployment) 

Default: -

### serviceOverride (*typeoverride.Service, optional) {#storegateway-serviceoverride}

See [Service override](../overrides/override/#service) 

Default: -

### metrics (*Metrics, optional) {#storegateway-metrics}

Default: -

### GRPCServerCertificate (string, optional) {#storegateway-grpcservercertificate}

Default: -

### logLevel (string, optional) {#storegateway-loglevel}

Default: -

### logFormat (string, optional) {#storegateway-logformat}

Default: -

### httpAddress (string, optional) {#storegateway-httpaddress}

Listen host:port for HTTP endpoints. 

Default: -

### http_grace_period (string, optional) {#storegateway-http_grace_period}

Time to wait after an interrupt received for HTTP Server. 

Default: -

### grpcAddress (string, optional) {#storegateway-grpcaddress}

Listen ip:port address for gRPC endpoints 

Default: -

### grpcGracePeriod (string, optional) {#storegateway-grpcgraceperiod}

Time to wait after an interrupt received for GRPC Server. 

Default: -

### indexCacheSize (string, optional) {#storegateway-indexcachesize}

Maximum size of items held in the in-memory index cache. 

Default: -

### indexCacheConfigFile (string, optional) {#storegateway-indexcacheconfigfile}

Path to YAML file that contains index cache configuration. See format details: https://thanos.io/tip/components/store.md/#index-cache 

Default: -

### indexCacheConfig (string, optional) {#storegateway-indexcacheconfig}

Alternative to 'index-cache.config-file' flag (lower priority). Content of YAML file that contains index cache configuration. See format details: https://thanos.io/tip/components/store.md/#index-cache 

Default: -

### chunkPoolSize (string, optional) {#storegateway-chunkpoolsize}

Maximum size of concurrently allocatable bytes for chunks. 

Default: -

### storeGRPCSeriesSampleLimit (string, optional) {#storegateway-storegrpcseriessamplelimit}

Maximum amount of samples returned via a single Series call. 0 means no limit. NOTE: For efficiency we take 120 as the number of samples in chunk (it cannot be bigger than that), so the actual number of samples might be lower, even though the maximum could be hit. 

Default: -

### storeGRPCTouchedSeriesSampleLimit (int, optional) {#storegateway-storegrpctouchedseriessamplelimit}

Maximum amount of touched series returned via a single Series call. The Series call fails if this limit is exceeded. 0 means no limit. 

Default: -

### storeGRPCSeriesMaxConcurrency (int, optional) {#storegateway-storegrpcseriesmaxconcurrency}

Maximum number of concurrent Series calls. 

Default: -

### syncBlockDuration (string, optional) {#storegateway-syncblockduration}

Repeat interval for syncing the blocks between local and remote view. 

Default: -

### blockSyncConcurrency (int, optional) {#storegateway-blocksyncconcurrency}

Number of goroutines to use when constructing index-cache.json blocks from object storage. 

Default: -

### blockMetaFetchConcurrency (int, optional) {#storegateway-blockmetafetchconcurrency}

Number of goroutines to use when fetching block metadata from object storage. 

Default: -

### selectorRelabelConfigFile (string, optional) {#storegateway-selectorrelabelconfigfile}

Path to YAML file that contains relabeling configuration that allows selecting blocks. It follows native Prometheus relabel-config syntax. See format details: https://prometheus.io/docs/prometheus/latest/configuration/configuration/#relabel_config 

Default: -

### selectorRelabelConfig (string, optional) {#storegateway-selectorrelabelconfig}

Alternative to 'selector.relabel-config-file' flag (lower priority). Content of YAML file that contains relabeling configuration that allows selecting blocks. It follows native Prometheus relabel-config syntax. See format details: https://prometheus.io/docs/prometheus/latest/configuration/configuration/#relabel_config 

Default: -

### consistencyDelay (string, optional) {#storegateway-consistencydelay}

Minimum age of all blocks before they are being read. Set it to safe value (e.g 30m) if your object storage is eventually consistent. GCS and S3 are (roughly) strongly consistent. 

Default: -

### ignoreDeletionMarksDelay (string, optional) {#storegateway-ignoredeletionmarksdelay}

Duration after which the blocks marked for deletion will be filtered out while fetching blocks. The idea of ignore-deletion-marks-delay is to ignore blocks that are marked for deletion with some delay. This ensures store can still serve blocks that are meant to be deleted but do not have a replacement yet. If delete-delay duration is provided to compactor or bucket verify component, it will upload deletion-mark.json file to mark after what duration the block should be deleted rather than deleting the block straight away. If delete-delay is non-zero for compactor or bucket verify component, ignore-deletion-marks-delay should be set to (delete-delay)/2 so that blocks marked for deletion are filtered out while fetching blocks before being deleted from bucket. Default is 24h, half of the default value for --delete-delay on compactor. 

Default: -

### storeEnableIndexHeaderLazyReader (*bool, optional) {#storegateway-storeenableindexheaderlazyreader}

If true, Store Gateway will lazy memory map index-header only once the block is required by a query. 

Default: -

### webExternalPrefix (string, optional) {#storegateway-webexternalprefix}

Static prefix for all HTML links and redirect URLs in the bucket web UI interface. Actual endpoints are still served on / or the web.route-prefix. This allows thanos bucket web UI to be served behind a reverse proxy that strips a URL sub-path. 

Default: -

### webPrefixHeader (string, optional) {#storegateway-webprefixheader}

Name of HTTP request header used for dynamic prefixing of UI links and redirects. This option is ignored if web.external-prefix argument is set. Security risk: enable this option only if a reverse proxy in front of thanos is resetting the header. The --web.prefix-header=X-Forwarded-Prefix option can be useful, for example, if Thanos UI is served via Traefik reverse proxy with PathPrefixStrip option enabled, which sends the stripped prefix value in X-Forwarded-Prefix header. This allows thanos UI to be served on a sub-path. 

Default: -

### timeRanges ([]TimeRange, optional) {#storegateway-timeranges}

TimeRanges is a list of TimeRange to partition Store Gateway 

Default: -


## Rule

### metaOverrides (*typeoverride.ObjectMeta, optional) {#rule-metaoverrides}

See [ObjectMeta override](../overrides/override/#objectmeta) 

Default: -

### statefulsetOverrides (*typeoverride.StatefulSet, optional) {#rule-statefulsetoverrides}

See [StatefulSet override](../overrides/override/#statefulset) 

Default: -

### serviceOverrides (*typeoverride.Service, optional) {#rule-serviceoverrides}

See [Service override](../overrides/override/#service) 

Default: -

### metrics (*Metrics, optional) {#rule-metrics}

Default: -

### HTTPIngress (*Ingress, optional) {#rule-httpingress}

Default: -

### GRPCIngress (*Ingress, optional) {#rule-grpcingress}

Default: -

### logLevel (string, optional) {#rule-loglevel}

Default: -

### logFormat (string, optional) {#rule-logformat}

Default: -

### httpAddress (string, optional) {#rule-httpaddress}

Listen host:port for HTTP endpoints. 

Default: -

### http_grace_period (string, optional) {#rule-http_grace_period}

Time to wait after an interrupt received for HTTP Server. 

Default: -

### dataDir (string, optional) {#rule-datadir}

Data directory. 

Default: -

### dataVolume (*volume.KubernetesVolume, optional) {#rule-datavolume}

Kubernetes volume abstraction refers to different types of volumes to be mounted to pods: emptyDir, hostPath, pvc. 

Default: -

### grpcAddress (string, optional) {#rule-grpcaddress}

Listen ip:port address for gRPC endpoints 

Default: -

### grpcGracePeriod (string, optional) {#rule-grpcgraceperiod}

Time to wait after an interrupt received for GRPC Server. 

Default: -

### labels (map[string]string, optional) {#rule-labels}

Labels to be applied to all generated metrics (repeated). Similar to external labels for Prometheus, used to identify ruler and its blocks as unique source. 

Default: -

### rules (string, optional) {#rule-rules}

Rules 

Default: -

### resendDelay (string, optional) {#rule-resenddelay}

Minimum amount of time to wait before resending an alert to Alertmanager. 

Default: -

### evalInterval (string, optional) {#rule-evalinterval}

The default evaluation interval to use. 

Default: -

### tsdbBlockDuration (string, optional) {#rule-tsdbblockduration}

Block duration for TSDB block. 

Default: -

### tsdbRetention (string, optional) {#rule-tsdbretention}

Block retention time on local disk. 

Default: -

### alertmanagersURLs ([]string, optional) {#rule-alertmanagersurls}

Alertmanager replica URLs to push firing alerts. Ruler claims success if push to at least one alertmanager from discovered succeeds. The scheme should not be empty e.g `http` might be used. The scheme may be prefixed with 'dns+' or 'dnssrv+' to detect Alertmanager IPs through respective DNS lookups. The port defaults to 9093 or the SRV record's value. The URL path is used as a prefix for the regular Alertmanager API path. 

Default: -

### alertmanagersSendTimeout (string, optional) {#rule-alertmanagerssendtimeout}

Timeout for sending alerts to Alertmanager 

Default: -

### alertmanagersSDDNSInterval (string, optional) {#rule-alertmanagerssddnsinterval}

Interval between DNS resolutions of Alertmanager hosts. 

Default: -

### alertQueryUrl (string, optional) {#rule-alertqueryurl}

The external Thanos Query URL that would be set in all alerts 'Source' field 

Default: -

### alertLabelDrop (map[string]string, optional) {#rule-alertlabeldrop}

Labels by name to drop before sending to alertmanager. This allows alert to be deduplicated on replica label (repeated). Similar Prometheus alert relabelling 

Default: -

### webRoutePrefix (string, optional) {#rule-webrouteprefix}

Prefix for API and UI endpoints. This allows thanos UI to be served on a sub-path. This option is analogous to --web.route-prefix of Promethus. 

Default: -

### webExternalPrefix (string, optional) {#rule-webexternalprefix}

Static prefix for all HTML links and redirect URLs in the UI query web interface. Actual endpoints are still served on / or the web.route-prefix. This allows thanos UI to be served behind a reverse proxy that strips a URL sub-path. 

Default: -

### webPrefixHeader (string, optional) {#rule-webprefixheader}

Name of HTTP request header used for dynamic prefixing of UI links and redirects. This option is ignored if web.external-prefix argument is set. Security risk: enable this option only if a reverse proxy in front of thanos is resetting the header. The --web.prefix-header=X-Forwarded-Prefix option can be useful, for example, if Thanos UI is served via Traefik reverse proxy with PathPrefixStrip option enabled, which sends the stripped prefix value in X-Forwarded-Prefix header. This allows thanos UI to be served on a sub-path. 

Default: -

### queries ([]string, optional) {#rule-queries}

Addresses of statically configured query API servers (repeatable). The scheme may be prefixed with 'dns+' or 'dnssrv+' to detect query API servers through respective DNS lookups. 

Default: -

### querySddnsInterval (string, optional) {#rule-querysddnsinterval}

Interval between DNS resolutions. 

Default: -


## ThanosStatus

ThanosStatus defines the observed state of Thanos


## Thanos

Thanos is the Schema for the thanos API

###  (metav1.TypeMeta, required) {#thanos-}

Default: -

### metadata (metav1.ObjectMeta, optional) {#thanos-metadata}

Default: -

### spec (ThanosSpec, optional) {#thanos-spec}

Default: -

### status (ThanosStatus, optional) {#thanos-status}

Default: -


## ThanosList

ThanosList contains a list of Thanos

###  (metav1.TypeMeta, required) {#thanoslist-}

Default: -

### metadata (metav1.ListMeta, optional) {#thanoslist-metadata}

Default: -

### items ([]Thanos, required) {#thanoslist-items}

Default: -


