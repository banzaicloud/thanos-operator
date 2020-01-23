# Thanos Operator

Thanos Operator is a Kubernetes operator to manage Thanos stack seployment
on Kubernetes.

## What is [Thanos](http://thanos.io)

Open source, highly available Prometheus setup with long term storage capabilities.

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
    externalLabels: thanos-operator-test
```

Remember to set `externalLabels` as it identifies the Prometheus instance for Thanos.

### Install prometheus-operator
```
helm install --name monitor stable/prometheus-operator -f thanos-sidecar.yaml
```