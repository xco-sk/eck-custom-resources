//go:build integration

package elasticsearch

import (
	"os"
	"testing"

	"github.com/elastic/go-elasticsearch/v9"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/xco-sk/eck-custom-resources/apis/es.eck/v1alpha1"
)

// Exercises the operator's real code paths against whatever Elasticsearch is at
// ES_TEST_URL, to establish whether the v9 client works against an 8.x server.
func testClient(t *testing.T) *elasticsearch.Client {
	url := os.Getenv("ES_TEST_URL")
	if url == "" {
		t.Skip("ES_TEST_URL not set")
	}
	c, err := elasticsearch.NewClient(elasticsearch.Config{Addresses: []string{url}})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	return c
}

func TestServerVersion(t *testing.T) {
	c := testClient(t)
	res, err := c.Info()
	if err != nil {
		t.Fatalf("Info() FAILED - client refuses this server: %v", err)
	}
	defer res.Body.Close()
	if res.IsError() {
		t.Fatalf("Info() error response: %s", res.String())
	}
	t.Logf("Info() OK: %s", res.String()[:120])
}

func TestIndexLifecyclePolicyRoundTrip(t *testing.T) {
	c := testClient(t)
	// The real betsys-30-days policy from production, freeze action included.
	p := v1alpha1.IndexLifecyclePolicy{
		ObjectMeta: metav1.ObjectMeta{Name: "compat-betsys-30-days"},
		Spec: v1alpha1.IndexLifecyclePolicySpec{Body: `{
  "policy": {
    "phases": {
      "hot": {"actions": {"rollover": {"max_primary_shard_size": "20gb","max_age": "3d","min_docs": "0"},"set_priority": {"priority": 50}}},
      "warm": {"min_age": "3d","actions": {"forcemerge": {"max_num_segments": 1},"shrink": {"number_of_shards": 1,"allow_write_after_shrink": false},"set_priority": {"priority": 25}}},
      "cold": {"min_age": "14d","actions": {"set_priority": {"priority": 0},"freeze": {},"allocate": {"number_of_replicas": 0}}},
      "delete": {"min_age": "30d","actions": {"delete": {"delete_searchable_snapshot": true}}}
    }
  }
}`},
	}
	if _, err := UpsertIndexLifecyclePolicy(c, p); err != nil {
		t.Fatalf("UpsertIndexLifecyclePolicy FAILED: %v", err)
	}
	t.Log("UpsertIndexLifecyclePolicy OK (incl. cold-phase freeze action)")
	if _, err := DeleteIndexLifecyclePolicy(c, "compat-betsys-30-days"); err != nil {
		t.Fatalf("DeleteIndexLifecyclePolicy FAILED: %v", err)
	}
	t.Log("DeleteIndexLifecyclePolicy OK")
}

func TestIndexTemplateRoundTrip(t *testing.T) {
	c := testClient(t)
	// The real k8s-template from production, logsdb mode included.
	tpl := v1alpha1.IndexTemplate{
		ObjectMeta: metav1.ObjectMeta{Name: "compat-k8s-template"},
		Spec: v1alpha1.IndexTemplateSpec{Body: `{
  "index_patterns": ["compat-k8s-*"],
  "data_stream": {},
  "priority": 500,
  "template": {
    "settings": {"number_of_shards": 2,"number_of_replicas": 1,
      "index": {"mode": "logsdb","sort": {"field": "@timestamp","order": "desc"}}},
    "mappings": {"properties": {"@timestamp": {"type": "date"}}}
  }
}`},
	}
	if _, err := UpsertIndexTemplate(c, tpl); err != nil {
		t.Fatalf("UpsertIndexTemplate FAILED: %v", err)
	}
	t.Log("UpsertIndexTemplate OK (incl. logsdb index mode)")
	if _, err := DeleteIndexTemplate(c, "compat-k8s-template"); err != nil {
		t.Fatalf("DeleteIndexTemplate FAILED: %v", err)
	}
	t.Log("DeleteIndexTemplate OK")
}

func TestComponentTemplateRoundTrip(t *testing.T) {
	c := testClient(t)
	ct := &v1alpha1.ComponentTemplate{
		ObjectMeta: metav1.ObjectMeta{Name: "compat-metrics-custom"},
		Spec: v1alpha1.ComponentTemplateSpec{
			Body: `{"template":{"settings":{"index.number_of_replicas":1}}}`,
		},
	}
	if _, err := UpsertComponentTemplate(c, ct); err != nil {
		t.Fatalf("UpsertComponentTemplate FAILED: %v", err)
	}
	t.Log("UpsertComponentTemplate OK")
	if _, err := DeleteComponentTemplate(c, "compat-metrics-custom"); err != nil {
		t.Fatalf("DeleteComponentTemplate FAILED: %v", err)
	}
	t.Log("DeleteComponentTemplate OK")
}

func TestIndexExistsAndEmpty(t *testing.T) {
	c := testClient(t)
	exists, err := VerifyIndexExists(c, "definitely-not-here")
	if err != nil {
		t.Fatalf("VerifyIndexExists FAILED: %v", err)
	}
	if exists {
		t.Error("phantom index reported as existing")
	}
	t.Log("VerifyIndexExists OK (404 handled)")

	// Must return an error, not panic (issue #79).
	if _, err := VerifyIndexEmpty(c, "definitely-not-here"); err == nil {
		t.Error("VerifyIndexEmpty on a missing index should error")
	} else {
		t.Logf("VerifyIndexEmpty OK (errors cleanly): %v", err)
	}
}
