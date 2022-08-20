# Data View (dataviews.kibana.eck.github.com)

Custom resource definition representing Data View resource in Kibana.

## Lifecycle

Data view lifecycle is simple - when it is deleted from K8s, it is also deleted from Kibana. Creation of
new resource is reconciled using `POST /api/data_views/data_view/` API. In case the `spec.space` is filled in, the URLs are prefixed
with `/s/<spec.space>`.

See [Data Views APIs](https://www.elastic.co/guide/en/kibana/current/data-views-api.html) in official documentation.

## Fields

| Key                         | Type            | Description                                                                                                                                     | Default                                              |
|-----------------------------|-----------------|-------------------------------------------------------------------------------------------------------------------------------------------------|------------------------------------------------------|
| `metadata.name`             | string          | Name of the Data View visualization, used also as its ID in Kibana                                                                                   | No default                                           |
| `spec.space`                | string          | Name of the Kibana namespace to which the Data View is deployed to                                                                                   | No default (will be deployed to "default" namespace) |
| `spec.body`                 | string          | Data View definition json                                                                                                                            | No default                                           |
| `spec.dependencies`         | List of objects | List of dependencies - the reconciler will wait for all resources from the list to be present in Kibana before deploying/updating this resource | -                                                    |                                                 |
| `spec.dependencies[].space` | string          | Kibana Space where to look for given resource                                                                                                   | -                                                    |
| `spec.dependencies[].type`  | string          | Type of resource - one of `visualization, dashboard, search, index-pattern, Data View`                                                               | -                                                    |
| `spec.dependencies[].name`  | string          | Name of resource                                                                                                                                | -                                                    |

## Example

```yaml
apiVersion: kibana.eck.github.com/v1alpha1
kind: DataView
metadata:
  name: dataview-sample
spec:
  space: my-space
  dependencies:
    - type: lens
      name: lens-sample
  body: |
    {
      "override": true,
      "refresh_fields": true,
      "data_view": {
         "title": "sample-view"
      }
    }
```
