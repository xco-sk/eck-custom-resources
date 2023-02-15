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

// DataViewSpec defines the desired state of DataView
type DataViewSpec struct {
	// +optional
	TargetConfig CommonKibanaConfig `json:"targetInstance,omitempty"`

	SavedObject `json:",inline"`
}

// DataViewStatus defines the observed state of DataView
type DataViewStatus struct {
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// DataView is the Schema for the dataviews API
type DataView struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DataViewSpec   `json:"spec,omitempty"`
	Status DataViewStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// DataViewList contains a list of DataView
type DataViewList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DataView `json:"items"`
}

func init() {
	SchemeBuilder.Register(&DataView{}, &DataViewList{})
}
