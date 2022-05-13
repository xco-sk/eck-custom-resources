# Index Pattern (indexpatterns.kibana.eck.github.com)

Custom resource definition representing Index Pattern in Kibana.

## Lifecycle

Index pattern lifecycle is simple - when the pattern is deleted from K8s, it is also deleted from Kibana.
Creation of new resource is reconciled using `POST /api/saved_objects/index-pattern/` API. Update is done using
`PUT /api/saved_objects/index-pattern/`.

See [Index patterns APIs](https://www.elastic.co/guide/en/kibana/8.2/index-patterns-api.html) in official documentation.

## Fields

| Key             | Type   | Description                                              |
|-----------------|--------|----------------------------------------------------------|
| `metadata.name` | string | Name of the Index Pattern, used also as its ID in Kibana |
| `spec.body`     | string | Index pattern definition json                            |

## Example

```yaml
apiVersion: kibana.eck.github.com/v1alpha1
kind: IndexPattern
metadata:
  name: indexpattern-sample
spec:
  body: |
    {
        "attributes": {
            "title": "index-*",
            "timeFieldName": "@timestamp"
        }
    }
```
