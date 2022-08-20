package kibana

import (
	"fmt"
	kibanaeckv1alpha1 "github.com/xco-sk/eck-custom-resources/apis/kibana.eck/v1alpha1"
	"github.com/xco-sk/eck-custom-resources/utils"
	"io/ioutil"
	"net/http"
	ctrl "sigs.k8s.io/controller-runtime"
)

func DeleteDataView(kClient Client, dataViewName string) (ctrl.Result, error) {
	_, deleteErr := kClient.DoDelete(fmt.Sprintf("/api/data_views/data_view/%s", dataViewName))
	return ctrl.Result{}, deleteErr
}

func UpsertDataView(kClient Client, dataView kibanaeckv1alpha1.DataView) (ctrl.Result, error) {

	exists, err := DataViewExists(kClient, dataView.Name)
	if err != nil {
		return utils.GetRequeueResult(), err
	}

	var res *http.Response

	if exists {
		res, err = kClient.DoPut(fmt.Sprintf("/api/data_views/data_view/%s", dataView.Name), dataView.Spec.Body)
	} else {
		res, err = kClient.DoPost("/api/data_views/dataView", dataView.Spec.Body)
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

func DataViewExists(kClient Client, dataViewName string) (bool, error) {
	res, err := kClient.DoGet(fmt.Sprintf("/api/data_views/data_view/%s", dataViewName))
	return err == nil && res.StatusCode == 200, err
}
