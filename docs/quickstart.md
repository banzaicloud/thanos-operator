# Single Cluster Thanos Install

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
cat << EOF > thanos-sidecar.yaml
prometheus:
  prometheusSpec:
    thanos:
      image: quay.io/thanos/thanos:v0.10.1
      version: v0.10.1
      objectStorageConfig:
        name: thanos
        key: object-store.yaml
    externalLabels:
      cluster: thanos-operator-test-1
EOF
```

Remember to set `externalLabels` as it identifies the Prometheus instance for Thanos.

### Install prometheus-operator
```
helm3 install monitor stable/prometheus-operator -f thanos-sidecar.yaml
```