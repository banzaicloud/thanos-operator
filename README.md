<p align="center"><img src="docs/img/logo/thanos_operator_vertical.svg" width="260"></p>
<p align="center">

  <a href="https://hub.docker.com/r/banzaicloud/thanos-operator/">
    <img src="https://img.shields.io/docker/automated/banzaicloud/thanos-operator.svg" alt="Docker Automated build">
  </a>

  <a href="https://hub.docker.com/r/banzaicloud/thanos-operator/">
    <img src="https://img.shields.io/docker/pulls/banzaicloud/thanos-operator.svg?style=shield" alt="Docker Pulls">
  </a>

  <a href="https://circleci.com/gh/banzaicloud/thanos-operator">
    <img src="https://circleci.com/gh/banzaicloud/thanos-operator.svg?style=shield" alt="CircleCI">
  </a>

  <a href="https://goreportcard.com/badge/github.com/banzaicloud/thanos-operator">
    <img src="https://goreportcard.com/badge/github.com/banzaicloud/thanos-operator" alt="Go Report Card">
  </a>

  <a href="https://github.com/banzaicloud/thanos-operator/">
    <img src="https://img.shields.io/badge/license-Apache%20v2-orange.svg" alt="license">
  </a>

</p>

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


## Documentation

 You can find the complete documentation of thanos operator [here](./docs/Readme.md) :blue_book: <br>

## Commercial support
If you are using the thanos operator in a production environment and [require commercial support, contact Banzai Cloud](https://banzaicloud.com/contact/), the company backing the development of the thanos operator.


## Contributing

If you find this project useful, help us:

- Support the development of this project and star this repo! :star:
- If you use the Logging operator in a production environment, add yourself to the list of production [adopters](https://github.com/banzaicloud/thanos-operator/blob/master/ADOPTERS.md).:metal: <br> 
- Help new users with issues they may encounter :muscle:
- Send a pull request with your new features and bug fixes :rocket: 

*For more information, read the [developer documentation](./docs/developers.md)*.

## License

Copyright (c) 2017-2020 [Banzai Cloud, Inc.](https://banzaicloud.com)

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

[http://www.apache.org/licenses/LICENSE-2.0](http://www.apache.org/licenses/LICENSE-2.0)

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
