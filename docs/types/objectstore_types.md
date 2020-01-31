### ObjectStoreSpec
#### ObjectStoreSpec defines the desired state of ObjectStore

| Variable Name | Type | Required | Default | Description |
|---|---|---|---|---|
| config | secret.Secret | Yes | - | Config<br> |
| compactor | *Compactor | No | - |  |
| bucketWeb | *BucketWeb | No | - |  |
### Compactor
| Variable Name | Type | Required | Default | Description |
|---|---|---|---|---|
|  | BaseObject | Yes | - |  |
| metrics | *Metrics | No | - |  |
| httpAddress | string | No | - | Listen host:port for HTTP endpoints.<br> |
| httpGracePeriod | metav1.Duration | No | - | Time to wait after an interrupt received for HTTP Server.<br> |
| dataDir | string | No | - | Data directory in which to cache blocks and process compactions.<br> |
| dataVolume | *volume.KubernetesVolume | No | - | Kubernetes volume abstraction refers to different types of volumes to be mounted to pods: emptyDir, hostPath, pvc.<br> |
| consistencyDelay | metav1.Duration | No | - | Minimum age of fresh (non-compacted) blocks before they are being processed.<br>Malformed blocks older than the maximum of consistency-delay and 48h0m0s will be removed.<br> |
| retentionResolutionRaw | metav1.Duration | No | - | How long to retain raw samples in bucket. 0d - disables this retention.<br> |
| retentionResolution5m | metav1.Duration | No | - | How long to retain samples of resolution 1 (5 minutes) in bucket. 0d - disables this retention.<br> |
| retentionResolution1h | metav1.Duration | No | - | How long to retain samples of resolution 2 (1 hour) in bucket. 0d - disables this retention.<br> |
| wait | bool | No | - | Do not exit after all compactions have been processed and wait for new work.<br> |
| downsamplingDisable | bool | No | - | Disables downsampling. This is not recommended as querying long time ranges without non-downsampleddata<br>is not efficient and useful e.g it is not possible to render all samples for a human eye anyway.<br> |
| blockSyncConcurrency | int | No | - | Number of goroutines to use when syncing block metadata from object storage.<br> |
| compactConcurrency | int | No | - | Number of goroutines to use when compacting groups.<br> |
### BucketWeb
| Variable Name | Type | Required | Default | Description |
|---|---|---|---|---|
|  | BaseObject | Yes | - |  |
| metrics | *Metrics | No | - |  |
| httpAddress | string | No | - | Listen host:port for HTTP endpoints.<br> |
| httpGracePeriod | metav1.Duration | No | - | Time to wait after an interrupt received for HTTP Server.<br> |
| web_external_prefix | string | No | - | Static prefix for all HTML links and redirect URLs in the bucket web UI interface. Actual endpoints are still served on / or the web.route-prefix. This allows thanos bucket web UI to be served behind a reverse proxy that strips a URL sub-path.<br> |
| web_prefix_header | string | No | - | Name of HTTP request header used for dynamic prefixing of UI links and redirects. This option is ignored if web.external-prefix argument is set. Security risk: enable this option only if a reverse proxy in front of thanos is resetting the header. The --web.prefix-header=X-Forwarded-Prefix option can be useful, for example, if Thanos UI is served via Traefik reverse proxy with PathPrefixStrip option enabled, which sends the stripped prefix value in X-Forwarded-Prefix header. This allows thanos UI to be served on a sub-path.<br> |
| refresh | metav1.Duration | No | - | Refresh interval to download metadata from remote storage.<br> |
| timeout | metav1.Duration | No | - | Timeout to download metadata from remote.<br> |
| label | string | No | - | Prometheus label to use as timeline title.<br> |
### ObjectStoreStatus
#### ObjectStoreStatus defines the observed state of ObjectStore

| Variable Name | Type | Required | Default | Description |
|---|---|---|---|---|
### ObjectStore
#### ObjectStore is the Schema for the objectstores API

| Variable Name | Type | Required | Default | Description |
|---|---|---|---|---|
|  | metav1.TypeMeta | Yes | - |  |
| metadata | metav1.ObjectMeta | No | - |  |
| spec | ObjectStoreSpec | No | - |  |
| status | ObjectStoreStatus | No | - |  |
### ObjectStoreList
#### ObjectStoreList contains a list of ObjectStore

| Variable Name | Type | Required | Default | Description |
|---|---|---|---|---|
|  | metav1.TypeMeta | Yes | - |  |
| metadata | metav1.ListMeta | No | - |  |
| items | []ObjectStore | Yes | - |  |
