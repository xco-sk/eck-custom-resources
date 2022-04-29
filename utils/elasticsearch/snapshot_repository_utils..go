package elasticsearch

import (
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/xco-sk/eck-custom-resources/apis/es.eck/v1alpha1"
	"github.com/xco-sk/eck-custom-resources/utils"
	ctrl "sigs.k8s.io/controller-runtime"
	"strings"
)

func DeleteSnapshotRepository(esClient *elasticsearch.Client, repositoryName string) (ctrl.Result, error) {
	res, err := esClient.Snapshot.DeleteRepository([]string{repositoryName})
	if err != nil || res.IsError() {
		return utils.GetRequeueResult(), err
	}
	return ctrl.Result{}, nil
}

func UpsertSnapshotRepository(esClient *elasticsearch.Client, snapshotRepository v1alpha1.SnapshotRepository) (ctrl.Result, error) {
	res, err := esClient.Snapshot.GetRepository(
		esClient.Snapshot.GetRepository.WithRepository(snapshotRepository.Name),
	)
	if err != nil {
		return utils.GetRequeueResult(), err
	}

	if res.StatusCode == 404 {
		return createSnapshotRepository(esClient, snapshotRepository)
	}
	return updateSnapshotRepository(esClient, snapshotRepository)
}

func createSnapshotRepository(esClient *elasticsearch.Client, snapshotRepository v1alpha1.SnapshotRepository) (ctrl.Result, error) {
	res, err := esClient.Snapshot.CreateRepository(snapshotRepository.Name, strings.NewReader(snapshotRepository.Spec.Body))

	if err != nil || res.IsError() {
		return utils.GetRequeueResult(), GetClientErrorOrResponseError(err, res)
	}

	return ctrl.Result{}, nil
}

func updateSnapshotRepository(esClient *elasticsearch.Client, snapshotRepository v1alpha1.SnapshotRepository) (ctrl.Result, error) {
	_, repoDeleteErr := DeleteSnapshotRepository(esClient, snapshotRepository.Name)
	if repoDeleteErr != nil {
		return ctrl.Result{}, repoDeleteErr
	}

	res, err := esClient.Snapshot.CreateRepository(snapshotRepository.Name, strings.NewReader(snapshotRepository.Spec.Body))
	if err != nil || res.IsError() {
		return utils.GetRequeueResult(), GetClientErrorOrResponseError(err, res)
	}

	return ctrl.Result{}, nil
}
