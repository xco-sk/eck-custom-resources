# User Role (elasticsearchroles.es.eck.github.com)

Representation of User role. The name `ElasticsearchRole` (instead of plain `Role`)
was chosen to avoid clash with RBAC resources.


## Lifecycle

No special lifecycle is applied for User role - when the role
is deleted from K8s, it is also deleted from ES.

Create and Update are done using the same `PUT /_security/role/` API.
See [Create or update roles API](https://www.elastic.co/guide/en/elasticsearch/reference/current/security-api-put-role.html)
in official documentation.

## Fields

| Key             | Type   | Description                                                               |
|-----------------|--------|---------------------------------------------------------------------------|
| `metadata.name` | string | Name of the Snapshot Lifecycle Policy                                     |
| `spec.targetInstance.name`| string | Name of the [Elasticsearch Instance](cr_elasticsearch_instance.md) to which this ElasticsearchRole will be deployed to |
| `spec.body`     | string | Role definition - same you would use when creating role using ES REST API |

## Example

```yaml
apiVersion: es.eck.github.com/v1alpha1
kind: ElasticsearchRole
metadata:
  name: elasticsearchrole-sample
spec:
  targetInstance:
    name: elasticsearch-quickstart
  body: |
    {
      "cluster": ["all"],
      "indices": [
        {
          "names": [ "index-sample"],
          "privileges": ["all"]
        }
      ],
      "metadata" : {
        "version" : 1
      }
    }
```
