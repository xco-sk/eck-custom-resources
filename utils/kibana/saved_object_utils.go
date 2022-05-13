package kibana

import (
	"fmt"
	kibanaeckv1alpha1 "github.com/xco-sk/eck-custom-resources/apis/kibana.eck/v1alpha1"
	"github.com/xco-sk/eck-custom-resources/utils"
	"io/ioutil"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
	ctrl "sigs.k8s.io/controller-runtime"
)

func DeleteSavedObject(kClient Client, savedObjectType string, savedObjectMeta metav1.ObjectMeta, savedObject kibanaeckv1alpha1.SavedObject) (ctrl.Result, error) {
	_, deleteErr := kClient.DoDelete(formatSavedObjectUrl(savedObjectType, savedObjectMeta.Name, savedObject.Space))
	return ctrl.Result{}, deleteErr
}

func UpsertSavedObject(kClient Client, savedObjectType string, savedObjectMeta metav1.ObjectMeta, savedObject kibanaeckv1alpha1.SavedObject) (ctrl.Result, error) {

	exists, err := SavedObjectExists(kClient, savedObjectType, savedObjectMeta, savedObject)
	if err != nil {
		return utils.GetRequeueResult(), err
	}

	var res *http.Response
	if exists {
		res, err = kClient.DoPut(formatSavedObjectUrl(savedObjectType, savedObjectMeta.Name, savedObject.Space), savedObject.Body)
	} else {
		res, err = kClient.DoPost(formatSavedObjectUrl(savedObjectType, savedObjectMeta.Name, savedObject.Space), savedObject.Body)
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

func SavedObjectExists(kClient Client, savedObjectType string, savedObjectMeta metav1.ObjectMeta, savedObject kibanaeckv1alpha1.SavedObject) (bool, error) {
	res, err := kClient.DoGet(formatSavedObjectUrl(savedObjectType, savedObjectMeta.Name, savedObject.Space))
	return err == nil && res.StatusCode == 200, err
}

func formatSavedObjectUrl(savedObjectType string, name string, space *string) string {
	if space == nil {
		return fmt.Sprintf("/api/saved_objects/%s/%s", savedObjectType, name)
	}
	return fmt.Sprintf("/s/%s/api/%s/%s", *space, savedObjectType, name)
}
