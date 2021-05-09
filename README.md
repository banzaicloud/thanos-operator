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

Thanos Operator is a Kubernetes operator to manage Thanos stack deployment
on Kubernetes.

## What is [Thanos](http://thanos.io)

Open source, highly available Prometheus setup with long term storage capabilities.


## Architecture
<p align="center"><img src="docs/img/thanos-single-cluster2.png" ></p>

## Feature highlights

- Auto discover endpoints
- Manage persistent volumes
- Metrics configuration
- Simple TLS configuration

Work in progress

- Tracing configuration
- Endpoint validation
- Certificate management
- Advanced secret configuration

## Documentation

 You can find the complete documentation of thanos operator [here](./docs/README.md) :blue_book: <br>

## Commercial support
If you are using the Thanos operator in a production environment and [require commercial support, contact Banzai Cloud](https://banzaicloud.com/contact/), the company backing the development of the Thanos operator. If you are looking for the ultimate observability tool for multi-cluster Kubernetes infrastructures to automate the collection, correlation, and storage of logs and metrics, check out [One Eye](https://banzaicloud.com/products/one-eye/).


## Contributing

If you find this project useful, help us:

- Support the development of this project and star this repo! :star:
- If you use the Thanos operator in a production environment, add yourself to the list of production [adopters](https://github.com/banzaicloud/thanos-operator/blob/master/ADOPTERS.md).:metal: <br> 
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
