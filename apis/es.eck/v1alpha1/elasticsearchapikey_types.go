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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ElasticsearchApikeySpec defines the desired state of ElasticsearchApikey
type ElasticsearchApikeySpec struct {
	// +optional
	TargetConfig CommonElasticsearchConfig `json:"targetInstance,omitempty"`

	Body string `json:"body"`
}

// ElasticsearchApikeyStatus defines the observed state of ElasticsearchApikey
type ElasticsearchApikeyStatus struct {
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// ElasticsearchApikey is the Schema for the elasticsearchApikeys API
type ElasticsearchApikey struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ElasticsearchApikeySpec   `json:"spec,omitempty"`
	Status ElasticsearchApikeyStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ElasticsearchApikeyList contains a list of ElasticsearchApikey
type ElasticsearchApikeyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ElasticsearchApikey `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ElasticsearchApikey{}, &ElasticsearchApikeyList{})
}
