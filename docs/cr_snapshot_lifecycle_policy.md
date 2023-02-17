# Snapshot Lifecycle Policy (snapshotlifecyclepolicies.es.eck.github.com)

CRD that represents Snapshot Lifecycle Policy.

## Lifecycle

No special lifecycle is applied for Snapshot Lifecycle Policies - when the policy
is deleted from K8s, it is also deleted from ES.

Create and Update are done using the same `PUT /_slm/policy/` API.
See [Create or update snapshot lifecycle policy API](https://www.elastic.co/guide/en/elasticsearch/reference/current/slm-api-put-policy.html)
in official documentation.

## Fields

| Key                       | Type   | Description                                                                                      |
|---------------------------|--------|--------------------------------------------------------------------------------------------------|
| `metadata.name`           | string | Name of the Snapshot Lifecycle Policy                                                            |
| `spec.targetInstance.name`| string | Name of the [Elasticsearch Instance](cr_elasticsearch_instance.md) to which this SnapshotLifecyclePolicy will be deployed to |
| `spec.body`               | string | Snapshot Lifecycle Policy definition - same you would use when creating policy using ES REST API |

## Example

```yaml
apiVersion: es.eck.github.com/v1alpha1
kind: SnapshotLifecyclePolicy
metadata:
  name: snapshotlifecyclepolicy-sample
spec:
  targetInstance:
    name: elasticsearch-quickstart
  body: |
    {
      "schedule": "0 30 1 * * ?", 
      "name": "<daily-snap-{now/d}>", 
      "repository": "snapshotrepository-sample", 
      "config": { 
        "indices": ["*"], 
        "ignore_unavailable": false,
        "include_global_state": true
      },
      "retention": { 
        "expire_after": "30d", 
        "min_count": 5, 
        "max_count": 50 
      }
    }
```
