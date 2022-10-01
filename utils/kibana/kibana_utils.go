package kibana

import (
	"encoding/json"
	"strings"
)

func InjectId(objectJson string, id string) (*string, error) {
	var body map[string]interface{}
	err := json.NewDecoder(strings.NewReader(objectJson)).Decode(&body)
	if err != nil {
		return nil, err
	}

	body["id"] = id

	marshalledBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	sBody := string(marshalledBody)
	return &sBody, nil
}
