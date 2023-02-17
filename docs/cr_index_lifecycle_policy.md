# Index Lifecycle Policy (indexlifecyclepolicies.es.eck.github.com)

Representation of Index Lifecycle policy.

## Lifecycle

Index Lifecycle policy resource lifecycle is simple - when the policy 
is deleted from K8s, it is also deleted from ES.
Create and Update are done using the same `PUT _ilm/policy` API.
See [Create or update lifecycle policy API](https://www.elastic.co/guide/en/elasticsearch/reference/current/ilm-put-lifecycle.html)
in official documentation.

## Fields

| Key                       | Type   | Description                                                                                       |
|---------------------------|--------|---------------------------------------------------------------------------------------------------|
| `metadata.name`           | string | Name of the Index Lifecycle Policy                                                                |
| `spec.targetInstance.name`| string | Name of the [Elasticsearch Instance](cr_elasticsearch_instance.md) to which this IndexLifecyclePolicy will be deployed to |
| `spec.body`               | string | Index Lifecycle Policy definition - same you would use when creating ILM policy using ES REST API |

## Example

```yaml
apiVersion: es.eck.github.com/v1alpha1
kind: IndexLifecyclePolicy
metadata:
  name: indexlifecyclepolicy-sample
spec:
  targetInstance:
    name: elasticsearch-quickstart
  body: |
    {
      "policy": {
        "_meta": {
          "description": "used as example for lifecycle policy",
          "project": {
            "name": "eck-custom-resources"
          }
        },
        "phases": {
          "warm": {
            "min_age": "10d",
            "actions": {
              "forcemerge": {
                "max_num_segments": 1
              }
            }
          },
          "delete": {
            "min_age": "30d",
            "actions": {
              "delete": {}
            }
          }
        }
      }
    }
```
