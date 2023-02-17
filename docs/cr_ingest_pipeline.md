# Ingest Pipeline (ingestpipelines.es.eck.github.com)

Representation of Ingest Pipeline.

## Lifecycle

No special lifecycle is applied for Ingest Pipeline - when the pipeline
is deleted from K8s, it is also deleted from ES.
Create and Update are done using the same `PUT /_ingest/pipeline/` API.
See [Create or update pipeline API](https://www.elastic.co/guide/en/elasticsearch/reference/current/put-pipeline-api.html)
in official documentation.

## Fields

| Key                       | Type   | Description                                                                                     |
|---------------------------|--------|-------------------------------------------------------------------------------------------------|
| `metadata.name`           | string | Name of the Ingest Pipeline                                                                     |
| `spec.targetInstance.name`| string | Name of the [Elasticsearch Instance](cr_elasticsearch_instance.md) to which this IngestPipeline will be deployed to |
| `spec.body`               | string | Ingest Pipeline definition - same you would use when creating ingest pipeline using ES REST API |

## Example

```yaml
apiVersion: es.eck.github.com/v1alpha1
kind: IngestPipeline
metadata:
  name: ingestpipeline-sample
spec:
  targetInstance:
    name: elasticsearch-quickstart
  body: |
    {
      "description" : "Ingest pipeline sample",
      "processors" : [
        {
          "set" : {
            "description" : "Ingest pipeline sample processor",
            "field": "eck-custom-resources",
            "value": "true"
          }
        }
      ]
    }
```
