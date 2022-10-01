package kibana

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	kibanaeckv1alpha1 "github.com/xco-sk/eck-custom-resources/apis/kibana.eck/v1alpha1"
	"github.com/xco-sk/eck-custom-resources/utils"
	ctrl "sigs.k8s.io/controller-runtime"
)

const OVERRIDE = true
const REFRESH_FIELDS = true

func DeleteDataView(kClient Client, dataView kibanaeckv1alpha1.DataView) (ctrl.Result, error) {
	_, deleteErr := kClient.DoDelete(formatExistingDataViewUrl(dataView.Name, dataView.Spec.Space))
	return ctrl.Result{}, deleteErr
}

func UpsertDataView(kClient Client, dataView kibanaeckv1alpha1.DataView) (ctrl.Result, error) {

	exists, err := DataViewExists(kClient, dataView)
	if err != nil {
		return utils.GetRequeueResult(), err
	}

	var res *http.Response

	modifiedBody, err := wrapDataView(dataView, exists)
	if err != nil {
		return ctrl.Result{}, err
	}

	if exists {
		res, err = kClient.DoPost(formatExistingDataViewUrl(dataView.Name, dataView.Spec.Space), *modifiedBody)
	} else {
		res, err = kClient.DoPost(formatDataViewUrl(dataView.Spec.Space), *modifiedBody)
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

func DataViewExists(kClient Client, dataView kibanaeckv1alpha1.DataView) (bool, error) {
	res, err := kClient.DoGet(formatExistingDataViewUrl(dataView.Name, dataView.Spec.Space))
	return err == nil && res.StatusCode == 200, err
}

func formatExistingDataViewUrl(name string, space *string) string {
	return fmt.Sprintf("%s/%s", formatDataViewUrl(space), name)
}

func formatDataViewUrl(space *string) string {
	if space == nil {
		return "/api/data_views/data_view"
	}
	return fmt.Sprintf("/s/%s/api/data_views/data_view", *space)
}

func wrapDataView(dataView kibanaeckv1alpha1.DataView, isUpdate bool) (*string, error) {
	var err error

	dataViewString := &dataView.Spec.Body

	if !isUpdate {
		dataViewString, err = InjectId(*dataViewString, dataView.Name)
		if err != nil {
			return nil, err
		}
	} else {
		dataViewString, err = removeName(*dataViewString, dataView.Name)
		if err != nil {
			return nil, err
		}
	}

	var dataViewMap map[string]interface{}
	err = json.NewDecoder(strings.NewReader(*dataViewString)).Decode(&dataViewMap)
	if err != nil {
		return nil, err
	}

	body := make(map[string]interface{})

	body["data_view"] = dataViewMap
	if isUpdate {
		body["refresh_fields"] = REFRESH_FIELDS
	} else {
		body["override"] = OVERRIDE
	}

	marshalledBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	sBody := string(marshalledBody)
	return &sBody, nil
}

func removeName(objectJson string, id string) (*string, error) {
	var body map[string]interface{}
	err := json.NewDecoder(strings.NewReader(objectJson)).Decode(&body)
	if err != nil {
		return nil, err
	}

	if _, exists := body["name"]; exists {
		delete(body, "name")
	}

	marshalledBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	sBody := string(marshalledBody)
	return &sBody, nil
}
