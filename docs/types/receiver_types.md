## ReceiverSpec

### receiverGroups ([]ReceiverGroup, optional) {#receiverspec-receivergroups}

Default: -


## ReceiverGroup

ReceiverGroup defines a Receiver group
Tenants are the Hard tenants of the receiver group
Replicas are the number of instances in this receiver group

### name (string, required) {#receivergroup-name}

Default: -

### tenants ([]string, optional) {#receivergroup-tenants}

Default: -

### config (secret.Secret, required) {#receivergroup-config}

Default: -

### replicas (int32, optional) {#receivergroup-replicas}

Default: -

### metaOverrides (*typeoverride.ObjectMeta, optional) {#receivergroup-metaoverrides}

See [ObjectMeta override](../overrides/override/#objectmeta) 

Default: -

### statefulSetOverrides (*typeoverride.StatefulSet, optional) {#receivergroup-statefulsetoverrides}

See [StatefulSet override](../overrides/override/#statefulset) 

Default: -

### serviceOverrides (*typeoverride.Service, optional) {#receivergroup-serviceoverrides}

See [Service override](../overrides/override/#service) 

Default: -

### httpIngress (*Ingress, optional) {#receivergroup-httpingress}

Default: -

### httpServerCertificate (string, optional) {#receivergroup-httpservercertificate}

Secret name for HTTP Server certificate (Kubernetes TLS secret type) 

Default: -

### httpClientCertificate (string, optional) {#receivergroup-httpclientcertificate}

Secret name for HTTP Client certificate (Kubernetes TLS secret type) 

Default: -

### grpcIngress (*Ingress, optional) {#receivergroup-grpcingress}

Default: -

### htpcClientCertificate (string, optional) {#receivergroup-htpcclientcertificate}

Secret name for GRPC Server certificate (Kubernetes TLS secret type) 

Default: -

### grpcServerCertificate (string, optional) {#receivergroup-grpcservercertificate}

Secret name for GRPC Client certificate (Kubernetes TLS secret type) 

Default: -

### remoteWriteClientServerName (string, optional) {#receivergroup-remotewriteclientservername}

Server name to verify the hostname on the returned gRPC certificates. See https://tools.ietf.org/html/rfc4366#section-3.1 

Default: -

### metrics (*Metrics, optional) {#receivergroup-metrics}

Default: -

### httpAddress (string, optional) {#receivergroup-httpaddress}

Listen host:port for HTTP endpoints. 

Default: -

### httpGracePeriod (metav1.Duration, optional) {#receivergroup-httpgraceperiod}

Time to wait after an interrupt received for HTTP Server. 

Default: -

### grpcAddress (string, optional) {#receivergroup-grpcaddress}

Listen ip:port address for gRPC endpoints 

Default: -

### grpcGracePeriod (string, optional) {#receivergroup-grpcgraceperiod}

Time to wait after an interrupt received for GRPC Server. 

Default: -

### remoteWriteAddress (string, optional) {#receivergroup-remotewriteaddress}

Address to listen on for remote write requests. 

Default: -

### labels (map[string]string, optional) {#receivergroup-labels}

External labels to announce. This flag will be removed in the future when handling multiple tsdb instances is added. 

Default: -

### dataVolume (*volume.KubernetesVolume, optional) {#receivergroup-datavolume}

Kubernetes volume abstraction refers to different types of volumes to be mounted to pods: emptyDir, hostPath, pvc. 

Default: -

### tsdbPath (string, optional) {#receivergroup-tsdbpath}

Default: -

### tsdbRetention (string, optional) {#receivergroup-tsdbretention}

How long to retain raw samples on local storage. 0d - disables this retention. 

Default: -

### tsdbMinBlockDuration (string, optional) {#receivergroup-tsdbminblockduration}

The --tsdb.min-block-duration and --tsdb.max-block-duration must be set to equal values to disable local compaction on order to use Thanos sidecar upload. Leave local compaction on if sidecar just exposes StoreAPI and your retention is normal. 

Default: -

### tsdbMaxBlockDuration (string, optional) {#receivergroup-tsdbmaxblockduration}

Default: -

### receiveHashringsFileRefreshInterval (string, optional) {#receivergroup-receivehashringsfilerefreshinterval}

Refresh interval to re-read the hashring configuration file. (used as a fallback) 

Default: -

### receiveTenantHeader (string, optional) {#receivergroup-receivetenantheader}

HTTP header to determine tenant for write requests. 

Default: -

### receiveDefaultTenantId (string, optional) {#receivergroup-receivedefaulttenantid}

Default tenant ID to use when none is provided via a header. 

Default: -

### receiveTenantLabelName (string, optional) {#receivergroup-receivetenantlabelname}

Label name through which the tenant will be announced. 

Default: -

### receiveReplicaHeader (string, optional) {#receivergroup-receivereplicaheader}

HTTP header specifying the replica number of a write request. 

Default: -

### receiveReplicationFactor (int, optional) {#receivergroup-receivereplicationfactor}

How many times to replicate incoming write requests. 

Default: -

### tsdbWalCompression (*bool, optional) {#receivergroup-tsdbwalcompression}

Compress the tsdb WAL. 

Default: -

### tsdbNoLockfile (*bool, optional) {#receivergroup-tsdbnolockfile}

Do not create lockfile in TSDB data directory. In any case, the lockfiles will be deleted on next startup. 

Default: -


## ReceiverStatus

ObjectStoreStatus defines the observed state of ObjectStore


## Receiver

Receiver is the Schema for the receiver cluster

###  (metav1.TypeMeta, required) {#receiver-}

Default: -

### metadata (metav1.ObjectMeta, optional) {#receiver-metadata}

Default: -

### spec (ReceiverSpec, optional) {#receiver-spec}

Default: -

### status (ReceiverStatus, optional) {#receiver-status}

Default: -


## ReceiverList

ObjectStoreList contains a list of ObjectStore

###  (metav1.TypeMeta, required) {#receiverlist-}

Default: -

### metadata (metav1.ListMeta, optional) {#receiverlist-metadata}

Default: -

### items ([]Receiver, required) {#receiverlist-items}

Default: -


