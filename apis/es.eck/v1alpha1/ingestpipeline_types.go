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

// IngestPipelineSpec defines the desired state of IngestPipeline
type IngestPipelineSpec struct {
	// +optional
	TargetConfig CommonElasticsearchConfig `json:"targetInstance,omitempty"`

	Body string `json:"body"`
}

// IngestPipelineStatus defines the observed state of IngestPipeline
type IngestPipelineStatus struct {
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// IngestPipeline is the Schema for the ingestpipelines API
type IngestPipeline struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   IngestPipelineSpec   `json:"spec,omitempty"`
	Status IngestPipelineStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// IngestPipelineList contains a list of IngestPipeline
type IngestPipelineList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []IngestPipeline `json:"items"`
}

func init() {
	SchemeBuilder.Register(&IngestPipeline{}, &IngestPipelineList{})
}
