package v1alpha1

type CommonKibanaConfig struct {
	// +optional
	KibanaInstance *string `json:"targetInstance,omitempty"`
}
