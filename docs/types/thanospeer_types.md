## ThanosPeer

###  (metav1.TypeMeta, required) {#thanospeer-}

Default: -

### metadata (metav1.ObjectMeta, optional) {#thanospeer-metadata}

Default: -

### spec (ThanosPeerSpec, optional) {#thanospeer-spec}

See [ThanosPeerSpec](#thanospeerspec) 

Default: -

### status (ThanosPeerStatus, optional) {#thanospeer-status}

See [ThanosPeerStatus](#thanospeerstatus) 

Default: -


## ThanosPeerList

###  (metav1.TypeMeta, required) {#thanospeerlist-}

Default: -

### metadata (metav1.ListMeta, optional) {#thanospeerlist-metadata}

Default: -

### items ([]ThanosPeer, required) {#thanospeerlist-items}

Default: -


## ThanosPeerSpec

### endpointAddress (string, required) {#thanospeerspec-endpointaddress}

Host (or IP) and port of the remote Thanos endpoint 

Default: -

### peerEndpointAlias (string, optional) {#thanospeerspec-peerendpointalias}

Optional alias for the remote endpoint in case we have to access it through a different name. This is typically needed if the remote endpoint has a certificate created for a predefined hostname. The controller should create an externalName service for this backed buy the actual peer endpoint host or a k8s service with a manually crafted k8s endpoint if EndpointAddress doesn't have a host but only an IP. 

Default: -

### certificate (string, optional) {#thanospeerspec-certificate}

The peer query should use this client certificate (tls.crt, tls.key) in the current namespace 

Default: -

### caBundle (string, optional) {#thanospeerspec-cabundle}

Name of the secret that contains the CA certificate in ca.crt to verify client certs in the current namespace 

Default: -

### replicaLabels ([]string, optional) {#thanospeerspec-replicalabels}

Custom replica labels if the default doesn't apply 

Default: -

### metaOverrides (typeoverride.ObjectMeta, optional) {#thanospeerspec-metaoverrides}

[Override metadata](../overrides/override/#objectmeta) for managed resources 

Default: -

### queryOverrides (*Query, optional) {#thanospeerspec-queryoverrides}

Override any of the [Query parameters](../thanos_types/#query) 

Default: -


## ThanosPeerStatus

### queryHTTPServiceURL (string, optional) {#thanospeerstatus-queryhttpserviceurl}

The peer query is available over HTTP on this internal service URL 

Default: -


