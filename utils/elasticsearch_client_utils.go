package utils

import (
	"context"
	b64 "encoding/base64"
	"github.com/elastic/go-elasticsearch/v8"
	eseckv1 "github.com/xco-sk/eck-custom-resources/api/v1alpha1"
	k8sv1 "k8s.io/api/core/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

const EckUserSecretSuffix = "-es-elastic-user"
const EckHttpServiceSuffix = "-es-http"

func GetElasticSecret(cli client.Client, ctx context.Context, namespace string, esSpec eseckv1.ElasticsearchSpec, secret *k8sv1.Secret) error {
	if err := cli.Get(ctx, client.ObjectKey{Namespace: namespace, Name: generateElasticSecretName(esSpec.EckCluster.ClusterName)}, secret); err != nil {
		return err
	}
	return nil
}

func GenerateElasticEndpoint(clusterName string, namespace string) string {
	return "https://" + clusterName + EckHttpServiceSuffix + "." + namespace + ":9200"
}

func generateElasticSecretName(clusterName string) string {
	return clusterName + EckUserSecretSuffix
}

func GetElasticsearchClient(cli client.Client, ctx context.Context, esSpec eseckv1.ElasticsearchSpec, req ctrl.Request) (*elasticsearch.Client, error) {
	logger := log.FromContext(ctx)

	var secret k8sv1.Secret
	if err := GetElasticSecret(cli, ctx, req.Namespace, esSpec, &secret); err != nil {
		logger.Error(err, "unable to fetch elastic-user secret")
	}

	decodedSecret, _ := b64.StdEncoding.DecodeString(string(secret.Data["elastic"]))

	config := elasticsearch.Config{
		Addresses: []string{GenerateElasticEndpoint(esSpec.EckCluster.ClusterName, req.Namespace)},
		Username:  "elastic",
		Password:  string(decodedSecret),
	}

	return elasticsearch.NewClient(config)
}
