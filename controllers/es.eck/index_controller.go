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
	configv2 "github.com/xco-sk/eck-custom-resources/apis/config/v2"
	eseckv1alpha1 "github.com/xco-sk/eck-custom-resources/apis/es.eck/v1alpha1"
	"github.com/xco-sk/eck-custom-resources/utils"
	esutils "github.com/xco-sk/eck-custom-resources/utils/elasticsearch"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"strings"
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

	esClient, createClientErr := esutils.GetElasticsearchClient(r.Client, ctx, r.ProjectConfig.Elasticsearch, req)
	if createClientErr != nil {
		logger.Error(createClientErr, "Failed to create Elasticsearch client")
		return utils.GetRequeueResult(), client.IgnoreNotFound(createClientErr)
	}

	var index eseckv1alpha1.Index
	if err := r.Get(ctx, req.NamespacedName, &index); err != nil {
		logger.Info("Deleting index", "index", req.Name)

		return esutils.DeleteIndexIfEmpty(esClient, req.Name)
	}

	indexExists, indexExistsErr := esutils.VerifyIndexExists(esClient, req.Name)
	if indexExistsErr != nil {
		logger.Error(indexExistsErr, "Failed to verify if index exists")
		return ctrl.Result{}, indexExistsErr
	}

	if indexExists {
		result, err, recreated := esutils.RecreateIndexIfEmpty(esClient, req)
		if recreated {
			return result, err
		} else {
			// TODO patch index
		}
	}

	var indicesCreateResponse, createIndexErr = esClient.Indices.Create(index.Name,
		esClient.Indices.Create.WithBody(strings.NewReader(index.Spec.Body)),
	)

	if createIndexErr != nil || indicesCreateResponse.IsError() {
		logger.Error(createIndexErr, "Error creating index")
		return utils.GetRequeueResult(), client.IgnoreNotFound(createIndexErr)
	}
	defer indicesCreateResponse.Body.Close()

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *IndexReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&eseckv1alpha1.Index{}).
		Complete(r)
}
