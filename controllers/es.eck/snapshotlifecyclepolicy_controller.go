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

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	eseckv1alpha1 "github.com/xco-sk/eck-custom-resources/apis/es.eck/v1alpha1"
)

// SnapshotLifecyclePolicyReconciler reconciles a SnapshotLifecyclePolicy object
type SnapshotLifecyclePolicyReconciler struct {
	client.Client
	Scheme        *runtime.Scheme
	ProjectConfig configv2.ProjectConfig
}

//+kubebuilder:rbac:groups=es.eck.github.com,resources=snapshotlifecyclepolicies,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=es.eck.github.com,resources=snapshotlifecyclepolicies/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=es.eck.github.com,resources=snapshotlifecyclepolicies/finalizers,verbs=update

func (r *SnapshotLifecyclePolicyReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// Define esclient as a singleton
	esClient, createClientErr := utils.GetElasticsearchClient(r.Client, ctx, r.ProjectConfig.TargetCluster, req)
	if createClientErr != nil {
		logger.Error(createClientErr, "Failed to create Elasticsearch client")
		return utils.GetRequeueResult(), client.IgnoreNotFound(createClientErr)
	}

	var snapshotLifecyclePolicy eseckv1alpha1.SnapshotLifecyclePolicy
	if err := r.Get(ctx, req.NamespacedName, &snapshotLifecyclePolicy); err != nil {
		logger.Info("Deleting Snapshot lifecycle policy", "snapshot lifecycle policy", req.Name)

		return utils.DeleteSnapshotLifecyclePolicy(esClient, req.Name)
	}

	return utils.UpsertSnapshotLifecyclePolicy(esClient, snapshotLifecyclePolicy)
}

// SetupWithManager sets up the controller with the Manager.
func (r *SnapshotLifecyclePolicyReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&eseckv1alpha1.SnapshotLifecyclePolicy{}).
		Complete(r)
}
