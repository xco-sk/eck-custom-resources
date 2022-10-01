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
  - [Space](docs/cr_space.md)
  - [Index pattern](docs/cr_index_pattern.md)
  - [Saved search](docs/cr_saved_search.md)
  - [Visualization](docs/cr_visualization.md)
  - [Lens](docs/cr_lens.md)
  - [Dashboard](docs/cr_dashboard.md)
  - [Data View](docs/cr_data_view.md)

## Installation

```shell
# Add eck-custom-resources helm repo
helm repo add eck-custom-resources https://xco-sk.github.io/eck-custom-resources/

# Install chart
helm install eck-cr eck-custom-resources/eck-custom-resources-operator
```
Configuration options are documented in [chart README file](charts/eck-custom-resources-operator/README.md)

## Upgrade guide

### From 0.3.2 to 0.4.0
There is new `DataView` CRD present. To apply the CRD, run:
```
kubectl apply --server-side -f https://raw.githubusercontent.com/xco-sk/eck-custom-resources/eck-custom-resources-operator-0.4.0/config/crd/bases/kibana.eck.github.com_dataviews.yaml
```


## Uninstallation
To uninstall the eck-cr from Kubernetes cluster, run:

```shell
helm uninstall eck-cr
```

This removes all resources related to eck-custom-resources operator. It won't remove the CRDs nor any deployed custom resource
(e.g. Index, Index Template ...), they will remain in K8s and also in Elasticsearch.

## Working with custom resources
After the operator is installed, you can deploy Elasticsearch/Kibana resources from the list above. The reconciler
will take care of propagating the change to Elasticsearch or Kibana, whether it is creation of new resource, deletion
or update.

For detailed documentation for each resource, see [List of supported resources](docs/cr_list.md)

## Help and Troubleshooting
In case you need help or found a bug, please create an [Issue on Github](https://github.com/xco-sk/eck-custom-resources/issues).

## License
Licensed under the Apache License, Version 2.0; see [LICENSE.md](LICENSE.md)
