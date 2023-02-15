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

// SnapshotRepositorySpec defines the desired state of SnapshotRepository
type SnapshotRepositorySpec struct {
	// +optional
	TargetConfig CommonElasticsearchConfig `json:"targetInstance,omitempty"`

	Body string `json:"body"`
}

// SnapshotRepositoryStatus defines the observed state of SnapshotRepository
type SnapshotRepositoryStatus struct {
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// SnapshotRepository is the Schema for the snapshotrepositories API
type SnapshotRepository struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SnapshotRepositorySpec   `json:"spec,omitempty"`
	Status SnapshotRepositoryStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// SnapshotRepositoryList contains a list of SnapshotRepository
type SnapshotRepositoryList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SnapshotRepository `json:"items"`
}

func init() {
	SchemeBuilder.Register(&SnapshotRepository{}, &SnapshotRepositoryList{})
}
