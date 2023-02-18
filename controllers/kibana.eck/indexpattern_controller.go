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
	"github.com/xco-sk/eck-custom-resources/utils"
	kibanaUtils "github.com/xco-sk/eck-custom-resources/utils/kibana"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

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

	indexPatternFinalizer := "indexpatterns.kibana.eck.github.com/finalizer"
	savedObjectType := "index-pattern"

	var indexPattern kibanaeckv1alpha1.IndexPattern
	if err := r.Get(ctx, req.NamespacedName, &indexPattern); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	targetInstance, err := r.getTargetInstance(&indexPattern, indexPattern.Spec.TargetConfig, ctx, req.Namespace)
	if err != nil {
		return utils.GetRequeueResult(), err
	}

	if !targetInstance.Enabled {
		logger.Info("Kibana reconciler disabled, not reconciling.", "Resource", req.NamespacedName)
		return ctrl.Result{}, nil
	}

	// Get the ElasticsearchInstance defined in target (if present and pass to the kibanaUtils.Client)
	kibanaClient := kibanaUtils.Client{
		Cli:        r.Client,
		Ctx:        ctx,
		KibanaSpec: *targetInstance,
		Req:        req,
	}

	if indexPattern.ObjectMeta.DeletionTimestamp.IsZero() {

		if err := kibanaUtils.DependenciesFulfilled(kibanaClient, indexPattern.Spec.GetSavedObject()); err != nil {
			r.Recorder.Event(&indexPattern, "Warning", "Missing dependencies",
				fmt.Sprintf("Some of declared dependencies are not present yet: %s", err.Error()))
			return utils.GetRequeueResult(), err
		}

		logger.Info("Creating/Updating index pattern", "id", req.Name)
		res, err := kibanaUtils.UpsertSavedObject(kibanaClient, savedObjectType, indexPattern.ObjectMeta, indexPattern.Spec.GetSavedObject())

		if err == nil {
			r.Recorder.Event(&indexPattern, "Normal", "Created",
				fmt.Sprintf("Created/Updated %s/%s %s", indexPattern.APIVersion, indexPattern.Kind, indexPattern.Name))
		} else {
			r.Recorder.Event(&indexPattern, "Warning", "Failed to create/update",
				fmt.Sprintf("Failed to create/update %s/%s %s: %s", indexPattern.APIVersion, indexPattern.Kind, indexPattern.Name, err.Error()))
		}

		if !controllerutil.ContainsFinalizer(&indexPattern, indexPatternFinalizer) {
			controllerutil.AddFinalizer(&indexPattern, indexPatternFinalizer)
			if err := r.Update(ctx, &indexPattern); err != nil {
				return ctrl.Result{}, err
			}
		}

		return res, err
	} else {
		// The object is being deleted
		if controllerutil.ContainsFinalizer(&indexPattern, indexPatternFinalizer) {
			if _, err := kibanaUtils.DeleteSavedObject(kibanaClient, savedObjectType, indexPattern.ObjectMeta, indexPattern.Spec.GetSavedObject()); err != nil {
				return ctrl.Result{}, err
			}

			controllerutil.RemoveFinalizer(&indexPattern, indexPatternFinalizer)
			if err := r.Update(ctx, &indexPattern); err != nil {
				return ctrl.Result{}, err
			}
		}

		return ctrl.Result{}, nil
	}
}

// SetupWithManager sets up the controller with the Manager.
func (r *IndexPatternReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&kibanaeckv1alpha1.IndexPattern{}).
		WithEventFilter(utils.CommonEventFilter()).
		Complete(r)
}

func (r *IndexPatternReconciler) getTargetInstance(object runtime.Object, TargetConfig kibanaeckv1alpha1.CommonKibanaConfig, ctx context.Context, namespace string) (*configv2.KibanaSpec, error) {
	targetInstance := r.ProjectConfig.Kibana
	if TargetConfig.KibanaInstance != "" {
		var resourceInstance kibanaeckv1alpha1.KibanaInstance
		if err := kibanaUtils.GetTargetInstance(r.Client, ctx, namespace, TargetConfig.KibanaInstance, &resourceInstance); err != nil {
			r.Recorder.Event(object, "Warning", "Failed to load target instance", fmt.Sprintf("Target instance not found: %s", err.Error()))
			return nil, err
		}

		targetInstance = resourceInstance.Spec
	}
	return &targetInstance, nil
}
