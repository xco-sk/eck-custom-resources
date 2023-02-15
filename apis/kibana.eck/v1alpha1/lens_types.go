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

// LensSpec defines the desired state of Lens
type LensSpec struct {
	// +optional
	TargetConfig CommonKibanaConfig `json:"targetInstance,omitempty"`

	SavedObject `json:",inline"`
}

// LensStatus defines the observed state of Lens
type LensStatus struct {
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Lens is the Schema for the lens API
type Lens struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   LensSpec   `json:"spec,omitempty"`
	Status LensStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// LensList contains a list of Lens
type LensList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Lens `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Lens{}, &LensList{})
}
