# Example deployment modes for Thanos

## Custom resources

TODO link CRD documentation

### Thanos
 - Query
 - Store
 - Rule 
 
### ObjectStore
 - Compactor
 - Bucket

### StoreEndpoint
 - Sidecar
 - StoreAPI addresses

## Single Cluster Install

#### Install secret

Example S3 configuration
```
type: S3
config:
  endpoint: "s3.eu-west-1.amazonaws.com"
  bucket: "test-bucket"
  region: "eu-west-1"
  access_key: "XXXXXXXXX"
  secret_key: "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
```

Deploy the secret on Kubernetes
```
kubectl create secret generic thanos --from-file=object-store.yaml=object-store.yaml
```

#### Install Prometheus Operator

> Note: Prometheus-operator and Thanos MUST be in the same namespace.

*thanos-sidecar.yaml*
```
prometheus:
  prometheusSpec:
    thanos:
      image: quay.io/thanos/thanos:v0.9.0
      version: v0.9.0
      objectStorageConfig:
        name: thanos
        key: object-store.yaml
    externalLabels: thanos-operator-test
```

Remember to set `externalLabels` as it identifies the Prometheus instance for Thanos.

#### Install prometheus-operator
```
helm install --name monitor stable/prometheus-operator -f thanos-sidecar.yaml
```

#### Install Thanos Operator

Thanos
```
apiVersion: monitoring.banzaicloud.io/v1alpha1
kind: Thanos
metadata:
  name: thanos-sample
spec:
  query: {}
  rule: {}
  storeGateway: {}
```

ObjectStore
```
apiVersion: monitoring.banzaicloud.io/v1alpha1
kind: ObjectStore
metadata:
  name: objectstore-sample
spec:
  config:
    mountFrom:
      secretKeyRef:
        name: thanos
        key: object-store.yaml
  bucketWeb: {}
  compactor: {}
```

StoreEndpoint
```
apiVersion: monitoring.banzaicloud.io/v1alpha1
kind: StoreEndpoint
metadata:
  name: storeendpoint-sample
spec:
  thanos: thanos-sample
  config:
    mountFrom:
      secretKeyRef:
        name: thanos
        key: object-store.yaml
  selector: {}
```

### Multiple Prometheus on single cluster

You can define different Prometheuses and Endpoints via `StoreEndpoint` CRs.

## Query discovery
Automatically discover all query on the cluster

```
apiVersion: monitoring.banzaicloud.io/v1alpha1
kind: Thanos
metadata:
  name: query
spec:
  query: {}
  queryDiscovery: true
```

## Remote Thanos for Prometheus

- Remote URLs
  - format: `http(s)://<fqdn>:<port>`
  - tls: ?

## Multi Thanos Observer

- Select endpoints based on labels
  - Namespace?
  - Labels?
  - Pods/Svc?

## Sidecar SD
As sidecars is only useful for Query it should be under query key
sidecars:
  - namespaces (optional) | []string | default: same namespace
    labels (optional) | default: app=prometheus
    url (optional) | default: ""
    
Based on `namespaces` and `labels` the operator creates services and adds
them to query parameters.

The `url` attribute has priority and requires FQDN