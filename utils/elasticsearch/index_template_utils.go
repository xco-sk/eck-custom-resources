package elasticsearch

import (
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/xco-sk/eck-custom-resources/apis/es.eck/v1alpha1"
	"github.com/xco-sk/eck-custom-resources/utils"
	ctrl "sigs.k8s.io/controller-runtime"
	"strings"
)

func DeleteIndexTemplate(esClient *elasticsearch.Client, indexTemplateName string) (ctrl.Result, error) {
	res, err := esClient.Indices.DeleteTemplate(indexTemplateName)
	if err != nil || res.IsError() {
		return utils.GetRequeueResult(), err
	}
	return ctrl.Result{}, nil
}

func UpsertIndexTemplate(esClient *elasticsearch.Client, indexTemplate v1alpha1.IndexTemplate) (ctrl.Result, error) {
	res, err := esClient.Indices.PutIndexTemplate(indexTemplate.Name, strings.NewReader(indexTemplate.Spec.Body))

	if err != nil || res.IsError() {
		return utils.GetRequeueResult(), err
	}

	return ctrl.Result{}, nil
}
