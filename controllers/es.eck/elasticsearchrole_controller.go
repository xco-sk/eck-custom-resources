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

	configv2 "github.com/xco-sk/eck-custom-resources/apis/config/v2"
	"github.com/xco-sk/eck-custom-resources/utils"
	esutils "github.com/xco-sk/eck-custom-resources/utils/elasticsearch"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	eseckv1alpha1 "github.com/xco-sk/eck-custom-resources/apis/es.eck/v1alpha1"
)

// ElasticsearchRoleReconciler reconciles a ElasticsearchRole object
type ElasticsearchRoleReconciler struct {
	client.Client
	Scheme        *runtime.Scheme
	ProjectConfig configv2.ProjectConfig
	Recorder      record.EventRecorder
}

//+kubebuilder:rbac:groups=es.eck.github.com,resources=elasticsearchroles,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=es.eck.github.com,resources=elasticsearchroles/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=es.eck.github.com,resources=elasticsearchroles/finalizers,verbs=update

func (r *ElasticsearchRoleReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	finalizer := "elasticsearchroles.es.eck.github.com/finalizer"

	var role eseckv1alpha1.ElasticsearchRole
	if err := r.Get(ctx, req.NamespacedName, &role); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	targetInstance, err := r.getTargetInstance(&role, role.Spec.TargetConfig, ctx, req.Namespace)
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

	if role.ObjectMeta.DeletionTimestamp.IsZero() {
		logger.Info("Creating/Updating Role", "role", req.Name)
		res, err := esutils.UpsertRole(esClient, role)

		if err == nil {
			r.Recorder.Event(&role, "Normal", "Created",
				fmt.Sprintf("Created/Updated %s/%s %s", role.APIVersion, role.Kind, role.Name))
		} else {
			r.Recorder.Event(&role, "Warning", "Failed to create/update",
				fmt.Sprintf("Failed to create/update %s/%s %s: %s", role.APIVersion, role.Kind, role.Name, err.Error()))
		}

		if !controllerutil.ContainsFinalizer(&role, finalizer) {
			controllerutil.AddFinalizer(&role, finalizer)
			if err := r.Update(ctx, &role); err != nil {
				return ctrl.Result{}, err
			}
		}
		return res, err
	} else {
		// The object is being deleted
		if controllerutil.ContainsFinalizer(&role, finalizer) {
			logger.Info("Deleting object", "role", role.Name)
			if _, err := esutils.DeleteRole(esClient, req.Name); err != nil {
				return ctrl.Result{}, err
			}

			controllerutil.RemoveFinalizer(&role, finalizer)
			if err := r.Update(ctx, &role); err != nil {
				return ctrl.Result{}, err
			}
		}

		return ctrl.Result{}, nil
	}
}

func (r *ElasticsearchRoleReconciler) getTargetInstance(object runtime.Object, TargetConfig eseckv1alpha1.CommonElasticsearchConfig, ctx context.Context, namespace string) (*configv2.ElasticsearchSpec, error) {
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

// SetupWithManager sets up the controller with the Manager.
func (r *ElasticsearchRoleReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&eseckv1alpha1.ElasticsearchRole{}).
		WithEventFilter(utils.CommonEventFilter()).
		Complete(r)
}
