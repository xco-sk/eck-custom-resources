# Index Template (indextemplates.es.eck.github.com)

Representation of the Index template resource.

## Lifecycle

Index template lifecycle is simple - when the template is deleted
from K8s, it is also deleted from ES.
Create and Update are done using the same `PUT /_index_template` API. 
See [Create or update index template API](https://www.elastic.co/guide/en/elasticsearch/reference/current/indices-put-template.html)
in official documentation.

## Fields

| Key                                | Type   | Description                                                                                     |
|------------------------------------|--------|-------------------------------------------------------------------------------------------------|
| `metadata.name`                    | string | Name of the Index Template                                                                      |
| `spec.targetInstance.name`         | string | Name of the [Elasticsearch Instance](cr_elasticsearch_instance.md) to which this IndexTemplate will be deployed to |
| `spec.body`                        | string | Index template definition - same you would use when creating index template using ES REST API   |
| `spec.dependencies.indexTemplates` | list   | List of index templates that have to be present in ES cluster before index is created / updated |
| `spec.dependencies.indices`        | list   | List of indices that have to be present in ES cluster before index created / updated            |

## Example

```yaml
apiVersion: es.eck.github.com/v1alpha1
kind: IndexTemplate
metadata:
  name: indextemplate-sample
spec:
  targetInstance:
    name: elasticsearch-quickstart
  dependencies:
    indexTemplates:
      - indextemplate-base
    indices:
      - index-base-sample
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
