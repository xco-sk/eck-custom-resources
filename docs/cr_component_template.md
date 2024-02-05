# Component Template (componenttemplates.es.eck.github.com)

Representation of the Component template resource.

## Lifecycle

Component template lifecycle is simple - when the template is deleted
from K8s, it is also deleted from ES.
See [Create or update component template API](https://www.elastic.co/guide/en/elasticsearch/reference/current/indices-component-template.html)
in official documentation.

## Fields

| Key                                    | Type   | Description                                                                                                            |
| -------------------------------------- | ------ | ---------------------------------------------------------------------------------------------------------------------- |
| `metadata.name`                        | string | Name of the Component Template                                                                                         |
| `spec.targetInstance.name`             | string | Name of the [Elasticsearch Instance](cr_elasticsearch_instance.md) to which this ComponentTemplate will be deployed to |
| `spec.body`                            | string | Component template definition - same you would use when creating component template using ES REST API                  |
| `spec.dependencies.indexTemplates`     | list   | List of index templates that have to be present in ES cluster before component template is created / updated           |
| `spec.dependencies.indices`            | list   | List of indices that have to be present in ES cluster before component template is created / updated                   |
| `spec.dependencies.conponentTemplates` | list   | List of component templates that have to be present in ES cluster before component template is created / updated       |

## Example

```yaml
apiVersion: es.eck.github.com/v1alpha1
kind: ComponentTemplate
metadata:
  labels:
    app.kubernetes.io/name: componenttemplate
    app.kubernetes.io/instance: componenttemplate-sample
    app.kubernetes.io/part-of: eck-custom-resources
    app.kubernetes.io/managed-by: kustomize
    app.kubernetes.io/created-by: eck-custom-resources
  name: componenttemplate-sample
spec:
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
