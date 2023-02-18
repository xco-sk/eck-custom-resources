# Dashboard (dashboards.kibana.eck.github.com)

Custom resource definition representing Dashboard in Kibana.

## Lifecycle

Dashboard lifecycle is simple - when it is deleted from K8s, it is also deleted from Kibana. Creation of
new resource is reconciled using `POST /api/saved_objects/dashboard/` API. Update is done using
`PUT /api/saved_objects/dashboard/`. In case the `spec.space` is filled in, the URLs are prefixed
with `/s/<spec.space>`.

See [Saved objects APIs](https://www.elastic.co/guide/en/kibana/master/saved-objects-api.html) in official documentation.

## Fields

| Key                         | Type            | Description                                                                                                                                     | Default                                              |
|-----------------------------|-----------------|-------------------------------------------------------------------------------------------------------------------------------------------------|------------------------------------------------------|
| `metadata.name`             | string          | Name of the Dashboard, used also as its ID in Kibana                                                                                            | No default                                           |
| `spec.space`                | string          | Name of the Kibana namespace to which the Dashboard is deployed to                                                                              | No default (will be deployed to "default" namespace) |
| `spec.targetInstance.name`  | string         | Name of the [Kibana Instance](cr_kibana_instance.md) to which this Dashboard will be deployed to | The operator configuration |
| `spec.body`                 | string          | Dashboard definition json (omitting everything except attributes and references)                                                                | No default                                           |
| `spec.dependencies`         | List of objects | List of dependencies - the reconciler will wait for all resources from the list to be present in Kibana before deploying/updating this resource | -                                                    |                                                 |
| `spec.dependencies[].space` | string          | Kibana Space where to look for given resource                                                                                                   | -                                                    |
| `spec.dependencies[].type`  | string          | Type of resource - one of `visualization, dashboard, search, index-pattern, lens`                                                               | -                                                    |
| `spec.dependencies[].name`  | string          | Name of resource                                                                                                                                | -                                                    |

## Example

```yaml
apiVersion: kibana.eck.github.com/v1alpha1
kind: Dashboard
metadata:
  name: dashboard-sample
spec:
  targetInstance:
    name: kibana-quickstart
  space: my-space
  dependencies:
    - type: lens
      name: lens-sample
  body: |
    {
      "attributes": {
        "title": "Sample dashboard",
        "hits": 0,
        "description": "",
        "panelsJSON": "[{\"version\":\"8.1.0\",\"type\":\"lens\",\"gridData\":{\"x\":0,\"y\":0,\"w\":24,\"h\":15,\"i\":\"4cfcdf6b-1729-4467-a911-c69be15d58f8\"},\"panelIndex\":\"4cfcdf6b-1729-4467-a911-c69be15d58f8\",\"embeddableConfig\":{\"enhancements\":{}},\"panelRefName\":\"panel_4cfcdf6b-1729-4467-a911-c69be15d58f8\"}]",
        "optionsJSON": "{\"useMargins\":true,\"syncColors\":false,\"hidePanelTitles\":false}",
        "version": 1,
        "timeRestore": false,
        "kibanaSavedObjectMeta": {
          "searchSourceJSON": "{\"query\":{\"query\":\"\",\"language\":\"kuery\"},\"filter\":[]}"
        }
      },
      "references": [
        {
          "name": "4cfcdf6b-1729-4467-a911-c69be15d58f8:panel_4cfcdf6b-1729-4467-a911-c69be15d58f8",
          "type": "lens",
          "id": "lens-sample"
        }
      ]
    }

```
