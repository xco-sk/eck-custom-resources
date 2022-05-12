/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package kibanaeck

import (
	"context"
	"fmt"
	configv2 "github.com/xco-sk/eck-custom-resources/apis/config/v2"
	kibanaUtils "github.com/xco-sk/eck-custom-resources/utils/kibana"
	"k8s.io/client-go/tools/record"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	kibanaeckv1alpha1 "github.com/xco-sk/eck-custom-resources/apis/kibana.eck/v1alpha1"
)

// IndexPatternReconciler reconciles a IndexPattern object
type IndexPatternReconciler struct {
	client.Client
	Scheme        *runtime.Scheme
	ProjectConfig configv2.ProjectConfig
	Recorder      record.EventRecorder
}

//+kubebuilder:rbac:groups=kibana.eck.github.com,resources=indexpatterns,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=kibana.eck.github.com,resources=indexpatterns/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=kibana.eck.github.com,resources=indexpatterns/finalizers,verbs=update

func (r *IndexPatternReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	kibanaClient := kibanaUtils.Client{
		Cli:        r.Client,
		Ctx:        ctx,
		KibanaSpec: r.ProjectConfig.Kibana,
		Req:        req,
	}

	var indexPattern kibanaeckv1alpha1.IndexPattern
	if err := r.Get(ctx, req.NamespacedName, &indexPattern); err != nil {
		logger.Info("Deleting index pattern", "role", req.Name)
		return kibanaUtils.DeleteIndexPattern(kibanaClient, req.Name)
	}

	logger.Info("Creating/Updating index pattern", "role", req.Name)
	res, err := kibanaUtils.UpsertIndexPattern(kibanaClient, indexPattern)

	if err == nil {
		r.Recorder.Event(&indexPattern, "Normal", "Created",
			fmt.Sprintf("Created/Updated %s/%s %s", indexPattern.APIVersion, indexPattern.Kind, indexPattern.Name))
	} else {
		r.Recorder.Event(&indexPattern, "Warning", "Failed to create/update",
			fmt.Sprintf("Failed to create/update %s/%s %s: %s", indexPattern.APIVersion, indexPattern.Kind, indexPattern.Name, err.Error()))
	}

	return res, err
}

// SetupWithManager sets up the controller with the Manager.
func (r *IndexPatternReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&kibanaeckv1alpha1.IndexPattern{}).
		Complete(r)
}
