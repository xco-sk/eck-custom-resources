# Data View (dataviews.kibana.eck.github.com)

Custom resource definition representing Data View resource in Kibana.

## Lifecycle

Data view lifecycle is simple - when it is deleted from K8s, it is also deleted from Kibana. Creation of
new resource is reconciled using `POST /api/data_views/data_view/` API. In case the `spec.space` is filled in, the URLs are prefixed
with `/s/<spec.space>`. There is however one field not being updated - `name`. The DataView update API does not allow to modify this field 
(it is actually not allowed to be part of the request body). To overcome this limitation, the reconciler removes this field on update, that mean
there might be inconsistency in DataView in K8s and actual deployed DataView in Kibana.

See [Data Views APIs](https://www.elastic.co/guide/en/kibana/current/data-views-api.html) in official documentation.

## Fields

| Key                         | Type            | Description                                                                                                                                     | Default                                              |
|-----------------------------|-----------------|-------------------------------------------------------------------------------------------------------------------------------------------------|------------------------------------------------------|
| `metadata.name`             | string          | Name of the Data View visualization, used also as its ID in Kibana                                                                                   | No default                                           |
| `spec.space`                | string          | Name of the Kibana namespace to which the Data View is deployed to                                                                                   | No default (will be deployed to "default" namespace) |
| `spec.targetInstance.name`  | string         | Name of the [Kibana Instance](cr_kibana_instance.md) to which this DataView will be deployed to | The operator configuration |
| `spec.body`                 | string          | Data View definition (the inner part of the requests) json                                                                                                                            | No default                                           |
| `spec.dependencies`         | List of objects | List of dependencies - the reconciler will wait for all resources from the list to be present in Kibana before deploying/updating this resource | -                                                    |                                                 |
| `spec.dependencies[].space` | string          | Kibana Space where to look for given resource                                                                                                   | -                                                    |
| `spec.dependencies[].type`  | string          | Type of resource - one of `visualization, dashboard, search, index-pattern, lens`                                                               | -                                                    |
| `spec.dependencies[].name`  | string          | Name of resource                                                                                                                                | -                                                    |

## Example

```yaml
apiVersion: kibana.eck.github.com/v1alpha1
kind: DataView
metadata:
  name: dataview-sample
spec:
  targetInstance:
    name: kibana-quickstart
  space: space-sample
  dependencies:
    - type: lens
      name: lens-sample
  body: |
    {
      "title": "sample-index-*",
      "type": "rollup",
      "typeMeta": {
        "params": {
          "rollup_index": "rollup_logstash"
        },
        "aggs": {
          "terms": {
            "geo.dest": { "agg": "terms" },
            "extension.keyword": { "agg": "terms" },
            "geo.src": { "agg": "terms" },
            "machine.os.keyword": { "agg": "terms" }
          },
          "date_histogram": {
            "@timestamp": {
              "agg": "date_histogram",
              "fixed_interval": "20m",
              "delay": "10m",
              "time_zone": "UTC"
            }
          },
          "avg": {
            "memory": { "agg": "avg" },
            "bytes": { "agg": "avg" }
          },
          "max": { "memory": { "agg": "max" } },
          "min": { "memory": { "agg": "min" } },
          "sum": { "memory": { "agg": "sum" } },
          "value_count": { "memory": { "agg": "value_count" } },
          "histogram": {
            "machine.ram": {
              "agg": "histogram",
              "interval": 5
            }
          }
        }
      }
    }
```
