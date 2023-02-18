# Elasticsearch instance (elasticsearchinstances.es.eck.github.com)

Representation of Elasticsearch instance

## Lifecycle

This resource is not reconciled, it is used only to hold the data about the target Elasticsearch instance.

## Fields

| Key                                                     | Type   | Description                                                                                       |
|---------------------------------------------------------|--------|---------------------------------------------------------------------------------------------------|
| `metadata.name`                                         | string | Name of the Elasticsearch instance, used in `targetInstance.name` field, that is present in otherES CRDs to reference the target ES instance |
| `spec.enabled`                                          | bool   | Defines whether this instance is enabled for resource reconciliation                              |
| `spec.url`                                              | string | The URL of Elasticsearch instance                                      |
| `spec.certificate.secretName`                           | string | Name of the secret with CA used for HTTPS communication with Elasticsearch, optional in case of "http://" prefixed URLs |
| `spec.certificate.certificateKey`                       | string | The key with actual certificate data inside the secret defined by `secretName` |
| `spec.authentication.usernamePasswordSecret.secretName` | string | Name of the secret containing user data in username:password form |
| `spec.authentication.usernamePasswordSecret.userName`   | string | The username that will be used for password lookup in secret and also for authentication with target instance |

## Example

```yaml
apiVersion: es.eck.github.com/v1alpha1
kind: ElasticsearchInstance
metadata:
  name: elasticsearch-quickstart
spec:  
  enabled: true
  url: https://quickstart-es-http:9200
  certificate:
    secretName: quickstart-es-http-certs-public
    certificateKey: ca.crt
  authentication:
    usernamePasswordSecret:
      secretName: quickstart-es-elastic-user
      userName: elastic
```
