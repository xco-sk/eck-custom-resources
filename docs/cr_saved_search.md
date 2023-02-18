# Saved Search (savedsearches.kibana.eck.github.com)

Custom resource definition representing (Saved) Search in Kibana.

## Lifecycle

Saved search lifecycle is simple - when the search is deleted from K8s, it is also deleted from Kibana. Creation of
new resource is reconciled using `POST /api/saved_objects/search/` API. Update is done using
`PUT /api/saved_objects/search/`. In case the `spec.space` is filled in, the URLs are prefixed
with `/s/<spec.space>`.

See [Saved objects APIs](https://www.elastic.co/guide/en/kibana/master/saved-objects-api.html) in official documentation.

## Fields

| Key                         | Type            | Description                                                                                                                                     | Default                                              |
|-----------------------------|-----------------|-------------------------------------------------------------------------------------------------------------------------------------------------|------------------------------------------------------|
| `metadata.name`             | string          | Name of the Saved search, used also as its ID in Kibana                                                                                         | No default                                           |
| `spec.space`                | string          | Name of the Kibana namespace to which the Search is deployed to                                                                                 | No default (will be deployed to "default" namespace) |
| `spec.targetInstance.name`  | string         | Name of the [Kibana Instance](cr_kibana_instance.md) to which this SavedSearch will be deployed to | The operator configuration |
| `spec.body`                 | string          | Saved search definition json                                                                                                                    | No default                                           |
| `spec.dependencies`         | List of objects | List of dependencies - the reconciler will wait for all resources from the list to be present in Kibana before deploying/updating this resource | -                                                    |                                                 |
| `spec.dependencies[].space` | string          | Kibana Space where to look for given resource                                                                                                   | -                                                    |
| `spec.dependencies[].type`  | string          | Type of resource - one of `visualization, dashboard, search, index-pattern, lens`                                                               | -                                                    |
| `spec.dependencies[].name`  | string          | Name of resource                                                                                                                                | -                                                    |

## Example

```yaml
apiVersion: kibana.eck.github.com/v1alpha1
kind: SavedSearch
metadata:
  name: savedsearch-sample
spec:
  targetInstance:
    name: kibana-quickstart
  dependencies:
    - type: index-pattern
      name: indexpattern-sample
  body: |
    {
      "attributes": {
         "columns": [
            "_source"
         ],
         "description": "",
         "hits": 0,
         "kibanaSavedObjectMeta": {
            "searchSourceJSON": "{\"highlightAll\":true,\"version\":true,\"query\":{\"query\":\"(message: test)\",\"language\":\"kuery\"},\"filter\":[],\"indexRefName\":\"kibanaSavedObjectMeta.searchSourceJSON.index\"}"
         },
         "sort": [],
         "title": "Message contains test",
         "version": 1
      },
      "references": [
         {
            "id": "indexpattern-sample",
            "name": "kibanaSavedObjectMeta.searchSourceJSON.index",
            "type": "index-pattern"
         }
      ]
    }
```
