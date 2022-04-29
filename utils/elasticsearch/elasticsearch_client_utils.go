package elasticsearch

import (
	"context"
	"github.com/elastic/go-elasticsearch/v8"
	configv2 "github.com/xco-sk/eck-custom-resources/apis/config/v2"
	k8sv1 "k8s.io/api/core/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

var esClient *elasticsearch.Client = nil

func GetElasticUserSecret(cli client.Client, ctx context.Context, namespace string, esSpec configv2.ElasticsearchSpec, secret *k8sv1.Secret) error {
	if err := cli.Get(ctx, client.ObjectKey{Namespace: namespace, Name: esSpec.Authentication.UsernamePassword.SecretName}, secret); err != nil {
		return err
	}
	return nil
}

func GetHttpCertificateSecret(cli client.Client, ctx context.Context, namespace string, esSpec configv2.ElasticsearchSpec, secret *k8sv1.Secret) error {
	if err := cli.Get(ctx, client.ObjectKey{Namespace: namespace, Name: esSpec.Certificate.SecretName}, secret); err != nil {
		return err
	}
	return nil
}

func GetElasticsearchClient(cli client.Client, ctx context.Context, esSpec configv2.ElasticsearchSpec, req ctrl.Request) (*elasticsearch.Client, error) {
	logger := log.FromContext(ctx)

	if esClient == nil {
		logger.Info("Elasticsearch client not initialized, initializing.")

		var userSecret k8sv1.Secret
		if err := GetElasticUserSecret(cli, ctx, req.Namespace, esSpec, &userSecret); err != nil {
			logger.Error(err, "unable to fetch elastic-user secret")
		}

		var certificateSecret k8sv1.Secret
		if err := GetHttpCertificateSecret(cli, ctx, req.Namespace, esSpec, &certificateSecret); err != nil {
			logger.Error(err, "unable to fetch certificate")
		}

		config := elasticsearch.Config{
			Addresses: []string{esSpec.Url},
			Username:  esSpec.Authentication.UsernamePassword.UserName,
			Password:  string(userSecret.Data[esSpec.Authentication.UsernamePassword.UserName]),
			CACert:    certificateSecret.Data[esSpec.Certificate.CertificateKey],
		}

		c, err := elasticsearch.NewClient(config)
		if err != nil {
			return nil, err
		}

		esClient = c
	}

	return esClient, nil
}
