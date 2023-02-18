package elasticsearch

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/elastic/elastic-transport-go/v8/elastictransport"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	configv2 "github.com/xco-sk/eck-custom-resources/apis/config/v2"
	"github.com/xco-sk/eck-custom-resources/apis/es.eck/v1alpha1"
	eseckv1alpha1 "github.com/xco-sk/eck-custom-resources/apis/es.eck/v1alpha1"
	"github.com/xco-sk/eck-custom-resources/utils"
	k8sv1 "k8s.io/api/core/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

func GetElasticsearchClient(cli client.Client, ctx context.Context, esSpec configv2.ElasticsearchSpec, req ctrl.Request) (*elasticsearch.Client, error) {
	logger := log.FromContext(ctx)

	logger.Info("Elasticsearch client not initialized, initializing.", "Spec", esSpec)

	config := elasticsearch.Config{
		Addresses:         []string{esSpec.Url},
		EnableDebugLogger: true,
		Logger:            &elastictransport.TextLogger{Output: os.Stdout},
	}

	if esSpec.Authentication != nil && esSpec.Authentication.UsernamePassword != nil {
		var userSecret k8sv1.Secret
		if err := utils.GetUserSecret(cli, ctx, req.Namespace, esSpec.Authentication.UsernamePassword, &userSecret); err != nil {
			return nil, err
		}
		config.Username = esSpec.Authentication.UsernamePassword.UserName
		config.Password = string(userSecret.Data[esSpec.Authentication.UsernamePassword.UserName])
	}

	if esSpec.Certificate != nil {
		var certificateSecret k8sv1.Secret
		if err := utils.GetCertificateSecret(cli, ctx, req.Namespace, esSpec.Certificate, &certificateSecret); err != nil {
			return nil, err
		}
		config.CACert = certificateSecret.Data[esSpec.Certificate.CertificateKey]
	}

	esClient, err := elasticsearch.NewClient(config)
	if err != nil {
		return nil, err
	}

	return esClient, nil
}

func GetClientErrorOrResponseError(err error, response *esapi.Response) error {
	if err != nil {
		return err
	}
	return fmt.Errorf("error(status: %d, response: %s)", response.StatusCode, response.String())
}

func DependenciesFulfilled(esClient *elasticsearch.Client, dependencies v1alpha1.Dependencies) error {

	var missingIdxTemplates []string
	var missingIdx []string
	var errors []string

	for _, idxTplDependency := range dependencies.IndexTemplates {
		exists, err := IndexTemplateExists(esClient, idxTplDependency)
		if err != nil {
			errors = append(errors, err.Error())
			continue
		}
		if !exists {
			missingIdxTemplates = append(missingIdxTemplates, idxTplDependency)
		}
	}
	for _, idxDependency := range dependencies.Indices {
		exists, err := VerifyIndexExists(esClient, idxDependency)
		if err != nil {
			errors = append(errors, err.Error())
			continue
		}
		if !exists {
			missingIdx = append(missingIdx, idxDependency)
		}
	}

	if len(missingIdx) > 0 || len(missingIdxTemplates) > 0 || len(errors) > 0 {
		return fmt.Errorf("dependencies not fulfilled. Missing indices:[%s]. Missing index templates:[%s]. Errors:[%s]",
			strings.Join(missingIdx[:], ","),
			strings.Join(missingIdxTemplates[:], ","),
			strings.Join(errors[:], ","))
	}
	return nil
}

func GetTargetElasticsearchInstance(cli client.Client, ctx context.Context, namespace string, targetName string, esInstance *eseckv1alpha1.ElasticsearchInstance) error {
	if err := cli.Get(ctx, client.ObjectKey{Namespace: namespace, Name: targetName}, esInstance); err != nil {
		return err
	}
	return nil
}
