# Space (spaces.kibana.eck.github.com)

Custom resource definition representing Space in Kibana.

## Lifecycle

Space lifecycle is simple - when it is deleted from K8s, it is also deleted from Kibana. Keep in mind, when you delete space,
all saved objects belonging to that space will be deleted as well. That might lead to inconsistent state when
resources deployed to K8s are still present in K8s, but they are deleted from Kibana - the synchronization of resources
goes only one way at the moment.

The value of `id` field in `body` is always added/replaced with value from `metadata.name`

See [Spaces APIs](https://www.elastic.co/guide/en/kibana/master/spaces-api.html) in official documentation.

## Fields

| Key             | Type   | Description                                                                                     | Default    |
|-----------------|--------|-------------------------------------------------------------------------------------------------|------------|
| `metadata.name` | string | Name of the Visualization, used also as its ID in Kibana                                        | No default |
| `spec.targetInstance.name`| string         | Name of the [Kibana Instance](cr_kibana_instance.md) to which this Space will be deployed to | The operator configuration |
| `spec.body`     | string | Space definition json, `id` field value is added/replaced with value from `metadata.name` field | No default |

## Example

```yaml
apiVersion: kibana.eck.github.com/v1alpha1
kind: Space
metadata:
  name: space-sample
spec:
  targetInstance:
    name: kibana-quickstart
  body: |
    {
      "name": "ECK Space sample",
      "description" : "This space was created using ECK-CR operator",
      "color": "#aabbcc",
      "initials": "CR",
      "disabledFeatures": [],
      "imageUrl": ""
    }
```
