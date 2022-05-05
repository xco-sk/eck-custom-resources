package elasticsearch

import (
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/xco-sk/eck-custom-resources/apis/es.eck/v1alpha1"
	"github.com/xco-sk/eck-custom-resources/utils"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"strings"
)

var UpdatableSettings = [...]string{
	"number_of_replicas",
	"refresh_interval",
}

func VerifyIndexExists(esClient *elasticsearch.Client, indexName string) (bool, error) {
	existsResponse, err := esClient.Indices.Exists([]string{indexName})
	if err != nil {
		return false, err
	}
	if existsResponse.StatusCode <= 299 {
		return true, nil
	}
	if existsResponse.StatusCode == 404 {
		return false, nil
	}

	return false, GetClientErrorOrResponseError(nil, existsResponse)
}

func VerifyIndexEmpty(esClient *elasticsearch.Client, indexName string) (bool, error) {
	var countResponse, countErr = esClient.Count(
		esClient.Count.WithIndex(indexName),
	)
	if countErr != nil {
		return false, countErr
	}

	var jsonResponse map[string]interface{}
	decodeErr := json.NewDecoder(countResponse.Body).Decode(&jsonResponse)
	if decodeErr != nil {
		return false, decodeErr
	}

	responseCloseErr := countResponse.Body.Close()
	if responseCloseErr != nil {
		return false, responseCloseErr
	}

	return int(jsonResponse["count"].(float64)) == 0, nil
}

func DeleteIndexIfEmpty(esClient *elasticsearch.Client, indexName string) (ctrl.Result, error) {
	var logger = log.Log
	indexExists, existsErr := VerifyIndexExists(esClient, indexName)
	if existsErr != nil {
		logger.Error(existsErr, "Failed to verify if index exists")
		return utils.GetRequeueResult(), existsErr
	}

	if !indexExists {
		logger.Info("Index does not exists, do nothing.", "Index name", indexName)
		return ctrl.Result{}, nil
	}

	isEmpty, emptyErr := VerifyIndexEmpty(esClient, indexName)
	if emptyErr != nil {
		logger.Error(emptyErr, "Failed to verify if index is empty")
		return utils.GetRequeueResult(), emptyErr
	}

	if isEmpty {
		return DeleteIndex(esClient, indexName)
	}

	logger.Info("Index not empty, skipping deletion", "Index name", indexName)
	return ctrl.Result{}, nil
}

func DeleteIndex(esClient *elasticsearch.Client, indexName string) (ctrl.Result, error) {
	res, deleteErr := esClient.Indices.Delete([]string{indexName})
	if deleteErr != nil || res.IsError() {
		return utils.GetRequeueResult(), deleteErr
	}
	return ctrl.Result{}, nil
}

func CreateIndex(esClient *elasticsearch.Client, index v1alpha1.Index) (ctrl.Result, error) {
	res, err := esClient.Indices.Create(index.Name,
		esClient.Indices.Create.WithBody(strings.NewReader(index.Spec.Body)),
	)

	if err != nil || res.IsError() {
		return utils.GetRequeueResult(), GetClientErrorOrResponseError(err, res)
	}

	return ctrl.Result{}, nil
}

func UpdateIndex(esClient *elasticsearch.Client, index v1alpha1.Index, eventRecorder record.EventRecorder) (ctrl.Result, error) {
	var updatedBody map[string]interface{}
	err := json.NewDecoder(strings.NewReader(index.Spec.Body)).Decode(&updatedBody)
	if err != nil {
		return ctrl.Result{}, err
	}

	whitelistedUpdatedBody := make(map[string]interface{})
	for _, updatable := range UpdatableSettings {
		whitelistedUpdatedBody[updatable] = updatedBody["settings"].(map[string]interface{})[updatable]
	}

	marshalledSettings, err := json.Marshal(whitelistedUpdatedBody)
	if err != nil {
		return ctrl.Result{}, err
	}
	settingsRes, settingsErr := esClient.Indices.PutSettings(
		strings.NewReader(string(marshalledSettings)),
		esClient.Indices.PutSettings.WithIndex(index.Name),
	)
	if settingsErr != nil || settingsRes.IsError() {
		return utils.GetRequeueResult(), GetClientErrorOrResponseError(settingsErr, settingsRes)
	}
	eventRecorder.Event(&index, "Normal", "Index settings updated", fmt.Sprintf("Index settings successfuly updated for %s", index.Name))

	marshalledMapping, err := json.Marshal(updatedBody["mappings"])
	if err != nil {
		return ctrl.Result{}, err
	}
	mappingRes, mappingErr := esClient.Indices.PutMapping(
		[]string{index.Name},
		strings.NewReader(string(marshalledMapping)),
	)
	if mappingErr != nil || mappingRes.IsError() {
		return utils.GetRequeueResult(), GetClientErrorOrResponseError(mappingErr, mappingRes)
	}
	eventRecorder.Event(&index, "Normal", "Index mapping updated", fmt.Sprintf("Index mapping successfuly updated for %s", index.Name))

	return ctrl.Result{}, nil
}
