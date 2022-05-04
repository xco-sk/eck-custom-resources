# Custom resources for ECK
[![docker-publish](https://github.com/xco-sk/eck-custom-resources/actions/workflows/docker-publish.yaml/badge.svg)](https://github.com/xco-sk/eck-custom-resources/actions/workflows/docker-publish.yaml)
[![helm-publish](https://github.com/xco-sk/eck-custom-resources/actions/workflows/helm-publish.yml/badge.svg)](https://github.com/xco-sk/eck-custom-resources/actions/workflows/helm-publish.yml)

Kubernetes operator that enables the installation of various resources for
Elasticsearch and Kibana.

Currently supported resources: 
- For Elasticsearch:
  - [Index](docs/cr_index.md)
  - [Index template](docs/cr_index_template.md)
  - [Index lifecycle policy](docs/cr_index_lifecycle_policy.md)
  - [Ingest pipeline](docs/cr_ingest_pipeline.md)
  - [Snapshot repository](docs/cr_snapshot_repo.md)
  - [Snapshot lifecycle policy](docs/cr_snapshot_lifecycle_policy.md)
  - [User](docs/cr_user.md)
  - [Role](docs/cr_role.md)
- For Kibana:
  - [Index pattern](docs/cr_index_pattern.md)
  - [Saved search](docs/cr_saved_search.md)
  - [Visualization](docs/cr_visualization.md)
  - [Dashboard](docs/cr_dashboard.md)

# Installation

```shell
# Add eck-custom-resources helm repo
helm repo add eck-custom-resources https://xco-sk.github.io/eck-custom-resources/

# Install chart
helm install eck-cr eck-custom-resources/eck-custom-resources-operator
```
Configuration options are documented in [docs/helm](docs/helm.md)



