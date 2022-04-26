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
	"github.com/xco-sk/eck-custom-resources/utils"
	esutils "github.com/xco-sk/eck-custom-resources/utils/elasticsearch"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	eseckv1alpha1 "github.com/xco-sk/eck-custom-resources/apis/es.eck/v1alpha1"
)

// SnapshotRepositoryReconciler reconciles a SnapshotRepository object
type SnapshotRepositoryReconciler struct {
	client.Client
	Scheme        *runtime.Scheme
	ProjectConfig configv2.ProjectConfig
}

//+kubebuilder:rbac:groups=es.eck.github.com,resources=snapshotrepositories,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=es.eck.github.com,resources=snapshotrepositories/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=es.eck.github.com,resources=snapshotrepositories/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the SnapshotRepository object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.0/pkg/reconcile
func (r *SnapshotRepositoryReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// Define esclient as a singleton
	esClient, createClientErr := esutils.GetElasticsearchClient(r.Client, ctx, r.ProjectConfig.TargetCluster, req)
	if createClientErr != nil {
		logger.Error(createClientErr, "Failed to create Elasticsearch client")
		return utils.GetRequeueResult(), client.IgnoreNotFound(createClientErr)
	}

	var snapshotRepository eseckv1alpha1.SnapshotRepository
	if err := r.Get(ctx, req.NamespacedName, &snapshotRepository); err != nil {
		logger.Info("Deleting Snapshot repository", "snapshot repository", req.Name)
		return esutils.DeleteSnapshotRepository(esClient, req.Name)
	}

	logger.Info("Creating/Updating Snapshot repository", "snapshot repository", req.Name)
	return esutils.UpsertSnapshotRepository(esClient, snapshotRepository)
}

// SetupWithManager sets up the controller with the Manager.
func (r *SnapshotRepositoryReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&eseckv1alpha1.SnapshotRepository{}).
		Complete(r)
}
