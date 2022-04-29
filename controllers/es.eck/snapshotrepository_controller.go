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
	Recorder      record.EventRecorder
}

//+kubebuilder:rbac:groups=es.eck.github.com,resources=snapshotrepositories,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=es.eck.github.com,resources=snapshotrepositories/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=es.eck.github.com,resources=snapshotrepositories/finalizers,verbs=update

func (r *SnapshotRepositoryReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	esClient, createClientErr := esutils.GetElasticsearchClient(r.Client, ctx, r.ProjectConfig.Elasticsearch, req)
	if createClientErr != nil {
		logger.Error(createClientErr, "Failed to create Elasticsearch client")
		return utils.GetRequeueResult(), client.IgnoreNotFound(createClientErr)
	}

	var snapshotRepository eseckv1alpha1.SnapshotRepository
	if err := r.Get(ctx, req.NamespacedName, &snapshotRepository); err != nil {
		logger.Info("Deleting Snapshot repository", "snapshot repository", req.Name)
		res, err := esutils.DeleteSnapshotRepository(esClient, req.Name)
		if err == nil {
			r.Recorder.Event(&snapshotRepository, "Normal", "Failed to delete",
				fmt.Sprintf("Failed to delete %s/%s %s", snapshotRepository.APIVersion, snapshotRepository.Kind, snapshotRepository.Name))
		}
		return res, err
	}

	logger.Info("Creating/Updating Snapshot repository", "snapshot repository", req.Name)
	res, err := esutils.UpsertSnapshotRepository(esClient, snapshotRepository)

	if err == nil {
		r.Recorder.Event(&snapshotRepository, "Normal", "Created",
			fmt.Sprintf("Created/Updated %s/%s %s", snapshotRepository.APIVersion, snapshotRepository.Kind, snapshotRepository.Name))
	} else {
		r.Recorder.Event(&snapshotRepository, "Warning", "Failed to create/update",
			fmt.Sprintf("Failed to create/update %s/%s %s: %s", snapshotRepository.APIVersion, snapshotRepository.Kind, snapshotRepository.Name, err.Error()))
	}

	return res, err
}

// SetupWithManager sets up the controller with the Manager.
func (r *SnapshotRepositoryReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&eseckv1alpha1.SnapshotRepository{}).
		Complete(r)
}
