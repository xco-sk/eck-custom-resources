/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v2

// KibanaSpec Definition of target elasticsearch cluster
type KibanaSpec struct {
	// +required
	Enabled bool `json:"enabled"`

	// +required
	// +kubebuilder:validation:MinLength=0
	Url string `json:"url,omitempty"`
	// +optional
	Certificate *PublicCertificate `json:"certificate,omitempty"`

	// +optional
	Authentication *KibanaAuthentication `json:"authentication,omitempty"`
}

// KibanaAuthentication Definition of Kibana authentication
type KibanaAuthentication struct {
	// +optional
	UsernamePassword *UsernamePasswordAuthentication `json:"usernamePasswordSecret,omitempty"`
}
