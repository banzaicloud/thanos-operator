<p align="center"><img src="../img/logo/thanos_operator_vertical.svg" width="260"></p>
<p align="center">


# Single Cluster Thanos Install

<p align="center"><img src="../img/Thanos-single-cluster.png" ></p>


## Prerequisites for Thanos


## Object Store secret

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

## Use Thanos with Prometheus Operator
Extra configuration for prometheus operator.

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
    externalLabels: 
      cluster: thanos-operator-test
```

Remember to set `externalLabels` as it identifies the Prometheus instance for Thanos.


### Add kubernetes stable helm repository
```
helm repo add stable https://kubernetes-charts.storage.googleapis.com
helm repo update
```

### Install prometheus-operator
```
helm install monitor stable/prometheus-operator -f thanos-sidecar.yaml --set manageCrds=false
```

### Install thanos-operator
```
helm install thanos-operator  ./charts/thanos-operator --set manageCrds=false
```
