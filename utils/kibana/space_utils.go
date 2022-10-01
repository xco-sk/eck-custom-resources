package kibana

import (
	"fmt"
	"io/ioutil"
	"net/http"

	kibanaeckv1alpha1 "github.com/xco-sk/eck-custom-resources/apis/kibana.eck/v1alpha1"
	"github.com/xco-sk/eck-custom-resources/utils"
	ctrl "sigs.k8s.io/controller-runtime"
)

func DeleteSpace(kClient Client, spaceName string) (ctrl.Result, error) {
	_, deleteErr := kClient.DoDelete(fmt.Sprintf("/api/spaces/space/%s", spaceName))
	return ctrl.Result{}, deleteErr
}

func UpsertSpace(kClient Client, space kibanaeckv1alpha1.Space) (ctrl.Result, error) {

	exists, err := SpaceExists(kClient, space.Name)
	if err != nil {
		return utils.GetRequeueResult(), err
	}

	var res *http.Response

	modifiedBody, err := InjectId(space.Spec.Body, space.Name)
	if err != nil {
		return ctrl.Result{}, err
	}

	if exists {
		res, err = kClient.DoPut(fmt.Sprintf("/api/spaces/space/%s", space.Name), *modifiedBody)
	} else {
		res, err = kClient.DoPost("/api/spaces/space", *modifiedBody)
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

func SpaceExists(kClient Client, spaceName string) (bool, error) {
	res, err := kClient.DoGet(fmt.Sprintf("/api/spaces/space/%s", spaceName))
	return err == nil && res.StatusCode == 200, err
}
