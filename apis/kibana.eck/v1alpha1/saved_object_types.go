package v1alpha1

type SavedObject struct {
	Space *string `json:"space,omitempty"`
	Body  string  `json:"body"`
}

func (t *SavedObject) GetSavedObject() SavedObject {
	return SavedObject{
		Space: t.Space,
		Body:  t.Body,
	}
}
