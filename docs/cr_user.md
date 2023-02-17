# ElasticsearchUser (elasticsearchusers.es.eck.github.com)

Representation of User. The name `ElasticsearchUser` (instead of plain `User`)
was chosen to avoid clash with RBAC resources.

## Lifecycle

User resource lifecycle is simple - when the user is deleted from K8s, 
it is also deleted from ES.
Create and Update are done using the same `PUT /_security/user/` API.
See [Create or update users API](https://www.elastic.co/guide/en/elasticsearch/reference/current/security-api-put-user.html)
in official documentation.

## Fields

| Key               | Type   | Description                                                                                                                                   |
|-------------------|--------|-----------------------------------------------------------------------------------------------------------------------------------------------|
| `metadata.name`   | string | Name of the Index Lifecycle Policy                                                                                                            |
| `spec.targetInstance.name`| string | Name of the [Elasticsearch Instance](cr_elasticsearch_instance.md) to which this ElasticsearchUser will be deployed to |
| `spec.secretName` | string | The name of the secret, from where the password is taken during create or update, the key has to be equal to username (`metadata.name` field) |
| `spec.body`       | string | User definition - same you would use when creating User using ES REST API                                                                     |

The changes in secret (e.g. password rotation) **are not** automatically propagated.

## Example

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: elasticsearchuser-secret
type: Opaque
data:
  elasticsearchuser-sample: sample.password

---
apiVersion: es.eck.github.com/v1alpha1
kind: ElasticsearchUser
metadata:
  name: elasticsearchuser-sample
spec:
  targetInstance:
    name: elasticsearch-quickstart
  secretName: elasticsearchuser-secret
  body: |
    {
      "roles" : [ "admin", "elasticsearchrole-sample" ],
      "full_name" : "Richard Feynman",
      "email" : "rfeynman@example.com",
      "metadata" : {
        "intelligence" : 7
      }
    }
```
