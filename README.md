# Custom resources for ECK
[![docker-publish](https://github.com/xco-sk/eck-custom-resources/actions/workflows/docker-publish.yaml/badge.svg)](https://github.com/xco-sk/eck-custom-resources/actions/workflows/docker-publish.yaml)
[![helm-publish](https://github.com/xco-sk/eck-custom-resources/actions/workflows/helm-publish.yml/badge.svg)](https://github.com/xco-sk/eck-custom-resources/actions/workflows/helm-publish.yml)

Kubernetes operator that enables the installation of various resources for
Elasticsearch and Kibana.

Currently supported resources: 
- For Elasticsearch:
  - Index
  - Index template
  - Index lifecycle policy
  - Ingest pipeline
  - Snapshot repository
  - Snapshot lifecycle policy
  - User
  - Role
- For Kibana:
  - Index pattern
  - Saved search
  - Visualization
  - Dashboard

# Installation

```shell
# Add eck-custom-resources helm repo
helm repo add eck-custom-resources https://xco-sk.github.io/eck-custom-resources/

# Install chart
helm install eck-cr eck-custom-resources/eck-custom-resources-operator
```
Configuration options are documented in [docs/helm](docs/helm.md)



