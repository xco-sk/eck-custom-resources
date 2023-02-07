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

// SnapshotLifecyclePolicySpec defines the desired state of SnapshotLifecyclePolicy
type SnapshotLifecyclePolicySpec struct {
	// +optional
	CommonConfig CommonElasticsearchConfig `json:",inline"`

	// +kubebuilder:validation:MinLength=0
	// +required
	Body string `json:"body"`
}

// SnapshotLifecyclePolicyStatus defines the observed state of SnapshotLifecyclePolicy
type SnapshotLifecyclePolicyStatus struct {
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// SnapshotLifecyclePolicy is the Schema for the snapshotlifecyclepolicies API
type SnapshotLifecyclePolicy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SnapshotLifecyclePolicySpec   `json:"spec,omitempty"`
	Status SnapshotLifecyclePolicyStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// SnapshotLifecyclePolicyList contains a list of SnapshotLifecyclePolicy
type SnapshotLifecyclePolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SnapshotLifecyclePolicy `json:"items"`
}

func init() {
	SchemeBuilder.Register(&SnapshotLifecyclePolicy{}, &SnapshotLifecyclePolicyList{})
}
