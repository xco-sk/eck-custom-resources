package kibana

import (
	"fmt"
	kibanaeckv1alpha1 "github.com/xco-sk/eck-custom-resources/apis/kibana.eck/v1alpha1"
	"github.com/xco-sk/eck-custom-resources/utils"
	"io/ioutil"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
	ctrl "sigs.k8s.io/controller-runtime"
	"strings"
)

func DeleteSavedObject(kClient Client, savedObjectType string, savedObjectMeta metav1.ObjectMeta, savedObject kibanaeckv1alpha1.SavedObject) (ctrl.Result, error) {
	_, deleteErr := kClient.DoDelete(formatSavedObjectUrl(savedObjectType, savedObjectMeta.Name, savedObject.Space))
	return ctrl.Result{}, deleteErr
}

func UpsertSavedObject(kClient Client, savedObjectType string, savedObjectMeta metav1.ObjectMeta, savedObject kibanaeckv1alpha1.SavedObject) (ctrl.Result, error) {

	exists, err := SavedObjectExists(kClient, savedObjectType, savedObjectMeta.Name, savedObject.Space)
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

func SavedObjectExists(kClient Client, savedObjectType string, name string, space *string) (bool, error) {
	res, err := kClient.DoGet(formatSavedObjectUrl(savedObjectType, name, space))
	return err == nil && res.StatusCode == 200, err
}

func DependenciesFulfilled(kClient Client, savedObject kibanaeckv1alpha1.SavedObject) error {

	var missingDependencies []string
	var errors []string

	for _, dependency := range savedObject.Dependencies {
		dSpace := savedObject.Space
		if dependency.Space != nil {
			dSpace = dependency.Space
		}

		exists, err := SavedObjectExists(kClient, string(dependency.ObjectType), dependency.Name, dSpace)

		if err != nil {
			errors = append(errors, err.Error())
			continue
		}
		if !exists {
			spaceInfoMessage := ""
			if dSpace != nil {
				spaceInfoMessage = fmt.Sprintf(" (in space: %s)", *dSpace)
			}
			missingDependencies = append(missingDependencies, fmt.Sprintf("%s/%s%s", dependency.ObjectType, dependency.Name, spaceInfoMessage))
		}
	}

	if len(missingDependencies) > 0 || len(errors) > 0 {
		return fmt.Errorf("dependencies not fulfilled. Missing dependencies:[%s]. Errors:[%s]",
			strings.Join(missingDependencies[:], ","),
			strings.Join(errors[:], ","))
	}
	return nil
}

func formatSavedObjectUrl(savedObjectType string, name string, space *string) string {
	if space == nil {
		return fmt.Sprintf("/api/saved_objects/%s/%s", savedObjectType, name)
	}
	return fmt.Sprintf("/s/%s/api/saved_objects/%s/%s", *space, savedObjectType, name)
}
