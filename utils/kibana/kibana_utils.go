package kibana

import (
	"context"
	"encoding/json"
	"strings"

	kibanaeckv1alpha1 "github.com/xco-sk/eck-custom-resources/apis/kibana.eck/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/client"
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

func GetTargetInstance(cli client.Client, ctx context.Context, namespace string, targetName string, kibanaInstance *kibanaeckv1alpha1.KibanaInstance) error {
	if err := cli.Get(ctx, client.ObjectKey{Namespace: namespace, Name: targetName}, kibanaInstance); err != nil {
		return err
	}
	return nil
}
