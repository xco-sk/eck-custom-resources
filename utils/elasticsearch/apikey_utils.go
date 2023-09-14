package elasticsearch

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/xco-sk/eck-custom-resources/apis/es.eck/v1alpha1"
	"github.com/xco-sk/eck-custom-resources/utils"
	k8sv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func DeleteApikey(cli client.Client, ctx context.Context, esClient *elasticsearch.Client, apikey v1alpha1.ElasticsearchApikey, req ctrl.Request) (ctrl.Result, error) {
	res, err := esClient.Security.InvalidateAPIKey(strings.NewReader(apikey.Spec.Body))
	if err != nil || res.IsError() {
		return utils.GetRequeueResult(), err
	}

	if err := DeleteApikeySecret(cli, ctx, req.Namespace, req.Name); err != nil {
		return utils.GetRequeueResult(), err
	}

	return ctrl.Result{}, nil
}

func CreateApikey(cli client.Client, ctx context.Context, esClient *elasticsearch.Client, apikey v1alpha1.ElasticsearchApikey, req ctrl.Request) (ctrl.Result, error) {
	response, err := esClient.Security.CreateAPIKey(strings.NewReader(apikey.Spec.Body))

	if err != nil || response.IsError() {
		return utils.GetRequeueResult(), GetClientErrorOrResponseError(err, response)
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return utils.GetRequeueResult(), err
	}

	var responseMap map[string]interface{}
	errs := json.Unmarshal([]byte(body), &responseMap)
	if err != nil {
		return utils.GetRequeueResult(), errs
	}

	apikeyEncoded, ok := responseMap["encoded"].(string)
	if !ok {
		fmt.Println("ApikeyEncoded's value conversion failed")
	}

	data := map[string][]byte{
		"apikey": []byte(apikeyEncoded),
	}

	if err := CreateApikeySecret(cli, ctx, req.Namespace, req.Name, data); err != nil {
		return utils.GetRequeueResult(), err
	}

	return ctrl.Result{}, nil
}

func CreateApikeySecret(cli client.Client, ctx context.Context, namespace string, secretName string, data map[string][]byte) error {
	newSecret := &k8sv1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      secretName,
		},
		Data: data,
		Type: k8sv1.SecretTypeOpaque,
	}

	if err := cli.Create(ctx, newSecret); err != nil {
		return err
	}
	return nil
}

func DeleteApikeySecret(cli client.Client, ctx context.Context, namespace string, secretName string) error {
	secret := &k8sv1.Secret{}

	if err := cli.Get(ctx, client.ObjectKey{Namespace: namespace, Name: secretName}, secret); err != nil {
		return err
	}

	if err := cli.Delete(ctx, secret); err != nil {
		return err
	}
	return nil
}
