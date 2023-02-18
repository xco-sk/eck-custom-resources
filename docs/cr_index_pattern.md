# Index Pattern (indexpatterns.kibana.eck.github.com)

Custom resource definition representing Index Pattern in Kibana.

## Lifecycle

Index pattern lifecycle is simple - when the pattern is deleted from K8s, it is also deleted from Kibana. Creation of
new resource is reconciled using `POST /api/saved_objects/index-pattern/` API. Update is done using
`PUT /api/saved_objects/index-pattern/`.

See [Index patterns APIs](https://www.elastic.co/guide/en/kibana/8.2/index-patterns-api.html) in official documentation.

## Fields

| Key                         | Type            | Description                                                                                                                                     | Default                                              |
|-----------------------------|-----------------|-------------------------------------------------------------------------------------------------------------------------------------------------|------------------------------------------------------|
| `metadata.name`             | string          | Name of the Index Pattern, used also as its ID in Kibana                                                                                        | No default                                           |
| `spec.space`                | string          | Name of the Kibana namespace to which the Index pattern is deployed to                                                                          | No default (will be deployed to "default" namespace) |
| `spec.targetInstance.name`  | string         | Name of the [Kibana Instance](cr_kibana_instance.md) to which this IndexPattern will be deployed to | The operator configuration |
| `spec.body`                 | string          | Index pattern definition json                                                                                                                   | No default                                           |
| `spec.dependencies`         | List of objects | List of dependencies - the reconciler will wait for all resources from the list to be present in Kibana before deploying/updating this resource | -                                                    |                                                 |
| `spec.dependencies[].space` | string          | Kibana Space where to look for given resource                                                                                                   | -                                                    |
| `spec.dependencies[].type`  | string          | Type of resource - one of `visualization, dashboard, search, index-pattern, lens`                                                               | -                                                    |
| `spec.dependencies[].name`  | string          | Name of resource                                                                                                                                | -                                                    |

## Example

```yaml
apiVersion: kibana.eck.github.com/v1alpha1
kind: IndexPattern
metadata:
  name: indexpattern-sample
spec:
  targetInstance:
    name: kibana-quickstart
  space: my-space
  dependencies:
    - type: index-pattern
      name: indexpattern-base-sample
  body: |
    {
        "attributes": {
            "title": "index-*",
            "timeFieldName": "@timestamp"
        }
    }
```
