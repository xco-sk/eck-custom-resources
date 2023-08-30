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

package main

import (
	"flag"
	"os"

	// Import all Kubernetes client auth plugins (e.g. Azure, GCP, OIDC, etc.)
	// to ensure that exec-entrypoint and run can make use of them.
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	configv2 "github.com/xco-sk/eck-custom-resources/apis/config/v2"
	eseckv1alpha1 "github.com/xco-sk/eck-custom-resources/apis/es.eck/v1alpha1"

	kibanaeckv1alpha1 "github.com/xco-sk/eck-custom-resources/apis/kibana.eck/v1alpha1"
	eseckcontrollers "github.com/xco-sk/eck-custom-resources/controllers/es.eck"
	kibanaeckcontrollers "github.com/xco-sk/eck-custom-resources/controllers/kibana.eck"
	//+kubebuilder:scaffold:imports
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))

	utilruntime.Must(eseckv1alpha1.AddToScheme(scheme))
	utilruntime.Must(configv2.AddToScheme(scheme))
	utilruntime.Must(kibanaeckv1alpha1.AddToScheme(scheme))
	//+kubebuilder:scaffold:scheme
}

func main() {
	var configFile string
	flag.StringVar(&configFile, "config", "",
		"The controller will load its initial configuration from this file. "+
			"Omit this flag to use the default configuration values. "+
			"Command-line flags override configuration from this file.")
	opts := zap.Options{
		Development: true,
	}
	opts.BindFlags(flag.CommandLine)
	flag.Parse()

	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))

	var err error
	ctrlConfig := configv2.ProjectConfig{}

	options := ctrl.Options{Scheme: scheme}
	if configFile != "" {
		options, err = options.AndFrom(ctrl.ConfigFile().AtPath(configFile).OfKind(&ctrlConfig))
		if err != nil {
			setupLog.Error(err, "unable to load the config file")
			os.Exit(1)
		}
	}

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), options)
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	if err = (&eseckcontrollers.IndexReconciler{
		Client:        mgr.GetClient(),
		Scheme:        mgr.GetScheme(),
		ProjectConfig: ctrlConfig,
		Recorder:      mgr.GetEventRecorderFor("index_controller"),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Index")
		os.Exit(1)
	}
	if err = (&eseckcontrollers.IndexTemplateReconciler{
		Client:        mgr.GetClient(),
		Scheme:        mgr.GetScheme(),
		ProjectConfig: ctrlConfig,
		Recorder:      mgr.GetEventRecorderFor("indextemplate_controller"),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "IndexTemplate")
		os.Exit(1)
	}
	if err = (&eseckcontrollers.IndexLifecyclePolicyReconciler{
		Client:        mgr.GetClient(),
		Scheme:        mgr.GetScheme(),
		ProjectConfig: ctrlConfig,
		Recorder:      mgr.GetEventRecorderFor("indexlifecyclepolicy_controller"),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "IndexLifecyclePolicy")
		os.Exit(1)
	}
	if err = (&eseckcontrollers.SnapshotLifecyclePolicyReconciler{
		Client:        mgr.GetClient(),
		Scheme:        mgr.GetScheme(),
		ProjectConfig: ctrlConfig,
		Recorder:      mgr.GetEventRecorderFor("snapshotlifecyclepolicy_controller"),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "SnapshotLifecyclePolicy")
		os.Exit(1)
	}
	if err = (&eseckcontrollers.IngestPipelineReconciler{
		Client:        mgr.GetClient(),
		Scheme:        mgr.GetScheme(),
		ProjectConfig: ctrlConfig,
		Recorder:      mgr.GetEventRecorderFor("ingestpipeline_controller"),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "IngestPipeline")
		os.Exit(1)
	}
	if err = (&eseckcontrollers.SnapshotRepositoryReconciler{
		Client:        mgr.GetClient(),
		Scheme:        mgr.GetScheme(),
		ProjectConfig: ctrlConfig,
		Recorder:      mgr.GetEventRecorderFor("snapshotrepository_controller"),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "SnapshotRepository")
		os.Exit(1)
	}
	if err = (&kibanaeckcontrollers.SavedSearchReconciler{
		Client:        mgr.GetClient(),
		Scheme:        mgr.GetScheme(),
		ProjectConfig: ctrlConfig,
		Recorder:      mgr.GetEventRecorderFor("savedsearch_controller"),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "SavedSearch")
		os.Exit(1)
	}
	if err = (&kibanaeckcontrollers.IndexPatternReconciler{
		Client:        mgr.GetClient(),
		Scheme:        mgr.GetScheme(),
		ProjectConfig: ctrlConfig,
		Recorder:      mgr.GetEventRecorderFor("indexpattern_controller"),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "IndexPattern")
		os.Exit(1)
	}
	if err = (&kibanaeckcontrollers.VisualizationReconciler{
		Client:        mgr.GetClient(),
		Scheme:        mgr.GetScheme(),
		ProjectConfig: ctrlConfig,
		Recorder:      mgr.GetEventRecorderFor("visualization_controller"),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Visualization")
		os.Exit(1)
	}
	if err = (&kibanaeckcontrollers.DashboardReconciler{
		Client:        mgr.GetClient(),
		Scheme:        mgr.GetScheme(),
		ProjectConfig: ctrlConfig,
		Recorder:      mgr.GetEventRecorderFor("dashboard_controller"),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Dashboard")
		os.Exit(1)
	}
	if err = (&eseckcontrollers.ElasticsearchRoleReconciler{
		Client:        mgr.GetClient(),
		Scheme:        mgr.GetScheme(),
		ProjectConfig: ctrlConfig,
		Recorder:      mgr.GetEventRecorderFor("elasticsearchrole_controller"),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "ElasticsearchRole")
		os.Exit(1)
	}
	if err = (&eseckcontrollers.ElasticsearchUserReconciler{
		Client:        mgr.GetClient(),
		Scheme:        mgr.GetScheme(),
		ProjectConfig: ctrlConfig,
		Recorder:      mgr.GetEventRecorderFor("elasticsearchuser_controller"),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "ElasticsearchUser")
		os.Exit(1)
	}
	if err = (&eseckcontrollers.ElasticsearchApikeyReconciler{
		Client:        mgr.GetClient(),
		Scheme:        mgr.GetScheme(),
		ProjectConfig: ctrlConfig,
		Recorder:      mgr.GetEventRecorderFor("elasticsearchapikey_controller"),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "ElasticsearchApikey")
		os.Exit(1)
	}
	if err = (&kibanaeckcontrollers.SpaceReconciler{
		Client:        mgr.GetClient(),
		Scheme:        mgr.GetScheme(),
		ProjectConfig: ctrlConfig,
		Recorder:      mgr.GetEventRecorderFor("kibanaspace_controller"),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Space")
		os.Exit(1)
	}
	if err = (&kibanaeckcontrollers.LensReconciler{
		Client:        mgr.GetClient(),
		Scheme:        mgr.GetScheme(),
		ProjectConfig: ctrlConfig,
		Recorder:      mgr.GetEventRecorderFor("kibanalens_controller"),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Lens")
		os.Exit(1)
	}
	if err = (&kibanaeckcontrollers.DataViewReconciler{
		Client:        mgr.GetClient(),
		Scheme:        mgr.GetScheme(),
		ProjectConfig: ctrlConfig,
		Recorder:      mgr.GetEventRecorderFor("kibanadataview_controller"),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "DataView")
		os.Exit(1)
	}
	//+kubebuilder:scaffold:builder

	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up health check")
		os.Exit(1)
	}
	if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up ready check")
		os.Exit(1)
	}

	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}
