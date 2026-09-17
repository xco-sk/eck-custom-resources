package elasticsearch

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/elastic/go-elasticsearch/v9"
)

type stubTransport struct {
	status int
	body   string
}

func (s stubTransport) RoundTrip(*http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: s.status,
		Header:     http.Header{"Content-Type": []string{"application/json"}, "X-Elastic-Product": []string{"Elasticsearch"}},
		Body:       io.NopCloser(strings.NewReader(s.body)),
	}, nil
}

func clientWith(status int, body string) *elasticsearch.Client {
	c, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{"http://localhost:9200"},
		Transport: stubTransport{status: status, body: body},
	})
	if err != nil {
		panic(err)
	}
	return c
}

func TestVerifyIndexEmpty(t *testing.T) {
	cases := []struct {
		name    string
		status  int
		body    string
		want    bool
		wantErr bool
	}{
		{"empty index", 200, `{"count":0}`, true, false},
		{"non-empty index", 200, `{"count":5}`, false, false},
		{"index_not_found 404 (issue #79 panic)", 404, `{"error":{"type":"index_not_found_exception"},"status":404}`, false, true},
		{"auth failure 403", 403, `{"error":{"type":"security_exception"},"status":403}`, false, true},
		{"200 but no count field", 200, `{}`, false, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Fatalf("PANIC (regression of issue #79): %v", r)
				}
			}()
			got, err := VerifyIndexEmpty(clientWith(tc.status, tc.body), "some-index")
			if tc.wantErr && err == nil {
				t.Errorf("expected an error, got nil (empty=%v)", got)
			}
			if !tc.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if !tc.wantErr && got != tc.want {
				t.Errorf("empty = %v, want %v", got, tc.want)
			}
		})
	}
}
