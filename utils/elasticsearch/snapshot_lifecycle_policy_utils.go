package elasticsearch

import (
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/xco-sk/eck-custom-resources/apis/es.eck/v1alpha1"
	"github.com/xco-sk/eck-custom-resources/utils"
	ctrl "sigs.k8s.io/controller-runtime"
	"strings"
)

func DeleteSnapshotLifecyclePolicy(esClient *elasticsearch.Client, snapshotLifecyclePolicyName string) (ctrl.Result, error) {
	res, err := esClient.SlmDeleteLifecycle(snapshotLifecyclePolicyName)
	if err != nil || res.IsError() {
		return utils.GetRequeueResult(), err
	}
	return ctrl.Result{}, nil
}

func UpsertSnapshotLifecyclePolicy(esClient *elasticsearch.Client, snapshotLifecyclePolicy v1alpha1.SnapshotLifecyclePolicy) (ctrl.Result, error) {
	res, err := esClient.SlmPutLifecycle(snapshotLifecyclePolicy.Name, esClient.SlmPutLifecycle.WithBody(strings.NewReader(snapshotLifecyclePolicy.Spec.Body)))
	if err != nil || res.IsError() {
		return utils.GetRequeueResult(), GetClientErrorOrResponseError(err, res)
	}
	return ctrl.Result{}, nil
}
