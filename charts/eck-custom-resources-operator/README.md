# Helm chart for eck-custom-resources

![Version: 0.5.6](https://img.shields.io/badge/Version-0.5.6-informational?style=flat-square) ![Type: application](https://img.shields.io/badge/Type-application-informational?style=flat-square) ![AppVersion: 0.5.6](https://img.shields.io/badge/AppVersion-0.5.6-informational?style=flat-square)

Helm chart for eck-custom-resources operator

**Homepage:** <https://github.com/xco-sk/eck-custom-resources>

## Maintainers

| Name | Email | Url |
| ---- | ------ | --- |
| Marek Hornak | <marek@xco.sk> | <https://github.com/xco-sk> |

## Source Code

* <https://github.com/xco-sk/eck-custom-resources>

## Installation

```shell
# Add eck-custom-resources helm repo
helm repo add eck-custom-resources https://xco-sk.github.io/eck-custom-resources/

# Install chart
helm install eck-cr eck-custom-resources/eck-custom-resources-operator
```

## Uninstallation
To uninstall the eck-cr named release from Kubernetes cluster, run:

```shell
helm uninstall eck-cr
```

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| affinity | object | `{}` | Affinity settings |
| autoscaling.enabled | bool | `false` | Flag if Horizontal pod autoscaling is used or not |
| autoscaling.maxReplicas | int | `100` | Maximum number of replicas |
| autoscaling.minReplicas | int | `1` | Minimum number of replicas |
| autoscaling.targetCPUUtilizationPercentage | int | `80` | Target CPU utilization percentage metric used for autoscaling decision |
| clusterRole.annotations | object | `{}` | Annotations to add to the service account |
| clusterRole.create | bool | `true` | Specifies whether a service account should be created |
| clusterRole.name | string | `""` | If not set and create is true, a name is generated using the fullname template |
| elasticsearch | object | `{}` | Configuration of Default Elasticsearch cluster to which the Custom resources are deployed. Can stay empty if you want to only use the ElasticsearchInstance CRD approach |
| elasticsearch.authentication.usernamePasswordSecret.secretName | string | `"quickstart-es-elastic-user"` | Name of the Secret containing password for user that is used to manage deployed resources. Should be in the `username: password` format. |
| elasticsearch.authentication.usernamePasswordSecret.userName | string | `"elastic"` | Username of user that is used to manage deployed resources |
| elasticsearch.certificate.certificateKey | string | `"ca.crt"` | Key in Secret that contain the PEM-encoded certificate |
| elasticsearch.certificate.secretName | string | `"quickstart-es-http-certs-public"` | Name of the Secret containing certificate used for communication with Elasticsearch |
| elasticsearch.enabled | bool | `true` | Flag to define if the Elasticsearch reconciler is enabled or not |
| elasticsearch.url | string | `"https://quickstart-es-http:9200"` | Url of Elasticsearch |
| fullnameOverride | string | `""` | Fully qualified app name |
| image.pullPolicy | string | `"IfNotPresent"` | Pull policy for docker image |
| image.repository | string | `"xcosk/eck-custom-resources"` | ECK Custom resources docker image registry |
| image.tag | string | `""` | Docker image tag. Overrides the image tag whose default is the chart appVersion. |
| imagePullSecrets | list | `[]` | Docker image pull secrets |
| kibana | object | `{}` | Configuration of Default Kibana to which the Custom resources are deployed. Can stay empty if you want to only use the KibanaInstance CRD approach |
| kibana.authentication.usernamePasswordSecret.secretName | string | `"quickstart-es-elastic-user"` | Name of the Secret containing password for user that is used to manage deployed resources. Should be in the `username: password` format. |
| kibana.authentication.usernamePasswordSecret.userName | string | `"elastic"` | Username of user that is used to manage deployed resources |
| kibana.certificate.certificateKey | string | `"ca.crt"` | Key in Secret that contain the PEM-encoded certificate |
| kibana.certificate.secretName | string | `"quickstart-kb-http-certs-public"` | Name of the Secret containing certificate used for communication with Kibana |
| kibana.enabled | bool | `true` | Flag to define if the Kibana reconciler is enabled or not |
| kibana.url | string | `"https://quickstart-kb-http:5601"` | Url of Kibana |
| manager.health.healthProbePort | int | `8081` | Port on which the health probe listens |
| manager.leaderElection.leaderElect | bool | `true` | If leader election is enabled |
| manager.webhook.port | int | `9443` | Port on which the webhook listens |
| metrics.enabled | bool | `false` | Flag to indicate if prometheus metrics are exported. If true, the Service and ServiceMonitor resources are deployed alongside the application |
| metrics.service.port | int | `8080` | Metrics service port |
| metrics.service.type | string | `"ClusterIP"` | Metrics service type |
| metrics.serviceMonitor.labels | object | `{}` | Labels to add to the ServiceMonitor |
| metrics.serviceMonitor.namespace | string | `""` | Namespace of the ServiceMonitor |
| nameOverride | string | `""` | Override for Chart.Name default value |
| nodeSelector | object | `{}` | Node selector |
| podAnnotations | object | `{}` | Pod annotation |
| podSecurityContext | object | `{}` | Pod security context |
| replicaCount | int | `1` | Desired number of replicas |
| resources | object | `{}` | Configuration of limits and requests for operator pod |
| securityContext | object | `{}` | Security context |
| serviceAccount.annotations | object | `{}` | Annotations to add to the service account |
| serviceAccount.create | bool | `true` | Specifies whether a service account should be created |
| serviceAccount.name | string | `""` | If not set and create is true, a name is generated using the fullname template |
| tolerations | list | `[]` | Tolerations |

----------------------------------------------
Autogenerated from chart metadata using [helm-docs v1.9.1](https://github.com/norwoodj/helm-docs/releases/v1.9.1)
