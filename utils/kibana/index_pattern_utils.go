package kibana

import (
	"fmt"
	kibanaeckv1alpha1 "github.com/xco-sk/eck-custom-resources/apis/kibana.eck/v1alpha1"
	ctrl "sigs.k8s.io/controller-runtime"
)

func DeleteIndexPattern(kClient Client, indexPattern string) (ctrl.Result, error) {
	_, deleteErr := kClient.DoDelete(fmt.Sprintf("/api/saved_objects/index-pattern/%s", indexPattern))
	return ctrl.Result{}, deleteErr
}

func UpsertIndexPattern(kClient Client, indexPattern kibanaeckv1alpha1.IndexPattern) (ctrl.Result, error) {

	return ctrl.Result{}, nil
}
