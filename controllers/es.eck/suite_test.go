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
	"path/filepath"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	configv2 "github.com/xco-sk/eck-custom-resources/apis/config/v2"
	eseckv1alpha1 "github.com/xco-sk/eck-custom-resources/apis/es.eck/v1alpha1"
	//+kubebuilder:scaffold:imports
)

// These tests use Ginkgo (BDD-style Go testing framework). Refer to
// http://onsi.github.io/ginkgo/ to learn more about Ginkgo.

var cfg *rest.Config
var k8sClient client.Client
var k8sManager ctrl.Manager
var testEnv *envtest.Environment

func TestAPIs(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecsWithDefaultAndCustomReporters(t,
		"Controller Suite",
		[]Reporter{})
}

var _ = BeforeSuite(func() {
	logf.SetLogger(zap.New(zap.WriteTo(GinkgoWriter), zap.UseDevMode(true)))

	By("bootstrapping test environment")
	testEnv = &envtest.Environment{
		CRDDirectoryPaths:     []string{filepath.Join("..", "..", "config", "crd", "bases")},
		ErrorIfCRDPathMissing: true,
	}

	var err error
	cfg, err = testEnv.Start()
	Expect(err).NotTo(HaveOccurred())
	Expect(cfg).NotTo(BeNil())

	err = eseckv1alpha1.AddToScheme(scheme.Scheme)
	Expect(err).NotTo(HaveOccurred())

	//+kubebuilder:scaffold:scheme

	k8sManager, err = ctrl.NewManager(cfg, ctrl.Options{
		Scheme: scheme.Scheme,
	})
	Expect(err).ToNot(HaveOccurred())

	k8sClient, err = client.New(cfg, client.Options{Scheme: scheme.Scheme})
	Expect(err).NotTo(HaveOccurred())
	Expect(k8sClient).NotTo(BeNil())

	pc := configv2.ProjectConfig{}
	mgrClient := k8sManager.GetClient()
	mgrScheme := k8sManager.GetScheme()

	err = (&ComponentTemplateReconciler{
		Client: mgrClient, Scheme: mgrScheme, ProjectConfig: pc,
		Recorder: k8sManager.GetEventRecorderFor("component-template"),
	}).SetupWithManager(k8sManager)
	Expect(err).ToNot(HaveOccurred())

	err = (&ElasticsearchApikeyReconciler{
		Client: mgrClient, Scheme: mgrScheme, ProjectConfig: pc,
		Recorder: k8sManager.GetEventRecorderFor("elasticsearch-apikey"),
	}).SetupWithManager(k8sManager)
	Expect(err).ToNot(HaveOccurred())

	err = (&ElasticsearchRoleReconciler{
		Client: mgrClient, Scheme: mgrScheme, ProjectConfig: pc,
		Recorder: k8sManager.GetEventRecorderFor("elasticsearch-role"),
	}).SetupWithManager(k8sManager)
	Expect(err).ToNot(HaveOccurred())

	err = (&ElasticsearchUserReconciler{
		Client: mgrClient, Scheme: mgrScheme, ProjectConfig: pc,
		Recorder: k8sManager.GetEventRecorderFor("elasticsearch-user"),
	}).SetupWithManager(k8sManager)
	Expect(err).ToNot(HaveOccurred())

	err = (&IndexReconciler{
		Client: mgrClient, Scheme: mgrScheme, ProjectConfig: pc,
		Recorder: k8sManager.GetEventRecorderFor("index"),
	}).SetupWithManager(k8sManager)
	Expect(err).ToNot(HaveOccurred())

	err = (&IndexLifecyclePolicyReconciler{
		Client: mgrClient, Scheme: mgrScheme, ProjectConfig: pc,
		Recorder: k8sManager.GetEventRecorderFor("index-lifecycle-policy"),
	}).SetupWithManager(k8sManager)
	Expect(err).ToNot(HaveOccurred())

	err = (&IndexTemplateReconciler{
		Client: mgrClient, Scheme: mgrScheme, ProjectConfig: pc,
		Recorder: k8sManager.GetEventRecorderFor("index-template"),
	}).SetupWithManager(k8sManager)
	Expect(err).ToNot(HaveOccurred())

	err = (&IngestPipelineReconciler{
		Client: mgrClient, Scheme: mgrScheme, ProjectConfig: pc,
		Recorder: k8sManager.GetEventRecorderFor("ingest-pipeline"),
	}).SetupWithManager(k8sManager)
	Expect(err).ToNot(HaveOccurred())

	err = (&SnapshotLifecyclePolicyReconciler{
		Client: mgrClient, Scheme: mgrScheme, ProjectConfig: pc,
		Recorder: k8sManager.GetEventRecorderFor("snapshot-lifecycle-policy"),
	}).SetupWithManager(k8sManager)
	Expect(err).ToNot(HaveOccurred())

	err = (&SnapshotRepositoryReconciler{
		Client: mgrClient, Scheme: mgrScheme, ProjectConfig: pc,
		Recorder: k8sManager.GetEventRecorderFor("snapshot-repository"),
	}).SetupWithManager(k8sManager)
	Expect(err).ToNot(HaveOccurred())

}, 60)

var _ = AfterSuite(func() {
	By("tearing down test environment")
	err := testEnv.Stop()
	Expect(err).ToNot(HaveOccurred())
})
