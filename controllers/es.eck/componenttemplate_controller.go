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

package eseck

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	configv2 "github.com/xco-sk/eck-custom-resources/apis/config/v2"
	eseckv1alpha1 "github.com/xco-sk/eck-custom-resources/apis/es.eck/v1alpha1"
	"github.com/xco-sk/eck-custom-resources/utils"
	esutils "github.com/xco-sk/eck-custom-resources/utils/elasticsearch"
)

// ComponentTemplateReconciler reconciles a ComponentTemplate object
type ComponentTemplateReconciler struct {
	client.Client
	Scheme        *runtime.Scheme
	ProjectConfig configv2.ProjectConfig
	Recorder      record.EventRecorder
}

//+kubebuilder:rbac:groups=es.eck.github.com,resources=componenttemplates,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=es.eck.github.com,resources=componenttemplates/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=es.eck.github.com,resources=componenttemplates/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the ComponentTemplate object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.4/pkg/reconcile
func (r *ComponentTemplateReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	finalizer := "componenttemplate.es.eck.github.com/finalizer"
	var comTem eseckv1alpha1.ComponentTemplate

	if err := r.Get(ctx, req.NamespacedName, &comTem); err != err {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	targetInstance, err := r.getTargetInstance(&comTem, comTem.Spec.TargetConfig, ctx, req.Namespace)
	if err != nil {
		return utils.GetRequeueResult(), err
	}
	if !targetInstance.Enabled {
		logger.Info("Elasticsearch reconciler disabled, not reconciling.", "Resource", req.NamespacedName)
		return ctrl.Result{}, nil
	}

	esClient, createClientErr := esutils.GetElasticsearchClient(r.Client, ctx, *targetInstance, req)
	if createClientErr != nil {
		logger.Error(createClientErr, "Failed to create Elasticsearch client")
		return utils.GetRequeueResult(), client.IgnoreNotFound(createClientErr)
	}
	if comTem.ObjectMeta.DeletionTimestamp.IsZero() {
		logger.Info("Creating/Updating component template", "componentTemplate", req.Name)
		res, err := esutils.UpsertComponentTemplate(esClient, comTem)
		if err == nil {
			r.Recorder.Event(&comTem, "Normal", "Created",
				fmt.Sprintf("Created/Updated %s/%s %s", comTem.APIVersion, comTem.Kind, comTem.Name))
		} else {
			r.Recorder.Event(&comTem, "Warning", "Failed to create/update",
				fmt.Sprintf("Failed to create/update %s/%s %s: %s", comTem.APIVersion, comTem.Kind, comTem.Name, err.Error()))
		}

		if err := r.addFinalizer(&comTem, finalizer, ctx); err != nil {
			return ctrl.Result{}, err
		}
		return res, err
	} else {
		if controllerutil.ContainsFinalizer(&comTem, finalizer) {
			logger.Info("Deleting object", "componentTemplate", comTem.Name)
			if _, err := esutils.DeleteComponentTemplate(esClient, comTem.Name); err != nil {
				return ctrl.Result{}, err
			}

			controllerutil.RemoveFinalizer(&comTem, finalizer)
			if err := r.Update(ctx, &comTem); err != nil {
				return ctrl.Result{}, err
			}
		}
	}

	return ctrl.Result{}, nil
}

func (r *ComponentTemplateReconciler) getTargetInstance(object runtime.Object, TargetConfig eseckv1alpha1.CommonElasticsearchConfig, ctx context.Context, namespace string) (*configv2.ElasticsearchSpec, error) {
	targetInstance := r.ProjectConfig.Elasticsearch
	if TargetConfig.ElasticsearchInstance != "" {
		var resourceInstance eseckv1alpha1.ElasticsearchInstance
		if err := esutils.GetTargetElasticsearchInstance(r.Client, ctx, namespace, TargetConfig.ElasticsearchInstance, &resourceInstance); err != nil {
			r.Recorder.Event(object, "Warning", "Failed to load target instance", fmt.Sprintf("Target instance not found: %s", err.Error()))
			return nil, err
		}

		targetInstance = resourceInstance.Spec
	}
	return &targetInstance, nil
}

func (r *ComponentTemplateReconciler) addFinalizer(o client.Object, finalizer string, ctx context.Context) error {
	if !controllerutil.ContainsFinalizer(o, finalizer) {
		controllerutil.AddFinalizer(o, finalizer)
		if err := r.Update(ctx, o); err != nil {
			return err
		}
	}
	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ComponentTemplateReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&eseckv1alpha1.ComponentTemplate{}).
		Complete(r)
}
