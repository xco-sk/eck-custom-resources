package v2

// UsernamePasswordAuthentication Definition of Username/Password authentication
type UsernamePasswordAuthentication struct {
	// +kubebuilder:validation:MinLength=0
	SecretName string `json:"secretName"`

	// +kubebuilder:validation:MinLength=0
	UserName string `json:"userName"`
}

// APIKey Definition of APIKey authentication
type APIKeyAuthentication struct {
	// +kubebuilder:validation:MinLength=0
	APIKey string `json:"apiKey"`
}

// PublicCertificate Configuration for public certificate used for communication with target
type PublicCertificate struct {
	// +required
	// +kubebuilder:validation:MinLength=0
	SecretName string `json:"secretName"`

	// +reqired
	// +kubebuilder:validation:MinLength=0
	CertificateKey string `json:"certificateKey"`
}
