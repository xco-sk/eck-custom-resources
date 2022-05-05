package elasticsearch

import (
	"context"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	configv2 "github.com/xco-sk/eck-custom-resources/apis/config/v2"
	"github.com/xco-sk/eck-custom-resources/apis/es.eck/v1alpha1"
	"io"
	k8sv1 "k8s.io/api/core/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"strings"
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

func GetClientErrorOrResponseError(err error, response *esapi.Response) error {
	if err != nil {
		return err
	}

	body, readErr := io.ReadAll(response.Body)
	if readErr != nil {
		return readErr
	}

	return fmt.Errorf("error response: %s", body)
}

func DependenciesFulfilled(esClient *elasticsearch.Client, dependencies v1alpha1.Dependencies) error {

	var missingIdxTemplates []string
	var missingIdx []string
	var errors []string

	for _, idxTplDependency := range dependencies.IndexTemplate {
		exists, err := IndexTemplateExists(esClient, idxTplDependency.IndexTemplateName)
		if err != nil {
			errors = append(errors, err.Error())
			continue
		}
		if !exists {
			missingIdxTemplates = append(missingIdxTemplates, idxTplDependency.IndexTemplateName)
		}
	}
	for _, idxDependency := range dependencies.Index {
		exists, err := VerifyIndexExists(esClient, idxDependency.IndexName)
		if err != nil {
			errors = append(errors, err.Error())
			continue
		}
		if !exists {
			missingIdx = append(missingIdx, idxDependency.IndexName)
		}
	}

	if len(missingIdx) > 0 || len(missingIdxTemplates) > 0 || len(errors) > 0 {
		return fmt.Errorf("dependencies not fulfilled. Missing indices: %s. Missing index templates: %s. Errors: %s",
			strings.Join(missingIdx[:], ","),
			strings.Join(missingIdxTemplates[:], ","),
			strings.Join(errors[:], ","))
	}
	return nil
}
