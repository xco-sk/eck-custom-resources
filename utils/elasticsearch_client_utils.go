package utils

import (
	"context"
	"github.com/elastic/go-elasticsearch/v8"
	configv2 "github.com/xco-sk/eck-custom-resources/apis/config/v2"
	"strconv"

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

func VerifyIndexExists(esClient *elasticsearch.Client, indexName string) (bool, error) {
	existsResponse, err := esClient.Indices.Exists([]string{indexName})
	if err != nil {
		return false, err
	}
	return existsResponse.StatusCode == 200, nil
}

func VerifyIndexEmpty(esClient *elasticsearch.Client, indexName string) (bool, error) {
	var countResponse, countErr = esClient.Count(
		esClient.Count.WithIndex(indexName),
	)
	var docsCount []byte
	if countErr != nil {
		return false, countErr
	}

	_, readErr := countResponse.Body.Read(docsCount)
	if readErr != nil {
		return false, readErr
	}

	count, conversionErr := strconv.Atoi(string(docsCount))
	countResponse.Body.Close()
	if conversionErr != nil {
		return false, conversionErr
	}

	return count == 0, nil
}

func DeleteIndexIfEmpty(esClient *elasticsearch.Client, indexName string) (ctrl.Result, error) {
	var logger = log.Log
	indexExists, existsErr := VerifyIndexExists(esClient, indexName)
	if existsErr != nil {
		logger.Error(existsErr, "Failed to verify if index exists")
		return GetRequeueResult(), existsErr
	}

	if !indexExists {
		logger.Info("Index does not exists, do nothing.", "Index name", indexName)
		return ctrl.Result{}, nil
	}

	isEmpty, emptyErr := VerifyIndexEmpty(esClient, indexName)
	if emptyErr != nil {
		logger.Error(emptyErr, "Failed to verify if index is empty")
		return GetRequeueResult(), emptyErr
	}

	if isEmpty {
		_, deleteErr := esClient.Indices.Delete([]string{indexName})
		if deleteErr != nil {
			return GetRequeueResult(), deleteErr
		}
	} else {
		logger.Info("Index not empty, skipping deletion", "Index name", indexName)
		return ctrl.Result{}, nil
	}
	return ctrl.Result{}, nil
}

func RecreateIndexIfEmpty(esClient *elasticsearch.Client, req ctrl.Request) (ctrl.Result, error, bool) {
	var logger = log.Log

	indexEmpty, indexEmptyErr := VerifyIndexEmpty(esClient, req.Name)
	if indexEmptyErr != nil {
		return ctrl.Result{}, indexEmptyErr, true
	}

	if indexEmpty {
		_, deleteErr := esClient.Indices.Delete([]string{req.Name})
		if deleteErr != nil {
			return ctrl.Result{}, deleteErr, true
		}
		logger.Info("Recreating index")
	} else {
		logger.Info("Index already exists and is not empty, doing nothing.", "Index name", req.Name)
		return ctrl.Result{}, nil, true
	}
	return ctrl.Result{}, nil, false
}
