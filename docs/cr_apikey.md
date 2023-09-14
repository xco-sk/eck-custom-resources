# ElasticsearchApikey (elasticsearchapikeys.es.eck.github.com)

Representation of API key. The name `ElasticsearchApikey` (instead of plain `APIkey`)
was chosen to avoid clash with RBAC resources.

## Lifecycle

API key resource lifecycle is simple - when the apikey is created, it will create k8s secret object with name defined in `metadata.name` that contained an encoded API key value from ES, when the apikey is deleted from K8s, 
it is also deleted from ES including the k8s secret object created during creation process.
Creation of new resource is reconciled using `POST /_security/apikey/` API. Deletion is done using `DELETE /_security/api_key` to invalidate the API key.

See [Create API keys API](https://www.elastic.co/guide/en/elasticsearch/reference/current/security-api-create-api-key.html) [Delete API keys API](https://www.elastic.co/guide/en/elasticsearch/reference/current/security-api-invalidate-api-key.html)
in official documentation.

## Fields

| Key               | Type   | Description                                                                                                                                   |
|-------------------|--------|-----------------------------------------------------------------------------------------------------------------------------------------------|
| `metadata.name`   | string | Name of the Index Lifecycle Policy                                                                                                            |
| `spec.targetInstance.name`| string | Name of the [Elasticsearch Instance](cr_elasticsearch_instance.md) to which this ElasticsearchApikey will be deployed to |
| `spec.body`       | string | API key definition - same you would use when creating API key using ES REST API                                                                     |


## Example

```yaml
apiVersion: es.eck.github.com/v1alpha1
kind: ElasticsearchApikey
metadata:
  name: elasticsearchapikey-sample
spec:
  targetInstance:
    name: elasticsearch-quickstart
  body: |
    {
      "name" : "elasticsearchapikey-sample"
    }
```
