<p align="center"><img src="./img/logo/thanos_operator_vertical.svg" width="260"></p>
<p align="center">

# Example deployment modes for Thanos

## Custom resources

The API reference of CRDs are documented [here](types/Readme.md)

# Single Cluster Install

In this scenario we deploy a full featured Prometheus with Thanos as long term storage.

## Prerequisites for Thanos

### Install secret

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

### Install Prometheus Operator

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

### Install prometheus-operator
```
helm install --name monitor stable/prometheus-operator -f thanos-sidecar.yaml
```

## Install Thanos Operator

For now the kustomize deploy model is available

```
make install
make deploy IMG=banzaicloud/thanos-operator:latest
```

### Apply CRDs for single cluster
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
Automatically discover all Query instances (created by CRD) on the cluster

```
apiVersion: monitoring.banzaicloud.io/v1alpha1
kind: Thanos
metadata:
  name: query
spec:
  query: {}
  queryDiscovery: true
```

## Remote Prometheus for Thanos
Remote URLs
  - format: `http(s)://<fqdn>:<port>`

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
  url: http://example.com:10901
```

## Multi Thanos Observer

- TODO