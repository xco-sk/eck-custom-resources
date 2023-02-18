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

// ElasticsearchRoleSpec defines the desired state of ElasticsearchRole
type ElasticsearchRoleSpec struct {
	// +optional
	TargetConfig CommonElasticsearchConfig `json:"targetInstance,omitempty"`

	// +kubebuilder:validation:MinLength=0
	// +required
	Body string `json:"body"`
}

// ElasticsearchRoleStatus defines the observed state of ElasticsearchRole
type ElasticsearchRoleStatus struct {
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// ElasticsearchRole is the Schema for the elasticsearchroles API
type ElasticsearchRole struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ElasticsearchRoleSpec   `json:"spec,omitempty"`
	Status ElasticsearchRoleStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ElasticsearchRoleList contains a list of ElasticsearchRole
type ElasticsearchRoleList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ElasticsearchRole `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ElasticsearchRole{}, &ElasticsearchRoleList{})
}
