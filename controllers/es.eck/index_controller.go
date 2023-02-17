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

	"github.com/elastic/go-elasticsearch/v8"
	configv2 "github.com/xco-sk/eck-custom-resources/apis/config/v2"
	eseckv1alpha1 "github.com/xco-sk/eck-custom-resources/apis/es.eck/v1alpha1"
	"github.com/xco-sk/eck-custom-resources/utils"
	esutils "github.com/xco-sk/eck-custom-resources/utils/elasticsearch"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// IndexReconciler reconciles a Index object
type IndexReconciler struct {
	client.Client
	Scheme        *runtime.Scheme
	ProjectConfig configv2.ProjectConfig
	Recorder      record.EventRecorder
}

//+kubebuilder:rbac:groups=core,resources=secrets,verbs=get
//+kubebuilder:rbac:groups=es.eck.github.com,resources=indices,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=es.eck.github.com,resources=indices/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=es.eck.github.com,resources=indices/finalizers,verbs=update

func (r *IndexReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	finalizer := "indices.es.eck.github.com/finalizer"

	var index eseckv1alpha1.Index
	if err := r.Get(ctx, req.NamespacedName, &index); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	targetInstance, err := r.getTargetInstance(&index, index.Spec.TargetConfig, ctx, req.Namespace)
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

	if index.ObjectMeta.DeletionTimestamp.IsZero() {
		res, err := r.createUpdate(ctx, req, esClient, index)

		if err := r.addFinalizer(&index, finalizer, ctx); err != nil {
			return ctrl.Result{}, err
		}

		return utils.RecordEventAndReturn(res, err, r.Recorder, utils.Event{
			Object:  &index,
			Name:    req.Name,
			Reason:  "Create/Update",
			Message: "",
		})
	} else {
		// The object is being deleted
		if controllerutil.ContainsFinalizer(&index, finalizer) {
			logger.Info("Deleting object", "index", index.Name)
			if _, err := esutils.DeleteIndexIfEmpty(esClient, req.Name); err != nil {
				return ctrl.Result{}, err
			}

			controllerutil.RemoveFinalizer(&index, finalizer)
			if err := r.Update(ctx, &index); err != nil {
				return ctrl.Result{}, err
			}
		}

		return ctrl.Result{}, nil
	}
}

func (r *IndexReconciler) createUpdate(ctx context.Context, req ctrl.Request, esClient *elasticsearch.Client, index eseckv1alpha1.Index) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	if err := esutils.DependenciesFulfilled(esClient, index.Spec.Dependencies); err != nil {
		r.Recorder.Event(&index, "Warning", "Missing dependencies",
			fmt.Sprintf("Some of declared dependencies are not present yet: %s", err.Error()))
		return utils.GetRequeueResult(), err
	}

	indexExists, indexExistsErr := esutils.VerifyIndexExists(esClient, req.Name)
	if indexExistsErr != nil {
		logger.Error(indexExistsErr, "Failed to verify if index exists")
		return ctrl.Result{}, indexExistsErr
	}

	if indexExists {
		isEmpty, indexEmptyErr := esutils.VerifyIndexEmpty(esClient, req.Name)
		if indexEmptyErr != nil {
			logger.Error(indexExistsErr, "Failed to verify if index is empty")
			return utils.GetRequeueResult(), client.IgnoreNotFound(indexEmptyErr)
		}

		if isEmpty {
			_, deleteErr := esutils.DeleteIndex(esClient, req.Name)
			if deleteErr != nil {
				logger.Error(deleteErr, "Failed to delete index")
				return utils.GetRequeueResult(), client.IgnoreNotFound(deleteErr)
			}

			return esutils.CreateIndex(esClient, index)
		}
		return esutils.UpdateIndex(esClient, index, r.Recorder)
	}
	return esutils.CreateIndex(esClient, index)
}

// SetupWithManager sets up the controller with the Manager.
func (r *IndexReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&eseckv1alpha1.Index{}).
		WithEventFilter(utils.CommonEventFilter()).
		Complete(r)
}

func (r *IndexReconciler) addFinalizer(o client.Object, finalizer string, ctx context.Context) error {
	if !controllerutil.ContainsFinalizer(o, finalizer) {
		controllerutil.AddFinalizer(o, finalizer)
		if err := r.Update(ctx, o); err != nil {
			return err
		}
	}
	return nil
}

func (r *IndexReconciler) getTargetInstance(object runtime.Object, TargetConfig eseckv1alpha1.CommonElasticsearchConfig, ctx context.Context, namespace string) (*configv2.ElasticsearchSpec, error) {
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
