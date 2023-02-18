# Lens visualization (lens.kibana.eck.github.com)

Custom resource definition representing Lens visualization in Kibana.

## Lifecycle

Lens lifecycle is simple - when it is deleted from K8s, it is also deleted from Kibana. Creation of
new resource is reconciled using `POST /api/saved_objects/visualization/` API. Update is done using
`PUT /api/saved_objects/visualization/`. In case the `spec.space` is filled in, the URLs are prefixed
with `/s/<spec.space>`.

See [Saved objects APIs](https://www.elastic.co/guide/en/kibana/master/saved-objects-api.html) in official documentation.

## Fields

| Key                         | Type            | Description                                                                                                                                     | Default                                              |
|-----------------------------|-----------------|-------------------------------------------------------------------------------------------------------------------------------------------------|------------------------------------------------------|
| `metadata.name`             | string          | Name of the Lens visualization, used also as its ID in Kibana                                                                                   | No default                                           |
| `spec.space`                | string          | Name of the Kibana namespace to which the Lens is deployed to                                                                                   | No default (will be deployed to "default" namespace) |
| `spec.targetInstance.name`  | string         | Name of the [Kibana Instance](cr_kibana_instance.md) to which this Lens will be deployed to | The operator configuration |
| `spec.body`                 | string          | Lens definition json                                                                                                                            | No default                                           |
| `spec.dependencies`         | List of objects | List of dependencies - the reconciler will wait for all resources from the list to be present in Kibana before deploying/updating this resource | -                                                    |                                                 |
| `spec.dependencies[].space` | string          | Kibana Space where to look for given resource                                                                                                   | -                                                    |
| `spec.dependencies[].type`  | string          | Type of resource - one of `visualization, dashboard, search, index-pattern, lens`                                                               | -                                                    |
| `spec.dependencies[].name`  | string          | Name of resource                                                                                                                                | -                                                    |

## Example

```yaml
apiVersion: kibana.eck.github.com/v1alpha1
kind: Lens
metadata:
  name: lens-sample
spec:
  targetInstance:
    name: kibana-quickstart
  space: my-space
  dependencies:
    - type: index-pattern
      name: indexpattern-sample
  body: |
    {
      "attributes": {
        "title": "Count of docs 2",
        ...
      },
      "references": [
        {
          "type": "index-pattern",
          "id": "indexpattern-sample",
          "name": "indexpattern-datasource-current-indexpattern"
        },
        {
          "type": "index-pattern",
          "id": "indexpattern-sample",
          "name": "indexpattern-datasource-layer-13e00170-fe50-4495-87ff-048934afece1"
        }
      ]
    }
```

[Complete example](../config/samples/kibana.eck_v1alpha1_lens.yaml)
