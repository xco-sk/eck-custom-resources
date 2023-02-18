# Visualization (visualizations.kibana.eck.github.com)

Custom resource definition representing Visualization in Kibana.

## Lifecycle

Visualization lifecycle is simple - when it is deleted from K8s, it is also deleted from Kibana. Creation of
new resource is reconciled using `POST /api/saved_objects/visualization/` API. Update is done using
`PUT /api/saved_objects/visualization/`. In case the `spec.space` is filled in, the URLs are prefixed
with `/s/<spec.space>`.

See [Saved objects APIs](https://www.elastic.co/guide/en/kibana/master/saved-objects-api.html) in official documentation.

## Fields

| Key                         | Type            | Description                                                                                                                                     | Default                                              |
|-----------------------------|-----------------|-------------------------------------------------------------------------------------------------------------------------------------------------|------------------------------------------------------|
| `metadata.name`             | string          | Name of the Visualization, used also as its ID in Kibana                                                                                        | No default                                           |
| `spec.space`                | string          | Name of the Kibana namespace to which the Visualization is deployed to                                                                          | No default (will be deployed to "default" namespace) |
| `spec.targetInstance.name`| string         | Name of the [Kibana Instance](cr_kibana_instance.md) to which this Visualization will be deployed to | The operator configuration |
| `spec.body`                 | string          | Visualization definition json                                                                                                                   | No default                                           |
| `spec.dependencies`         | List of objects | List of dependencies - the reconciler will wait for all resources from the list to be present in Kibana before deploying/updating this resource | -                                                    |                                                 |
| `spec.dependencies[].space` | string          | Kibana Space where to look for given resource                                                                                                   | -                                                    |
| `spec.dependencies[].type`  | string          | Type of resource - one of `visualization, dashboard, search, index-pattern, lens`                                                               | -                                                    |
| `spec.dependencies[].name`  | string          | Name of resource                                                                                                                                | -                                                    |

## Example

```yaml
apiVersion: kibana.eck.github.com/v1alpha1
kind: Visualization
metadata:
  name: visualization-sample
spec:
  targetInstance:
    name: kibana-quickstart
  dependencies:
    - type: index-pattern
      name: indexpattern-sample
  body: |
    {
      "attributes": {
        "visState": "{\"title\":\"Count visualization\",\"type\":\"metric\",\"aggs\":[{\"id\":\"1\",\"enabled\":true,\"type\":\"count\",\"params\":{},\"schema\":\"metric\"}],\"params\":{\"addTooltip\":true,\"addLegend\":false,\"type\":\"metric\",\"metric\":{\"percentageMode\":false,\"useRanges\":false,\"colorSchema\":\"Green to Red\",\"metricColorMode\":\"None\",\"colorsRange\":[{\"from\":0,\"to\":10000}],\"labels\":{\"show\":true},\"invertColors\":false,\"style\":{\"bgFill\":\"#000\",\"bgColor\":false,\"labelColor\":false,\"subText\":\"\",\"fontSize\":60}}}}",
        "title": "Count visualization",
        "uiStateJSON": "{}",
        "description": "",
        "version": 1,
        "kibanaSavedObjectMeta": {
          "searchSourceJSON": "{\"query\":{\"query\":\"\",\"language\":\"kuery\"},\"filter\":[],\"indexRefName\":\"kibanaSavedObjectMeta.searchSourceJSON.index\"}"
        }
      },
      "references": [
        {
          "name": "kibanaSavedObjectMeta.searchSourceJSON.index",
          "type": "index-pattern",
          "id": "indexpattern-sample"
        }
      ]
    }
```
