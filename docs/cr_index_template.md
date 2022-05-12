# Index Template (indextemplates.es.eck.github.com)

Representation of Index template resource.

## Lifecycle

Index template lifecycle is simple - when the template is deleted
from K8s, it is also deleted from ES.
Create and Update are done using the same `PUT /_index_template` API. 
See [Create or update index template API](https://www.elastic.co/guide/en/elasticsearch/reference/current/indices-put-template.html)
in official documentation.

## Fields

| Key             | Type   | Description                                                                                   |
|-----------------|--------|-----------------------------------------------------------------------------------------------|
| `metadata.name` | string | Name of the Index Template                                                                    |
| `spec.body`     | string | Index template definition - same you would use when creating index template using ES REST API |
| `spec.dependsOn.indexTemplates` | list | List of index templates that have to be present in ES cluster before index is created / updated |
| `spec.dependsOn.indices`        | list | List of indices that have to be present in ES cluster before index created / updated            |

## Example

```yaml
apiVersion: es.eck.github.com/v1alpha1
kind: IndexTemplate
metadata:
  name: indextemplate-sample
spec:
  dependsOn:
    indexTemplates:
      - indexTemplateName: indextemplate-base
    indices:
      - indexName: index-base-sample
  body: |
    {
      "index_patterns" : ["index-*"],
      "priority" : 1,
      "template": {
        "settings" : {
          "number_of_shards" : 2,
          "number_of_replicas" : 0
        }
      }
    }
```