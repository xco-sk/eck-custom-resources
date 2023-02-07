package v1alpha1

type CommonElasticsearchConfig struct {
	// +optional
	ElasticsearchInstance string `json:"targetInstance,omitempty"`
}
