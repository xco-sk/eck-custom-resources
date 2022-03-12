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

package controllers

import (
	"context"
	eseckv1 "github.com/xco-sk/eck-custom-resources/api/v1alpha1"
	esCliUtils "github.com/xco-sk/eck-custom-resources/utils"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// IndexReconciler reconciles a Index object
type IndexReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=core,resources=secrets,verbs=get
//+kubebuilder:rbac:groups=es.eck.github.com,resources=indices,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=es.eck.github.com,resources=indices/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=es.eck.github.com,resources=indices/finalizers,verbs=update

func (r *IndexReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	esClient, _ := esCliUtils.GetElasticsearchClient(r.Client, ctx, index.Spec.TargetCluster, req)
	// TODO remove target cluster - it has to be defined on operator level, not resource level.

	var index eseckv1.Index
	if err := r.Get(ctx, req.NamespacedName, &index); err != nil {
		logger.Info("Index does not exists - deleting", "index", req.Name)

		return ctrl.Result{}, nil
	}

	var indicesCreateResponse, err = esClient.Indices.Create("test",
		esClient.Indices.Create.WithHuman(),
	)
	if err != nil {
		logger.Error(err, "Error querying")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	defer indicesCreateResponse.Body.Close()
	logger.Info("Got", "response", indicesCreateResponse)

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *IndexReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&eseckv1.Index{}).
		Complete(r)
}
