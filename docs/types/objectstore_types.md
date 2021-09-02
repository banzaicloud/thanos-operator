## ObjectStoreSpec

ObjectStoreSpec defines the desired state of ObjectStore

### config (secret.Secret, required) {#objectstorespec-config}

Config 

Default: -

### compactor (*Compactor, optional) {#objectstorespec-compactor}

See [Compactor](#compactor) 

Default: -

### bucketWeb (*BucketWeb, optional) {#objectstorespec-bucketweb}

See [BucketWeb](#bucketweb) 

Default: -


## Compactor

### metaOverrides (*typeoverride.ObjectMeta, optional) {#compactor-metaoverrides}

See [ObjectMeta override](../overrides/override/#objectmeta) 

Default: -

### deploymentOverrides (*typeoverride.Deployment, optional) {#compactor-deploymentoverrides}

See [Deployment override](../overrides/override/#deployment) 

Default: -

### serviceOverrides (*typeoverride.Service, optional) {#compactor-serviceoverrides}

See [Service override](../overrides/override/#service) 

Default: -

### metrics (*Metrics, optional) {#compactor-metrics}

Default: -

### httpAddress (string, optional) {#compactor-httpaddress}

Listen host:port for HTTP endpoints. 

Default: -

### httpGracePeriod (metav1.Duration, optional) {#compactor-httpgraceperiod}

Time to wait after an interrupt received for HTTP Server. 

Default: -

### dataDir (string, optional) {#compactor-datadir}

Data directory in which to cache blocks and process compactions. 

Default: -

### dataVolume (*volume.KubernetesVolume, optional) {#compactor-datavolume}

Kubernetes volume abstraction refers to different types of volumes to be mounted to pods: emptyDir, hostPath, pvc. 

Default: -

### consistencyDelay (metav1.Duration, optional) {#compactor-consistencydelay}

Minimum age of fresh (non-compacted) blocks before they are being processed. Malformed blocks older than the maximum of consistency-delay and 48h0m0s will be removed. 

Default: -

### retentionResolutionRaw (metav1.Duration, optional) {#compactor-retentionresolutionraw}

How long to retain raw samples in bucket. 0d - disables this retention. 

Default: -

### retentionResolution5m (metav1.Duration, optional) {#compactor-retentionresolution5m}

How long to retain samples of resolution 1 (5 minutes) in bucket. 0d - disables this retention. 

Default: -

### retentionResolution1h (metav1.Duration, optional) {#compactor-retentionresolution1h}

How long to retain samples of resolution 2 (1 hour) in bucket. 0d - disables this retention. 

Default: -

### wait (bool, optional) {#compactor-wait}

Do not exit after all compactions have been processed and wait for new work. 

Default: -

### downsamplingDisable (bool, optional) {#compactor-downsamplingdisable}

Disables downsampling. This is not recommended as querying long time ranges without non-downsampleddata is not efficient and useful e.g it is not possible to render all samples for a human eye anyway. 

Default: -

### blockSyncConcurrency (int, optional) {#compactor-blocksyncconcurrency}

Number of goroutines to use when syncing block metadata from object storage. 

Default: -

### compactConcurrency (int, optional) {#compactor-compactconcurrency}

Number of goroutines to use when compacting groups. 

Default: -


## BucketWeb

### metaOverrides (*typeoverride.ObjectMeta, optional) {#bucketweb-metaoverrides}

See [ObjectMeta override](../overrides/override/#objectmeta) 

Default: -

### deploymentOverrides (*typeoverride.Deployment, optional) {#bucketweb-deploymentoverrides}

See [Deployment override](../overrides/override/#deployment) 

Default: -

### serviceOverrides (*typeoverride.Service, optional) {#bucketweb-serviceoverrides}

See [Service override](../overrides/override/#service) 

Default: -

### metrics (*Metrics, optional) {#bucketweb-metrics}

Default: -

### HTTPIngress (*Ingress, optional) {#bucketweb-httpingress}

Default: -

### httpAddress (string, optional) {#bucketweb-httpaddress}

Listen host:port for HTTP endpoints. 

Default: -

### httpGracePeriod (metav1.Duration, optional) {#bucketweb-httpgraceperiod}

Time to wait after an interrupt received for HTTP Server. 

Default: -

### web_external_prefix (string, optional) {#bucketweb-web_external_prefix}

Static prefix for all HTML links and redirect URLs in the bucket web UI interface. Actual endpoints are still served on / or the web.route-prefix. This allows thanos bucket web UI to be served behind a reverse proxy that strips a URL sub-path. 

Default: -

### web_prefix_header (string, optional) {#bucketweb-web_prefix_header}

Name of HTTP request header used for dynamic prefixing of UI links and redirects. This option is ignored if web.external-prefix argument is set. Security risk: enable this option only if a reverse proxy in front of thanos is resetting the header. The --web.prefix-header=X-Forwarded-Prefix option can be useful, for example, if Thanos UI is served via Traefik reverse proxy with PathPrefixStrip option enabled, which sends the stripped prefix value in X-Forwarded-Prefix header. This allows thanos UI to be served on a sub-path. 

Default: -

### refresh (metav1.Duration, optional) {#bucketweb-refresh}

Refresh interval to download metadata from remote storage. 

Default: -

### timeout (metav1.Duration, optional) {#bucketweb-timeout}

Timeout to download metadata from remote. 

Default: -

### label (string, optional) {#bucketweb-label}

Prometheus label to use as timeline title. 

Default: -


## ObjectStoreStatus

ObjectStoreStatus defines the observed state of ObjectStore


## ObjectStore

ObjectStore is the Schema for the objectstores API

###  (metav1.TypeMeta, required) {#objectstore-}

Default: -

### metadata (metav1.ObjectMeta, optional) {#objectstore-metadata}

Default: -

### spec (ObjectStoreSpec, optional) {#objectstore-spec}

Default: -

### status (ObjectStoreStatus, optional) {#objectstore-status}

Default: -


## ObjectStoreList

ObjectStoreList contains a list of ObjectStore

###  (metav1.TypeMeta, required) {#objectstorelist-}

Default: -

### metadata (metav1.ListMeta, optional) {#objectstorelist-metadata}

Default: -

### items ([]ObjectStore, required) {#objectstorelist-items}

Default: -


