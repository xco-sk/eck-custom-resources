package utils

import (
	"context"
	"github.com/elastic/go-elasticsearch/v8"
	configv2 "github.com/xco-sk/eck-custom-resources/apis/config/v2"
	k8sv1 "k8s.io/api/core/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

const eckUserSecretSuffix = "-es-elastic-user"
const eckHttpCertificateSecretSuffix = "-es-http-certs-public"
const eckHttpCertificateCAKey = "ca.crt"
const eckHttpServiceSuffix = "-es-http"

func GetElasticUserSecret(cli client.Client, ctx context.Context, namespace string, esSpec configv2.ElasticsearchSpec, secret *k8sv1.Secret) error {
	if err := cli.Get(ctx, client.ObjectKey{Namespace: namespace, Name: generateElasticUserSecretName(esSpec.EckCluster.ClusterName)}, secret); err != nil {
		return err
	}
	return nil
}

func GetHttpCertificateSecret(cli client.Client, ctx context.Context, namespace string, esSpec configv2.ElasticsearchSpec, secret *k8sv1.Secret) error {
	if err := cli.Get(ctx, client.ObjectKey{Namespace: namespace, Name: generateHttpCertificateSecretName(esSpec.EckCluster.ClusterName)}, secret); err != nil {
		return err
	}
	return nil
}

func GenerateElasticEndpoint(clusterName string, namespace string) string {
	return "https://" + clusterName + eckHttpServiceSuffix + ":9200"
	//return "https://" + clusterName + EckHttpServiceSuffix + "." + namespace + ":9200"
}

func generateElasticUserSecretName(clusterName string) string {
	return clusterName + eckUserSecretSuffix
}

func generateHttpCertificateSecretName(clusterName string) string {
	return clusterName + eckHttpCertificateSecretSuffix
}

func GetElasticsearchClient(cli client.Client, ctx context.Context, esSpec configv2.ElasticsearchSpec, req ctrl.Request) (*elasticsearch.Client, error) {
	logger := log.FromContext(ctx)

	var userSecret k8sv1.Secret
	if err := GetElasticUserSecret(cli, ctx, req.Namespace, esSpec, &userSecret); err != nil {
		logger.Error(err, "unable to fetch elastic-user secret")
	}

	var certificateSecret k8sv1.Secret
	if err := GetHttpCertificateSecret(cli, ctx, req.Namespace, esSpec, &certificateSecret); err != nil {
		logger.Error(err, "unable to fetch certificate")
	}

	config := elasticsearch.Config{
		Addresses: []string{GenerateElasticEndpoint(esSpec.EckCluster.ClusterName, req.Namespace)},
		Username:  "elastic",
		Password:  string(userSecret.Data["elastic"]),
		CACert:    certificateSecret.Data[eckHttpCertificateCAKey],
	}

	return elasticsearch.NewClient(config)
}
