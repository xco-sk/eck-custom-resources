package kibana

import (
	"fmt"
	kibanaeckv1alpha1 "github.com/xco-sk/eck-custom-resources/apis/kibana.eck/v1alpha1"
	"github.com/xco-sk/eck-custom-resources/utils"
	"io/ioutil"
	"net/http"
	ctrl "sigs.k8s.io/controller-runtime"
)

func DeleteIndexPattern(kClient Client, indexPattern string) (ctrl.Result, error) {
	_, deleteErr := kClient.DoDelete(fmt.Sprintf("/api/saved_objects/index-pattern/%s", indexPattern))
	return ctrl.Result{}, deleteErr
}

func UpsertIndexPattern(kClient Client, indexPattern kibanaeckv1alpha1.IndexPattern) (ctrl.Result, error) {

	exists, err := IndexPatternExists(kClient, indexPattern)
	if err != nil {
		return utils.GetRequeueResult(), err
	}

	var res *http.Response
	if exists {
		res, err = kClient.DoPut(fmt.Sprintf("/api/saved_objects/index-pattern/%s", indexPattern.Name), indexPattern.Spec.Body)
	} else {
		res, err = kClient.DoPost(fmt.Sprintf("/api/saved_objects/index-pattern/%s", indexPattern.Name), indexPattern.Spec.Body)
	}

	if err != nil {
		return utils.GetRequeueResult(), err
	}
	if res.StatusCode > 299 {
		resBody, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return utils.GetRequeueResult(), err
		}
		return utils.GetRequeueResult(), fmt.Errorf("Non-success (%d) response: %s, ", res.StatusCode, string(resBody))
	}

	return ctrl.Result{}, nil
}

func IndexPatternExists(kClient Client, indexPattern kibanaeckv1alpha1.IndexPattern) (bool, error) {
	res, err := kClient.DoGet(fmt.Sprintf("/api/saved_objects/index-pattern/%s", indexPattern.Name))
	return err == nil && res.StatusCode == 200, err
}
