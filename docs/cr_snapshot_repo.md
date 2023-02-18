# Snapshot Repository (snapshotrepositories.es.eck.github.com)

CRD that represents Snapshot Repository.

## Lifecycle

No special lifecycle is applied for Snapshot Repositories - when the repo
is deleted from K8s, it is also deleted from ES. As stated in ES documentation,
the data stored in repository are left untouched.

Create and Update are done using the same `PUT /_snapshot/` API.
See [Create or update snapshot repository API](https://www.elastic.co/guide/en/elasticsearch/reference/current/put-snapshot-repo-api.html)
in official documentation.

## Fields

| Key             | Type   | Description                                                                              |
|-----------------|--------|------------------------------------------------------------------------------------------|
| `metadata.name` | string | Name of the Snapshot Repository                                                          |
| `spec.targetInstance.name`| string | Name of the [Elasticsearch Instance](cr_elasticsearch_instance.md) to which this SnapshotRepository will be deployed to |
| `spec.body`     | string | Snapshot repository definition - same you would use when creating repo using ES REST API |

Please keep in mind, the repository location has to be accessible from each and
every cluster node. For `fs` repository type, the `location` needs to be
enabled in `path.repo` field of `elasticsearch.yaml`, see [docs](https://www.elastic.co/guide/en/elasticsearch/reference/current/snapshots-filesystem-repository.html).

## Example

```yaml
apiVersion: es.eck.github.com/v1alpha1
kind: SnapshotRepository
metadata:
  name: snapshotrepository-sample
spec:
  targetInstance:
    name: elasticsearch-quickstart
  body: |
    {
      "type": "fs",
      "settings": {
        "location": "/tmp/"
      }
    }
```
