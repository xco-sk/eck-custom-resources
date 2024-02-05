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

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ComponentTemplateSpec defines the desired state of ComponentTemplate
type ComponentTemplateSpec struct {
	// +optional
	TargetConfig CommonElasticsearchConfig `json:"targetInstance,omitempty"`
	// +optional
	Dependencies Dependencies `json:"dependencies,omitempty"`

	Body string `json:"body"`
}

// ComponentTemplateStatus defines the observed state of ComponentTemplate
type ComponentTemplateStatus struct {
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// ComponentTemplate is the Schema for the componenttemplates API
type ComponentTemplate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ComponentTemplateSpec   `json:"spec,omitempty"`
	Status ComponentTemplateStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ComponentTemplateList contains a list of ComponentTemplate
type ComponentTemplateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ComponentTemplate `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ComponentTemplate{}, &ComponentTemplateList{})
}
