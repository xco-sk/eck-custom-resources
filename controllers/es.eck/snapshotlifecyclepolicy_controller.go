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
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	eseckv1alpha1 "github.com/xco-sk/eck-custom-resources/apis/es.eck/v1alpha1"
)

// SnapshotLifecyclePolicyReconciler reconciles a SnapshotLifecyclePolicy object
type SnapshotLifecyclePolicyReconciler struct {
	client.Client
	Scheme        *runtime.Scheme
	ProjectConfig configv2.ProjectConfig
	Recorder      record.EventRecorder
}

//+kubebuilder:rbac:groups=es.eck.github.com,resources=snapshotlifecyclepolicies,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=es.eck.github.com,resources=snapshotlifecyclepolicies/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=es.eck.github.com,resources=snapshotlifecyclepolicies/finalizers,verbs=update

func (r *SnapshotLifecyclePolicyReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	finalizer := "snapshotlifecyclepolicies.es.eck.github.com/finalizer"

	var snapshotLifecyclePolicy eseckv1alpha1.SnapshotLifecyclePolicy
	if err := r.Get(ctx, req.NamespacedName, &snapshotLifecyclePolicy); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	targetInstance, err := r.getTargetInstance(&snapshotLifecyclePolicy, snapshotLifecyclePolicy.Spec.TargetConfig, ctx, req.Namespace)
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

	if snapshotLifecyclePolicy.ObjectMeta.DeletionTimestamp.IsZero() {
		res, err := esutils.UpsertSnapshotLifecyclePolicy(esClient, snapshotLifecyclePolicy)

		if err == nil {
			r.Recorder.Event(&snapshotLifecyclePolicy, "Normal", "Created",
				fmt.Sprintf("Created/Updated %s/%s %s", snapshotLifecyclePolicy.APIVersion, snapshotLifecyclePolicy.Kind, snapshotLifecyclePolicy.Name))
		} else {
			r.Recorder.Event(&snapshotLifecyclePolicy, "Warning", "Failed to create/update",
				fmt.Sprintf("Failed to create/update %s/%s %s: %s", snapshotLifecyclePolicy.APIVersion, snapshotLifecyclePolicy.Kind, snapshotLifecyclePolicy.Name, err.Error()))
		}

		if err := r.addFinalizer(&snapshotLifecyclePolicy, finalizer, ctx); err != nil {
			return ctrl.Result{}, err
		}
		return res, err
	} else {
		// The object is being deleted
		if controllerutil.ContainsFinalizer(&snapshotLifecyclePolicy, finalizer) {
			logger.Info("Deleting object", "snapshotLifecyclePolicy", snapshotLifecyclePolicy.Name)
			if _, err := esutils.DeleteSnapshotLifecyclePolicy(esClient, req.Name); err != nil {
				return ctrl.Result{}, err
			}

			controllerutil.RemoveFinalizer(&snapshotLifecyclePolicy, finalizer)
			if err := r.Update(ctx, &snapshotLifecyclePolicy); err != nil {
				return ctrl.Result{}, err
			}
		}

		return ctrl.Result{}, nil
	}
}

// SetupWithManager sets up the controller with the Manager.
func (r *SnapshotLifecyclePolicyReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&eseckv1alpha1.SnapshotLifecyclePolicy{}).
		WithEventFilter(utils.CommonEventFilter()).
		Complete(r)
}

func (r *SnapshotLifecyclePolicyReconciler) addFinalizer(o client.Object, finalizer string, ctx context.Context) error {
	if !controllerutil.ContainsFinalizer(o, finalizer) {
		controllerutil.AddFinalizer(o, finalizer)
		if err := r.Update(ctx, o); err != nil {
			return err
		}
	}
	return nil
}

func (r *SnapshotLifecyclePolicyReconciler) getTargetInstance(object runtime.Object, TargetConfig eseckv1alpha1.CommonElasticsearchConfig, ctx context.Context, namespace string) (*configv2.ElasticsearchSpec, error) {
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
