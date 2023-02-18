# Kibana instance (kibanainstances.kibana.eck.github.com)

Representation of Kibana instance

## Lifecycle

This resource is not reconciled, it is used only to hold the data about the target Kibana instance.

## Fields

| Key                                                     | Type   | Description                                                                                       |
|---------------------------------------------------------|--------|---------------------------------------------------------------------------------------------------|
| `metadata.name`                                         | string | Name of the Kibana instance, used in `targetInstance.name` field, that is present in other Kibana CRDs to reference the target Kibana instance |
| `spec.enabled`                                          | bool   | Defines whether this instance is enabled for resource reconciliation                              |
| `spec.url`                                              | string | The URL of Kibana instance                                      |
| `spec.certificate.secretName`                           | string | Name of the secret with CA used for HTTPS communication with Kibana, optional in case of "http://" prefixed URLs |
| `spec.certificate.certificateKey`                       | string | The key with actual certificate data inside the secret defined by `secretName` |
| `spec.authentication.usernamePasswordSecret.secretName` | string | Name of the secret containing user data in username:password form |
| `spec.authentication.usernamePasswordSecret.userName`   | string | The username that will be used for password lookup in secret and also for authentication with target instance |

## Example

```yaml
apiVersion: kibana.eck.github.com/v1alpha1
kind: KibanaInstance
metadata:
  name: kibana-quickstart
spec:
  enabled: true
  url: https://quickstart-kb-http:5601
  certificate:
    secretName: quickstart-kb-http-certs-public
    certificateKey: ca.crt
  authentication:
    usernamePasswordSecret:
      secretName: quickstart-es-elastic-user
      userName: elastic
```
