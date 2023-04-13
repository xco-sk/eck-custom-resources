package elasticsearch

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/xco-sk/eck-custom-resources/apis/es.eck/v1alpha1"
	"github.com/xco-sk/eck-custom-resources/utils"
	k8sv1 "k8s.io/api/core/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func DeleteUser(esClient *elasticsearch.Client, userName string) (ctrl.Result, error) {
	res, err := esClient.Security.DeleteUser(userName)
	if err != nil || res.IsError() {
		return utils.GetRequeueResult(), err
	}
	return ctrl.Result{}, nil
}

func UpsertUser(esClient *elasticsearch.Client, cli client.Client, ctx context.Context, user v1alpha1.ElasticsearchUser) (ctrl.Result, error) {
	var secret k8sv1.Secret

	// Inject password field with data from given secret
	err := getUserSecret(cli, ctx, user.Namespace, user, &secret)
	if err != nil {
		return utils.GetRequeueResult(), err
	}
	var password = secret.Data[user.Name]

	var userBody map[string]interface{}
	unmarshallErr := json.Unmarshal([]byte(user.Spec.Body), &userBody)
	if unmarshallErr != nil {
		return ctrl.Result{}, unmarshallErr
	}

	userBody["password"] = string(password)
	userWithPassword, marshallErr := json.Marshal(userBody)
	if marshallErr != nil {
		return ctrl.Result{}, marshallErr
	}

	res, err := esClient.Security.PutUser(user.Name, strings.NewReader(string(userWithPassword)))
	if err != nil || res.IsError() {
		return utils.GetRequeueResult(), GetClientErrorOrResponseError(err, res)
	}
	return ctrl.Result{}, nil
}

func getUserSecret(cli client.Client, ctx context.Context, namespace string, user v1alpha1.ElasticsearchUser, secret *k8sv1.Secret) error {
	if err := cli.Get(ctx, client.ObjectKey{Namespace: namespace, Name: user.Spec.SecretName}, secret); err != nil {
		return err
	}
	return nil
}
