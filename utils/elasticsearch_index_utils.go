package utils

import (
	"encoding/json"
	"github.com/elastic/go-elasticsearch/v8"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

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
		res, deleteErr := esClient.Indices.Delete([]string{indexName})
		if deleteErr != nil || res.IsError() {
			return GetRequeueResult(), deleteErr
		}
	} else {
		logger.Info("Index not empty, skipping deletion", "Index name", indexName)
		return ctrl.Result{}, nil
	}
	return ctrl.Result{}, nil
}

// TODO rename/refactor
func RecreateIndexIfEmpty(esClient *elasticsearch.Client, req ctrl.Request) (ctrl.Result, error, bool) {
	var logger = log.Log

	indexEmpty, indexEmptyErr := VerifyIndexEmpty(esClient, req.Name)
	if indexEmptyErr != nil {
		return ctrl.Result{}, indexEmptyErr, true
	}

	if indexEmpty {
		res, deleteErr := esClient.Indices.Delete([]string{req.Name})
		if deleteErr != nil || res.IsError() {
			return ctrl.Result{}, deleteErr, true
		}
		logger.Info("Recreating index")
	} else {
		logger.Info("Index already exists and is not empty, doing nothing.", "Index name", req.Name)
		return ctrl.Result{}, nil, true
	}
	return ctrl.Result{}, nil, false
}
