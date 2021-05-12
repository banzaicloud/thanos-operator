### ReceiverSpec
| Variable Name | Type | Required | Default | Description |
|---|---|---|---|---|
| receiverGroups | []ReceiverGroup | No | - |  |
### ReceiverGroup
#### ReceiverGroup defines a Receiver group
Tenants are the Hard tenants of the receiver group
Replicas are the number of instances in this receiver group

| Variable Name | Type | Required | Default | Description |
|---|---|---|---|---|
| name | string | Yes | - |  |
| tenants | []string | No | - |  |
| config | secret.Secret | Yes | - |  |
| replicas | int32 | No | - |  |
| metaOverrides | *typeoverride.ObjectMeta | No | - | See [ObjectMeta override](../overrides/override/#objectmeta)<br> |
| statefulSetOverrides | *typeoverride.StatefulSet | No | - | See [StatefulSet override](../overrides/override/#statefulset)<br> |
| serviceOverrides | *typeoverride.Service | No | - | See [Service override](../overrides/override/#service)<br> |
| httpIngress | *Ingress | No | - |  |
| httpServerCertificate | string | No | - | Secret name for HTTP Server certificate (Kubernetes TLS secret type)<br> |
| httpClientCertificate | string | No | - | Secret name for HTTP Client certificate (Kubernetes TLS secret type)<br> |
| grpcIngress | *Ingress | No | - |  |
| htpcClientCertificate | string | No | - | Secret name for GRPC Server certificate (Kubernetes TLS secret type)<br> |
| grpcServerCertificate | string | No | - | Secret name for GRPC Client certificate (Kubernetes TLS secret type)<br> |
| remoteWriteClientServerName | string | No | - | Server name to verify the hostname on the returned gRPC certificates. See https://tools.ietf.org/html/rfc4366#section-3.1<br> |
| metrics | *Metrics | No | - |  |
| httpAddress | string | No | - | Listen host:port for HTTP endpoints.<br> |
| httpGracePeriod | metav1.Duration | No | - | Time to wait after an interrupt received for HTTP Server.<br> |
| grpcAddress | string | No | - | Listen ip:port address for gRPC endpoints<br> |
| grpcGracePeriod | string | No | - | Time to wait after an interrupt received for GRPC Server.<br> |
| remoteWriteAddress | string | No | - | Address to listen on for remote write requests.<br> |
| labels | map[string]string | No | - | External labels to announce. This flag will be removed in the future when handling multiple tsdb instances is added.<br> |
| dataVolume | *volume.KubernetesVolume | No | - | Kubernetes volume abstraction refers to different types of volumes to be mounted to pods: emptyDir, hostPath, pvc.<br> |
| tsdbPath | string | No | - |  |
| tsdbRetention | string | No | - | How long to retain raw samples on local storage. 0d - disables this retention.<br> |
| tsdbMinBlockDuration | string | No | - | The --tsdb.min-block-duration and --tsdb.max-block-duration must be set to equal values to disable local compaction<br>on order to use Thanos sidecar upload. Leave local compaction on if sidecar just exposes StoreAPI and your retention is normal.<br> |
| tsdbMaxBlockDuration | string | No | - |  |
| receiveHashringsFileRefreshInterval | string | No | - | Refresh interval to re-read the hashring configuration file. (used as a fallback)<br> |
| receiveTenantHeader | string | No | - | HTTP header to determine tenant for write requests.<br> |
| receiveDefaultTenantId | string | No | - | Default tenant ID to use when none is provided via a header.<br> |
| receiveTenantLabelName | string | No | - | Label name through which the tenant will be announced.<br> |
| receiveReplicaHeader | string | No | - | HTTP header specifying the replica number of a write request.<br> |
| receiveReplicationFactor | int | No | - | How many times to replicate incoming write requests.<br> |
| tsdbWalCompression | *bool | No | - | Compress the tsdb WAL.<br> |
| tsdbNoLockfile | *bool | No | - | Do not create lockfile in TSDB data directory. In any case, the lockfiles will be deleted on next startup.<br> |
### ReceiverStatus
#### ObjectStoreStatus defines the observed state of ObjectStore

| Variable Name | Type | Required | Default | Description |
|---|---|---|---|---|
### Receiver
#### Receiver is the Schema for the receiver cluster

| Variable Name | Type | Required | Default | Description |
|---|---|---|---|---|
|  | metav1.TypeMeta | Yes | - |  |
| metadata | metav1.ObjectMeta | No | - |  |
| spec | ReceiverSpec | No | - |  |
| status | ReceiverStatus | No | - |  |
### ReceiverList
#### ObjectStoreList contains a list of ObjectStore

| Variable Name | Type | Required | Default | Description |
|---|---|---|---|---|
|  | metav1.TypeMeta | Yes | - |  |
| metadata | metav1.ListMeta | No | - |  |
| items | []Receiver | Yes | - |  |
