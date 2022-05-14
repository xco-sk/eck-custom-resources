package v1alpha1

type SavedObject struct {
	Space        *string      `json:"space,omitempty"`
	Body         string       `json:"body"`
	Dependencies []Dependency `json:"dependencies,omitempty"`
}

type Dependency struct {
	ObjectType SavedObjectType `json:"type"`
	Name       string          `json:"name"`
	Space      *string         `json:"space,omitempty"`
}

// +kubebuilder:validation:Enum=visualization;dashboard;search;index-pattern;lens
type SavedObjectType string

func (in *SavedObject) GetSavedObject() SavedObject {
	return SavedObject{
		Space:        in.Space,
		Body:         in.Body,
		Dependencies: in.Dependencies,
	}
}
