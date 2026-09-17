package main

import (
	"testing"
	"time"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/yaml"

	configv2 "github.com/xco-sk/eck-custom-resources/apis/config/v2"
)

// Every key the removed ControllerManagerConfigurationSpec accepted.
// cacheSyncTimeout is integer nanoseconds because upstream typed it as a bare
// time.Duration; the metav1.Duration fields take strings.
const fullConfig = `
type:
  apiVersion: controller-runtime.sigs.k8s.io/v1alpha1
  kind: ControllerManagerConfig
manager:
  syncPeriod: 10m
  cacheNamespace: some-namespace
  gracefulShutDown: 25s
  leaderElection:
    leaderElect: true
    leaseDuration: 20s
    renewDeadline: 15s
    retryPeriod: 3s
    resourceLock: leases
    resourceName: my-lock
    resourceNamespace: lock-ns
  controller:
    groupKindConcurrency:
      Index.es.eck.github.com: 7
    cacheSyncTimeout: 30000000000
    recoverPanic: true
  metrics:
    bindAddress: :9090
  health:
    healthProbeBindAddress: :9091
    readinessEndpointName: /ready
    livenessEndpointName: /alive
  webhook:
    port: 9444
    host: 0.0.0.0
    certDir: /tmp/certs
`

func TestApplyManagerConfigFullSurface(t *testing.T) {
	var c configv2.ProjectConfig
	if err := yaml.UnmarshalStrict([]byte(fullConfig), &c); err != nil {
		t.Fatalf("strict decode of full config failed: %v", err)
	}
	o := applyManagerConfig(ctrl.Options{}, c.Manager)

	if o.Cache.SyncPeriod == nil || *o.Cache.SyncPeriod != 10*time.Minute {
		t.Errorf("syncPeriod: got %v", o.Cache.SyncPeriod)
	}
	if _, ok := o.Cache.DefaultNamespaces["some-namespace"]; !ok {
		t.Errorf("cacheNamespace: got %v", o.Cache.DefaultNamespaces)
	}
	if o.GracefulShutdownTimeout == nil || *o.GracefulShutdownTimeout != 25*time.Second {
		t.Errorf("gracefulShutDown: got %v (upstream AndFrom never applied this)", o.GracefulShutdownTimeout)
	}
	if !o.LeaderElection {
		t.Error("leaderElect")
	}
	if o.LeaseDuration == nil || *o.LeaseDuration != 20*time.Second {
		t.Errorf("leaseDuration: got %v", o.LeaseDuration)
	}
	if o.RenewDeadline == nil || *o.RenewDeadline != 15*time.Second {
		t.Errorf("renewDeadline: got %v", o.RenewDeadline)
	}
	if o.RetryPeriod == nil || *o.RetryPeriod != 3*time.Second {
		t.Errorf("retryPeriod: got %v", o.RetryPeriod)
	}
	if o.LeaderElectionResourceLock != "leases" {
		t.Errorf("resourceLock: got %q", o.LeaderElectionResourceLock)
	}
	if o.LeaderElectionID != "my-lock" {
		t.Errorf("resourceName: got %q", o.LeaderElectionID)
	}
	if o.LeaderElectionNamespace != "lock-ns" {
		t.Errorf("resourceNamespace: got %q", o.LeaderElectionNamespace)
	}
	if o.Controller.GroupKindConcurrency["Index.es.eck.github.com"] != 7 {
		t.Errorf("groupKindConcurrency: got %v", o.Controller.GroupKindConcurrency)
	}
	if o.Controller.CacheSyncTimeout != 30*time.Second {
		t.Errorf("cacheSyncTimeout: got %v", o.Controller.CacheSyncTimeout)
	}
	if o.Controller.RecoverPanic == nil || !*o.Controller.RecoverPanic {
		t.Errorf("recoverPanic: got %v (upstream AndFrom never applied this)", o.Controller.RecoverPanic)
	}
	if o.Metrics.BindAddress != ":9090" {
		t.Errorf("metrics.bindAddress: got %q", o.Metrics.BindAddress)
	}
	if o.HealthProbeBindAddress != ":9091" {
		t.Errorf("healthProbeBindAddress: got %q", o.HealthProbeBindAddress)
	}
	if o.ReadinessEndpointName != "/ready" {
		t.Errorf("readinessEndpointName: got %q", o.ReadinessEndpointName)
	}
	if o.LivenessEndpointName != "/alive" {
		t.Errorf("livenessEndpointName: got %q", o.LivenessEndpointName)
	}
	if o.WebhookServer == nil {
		t.Error("webhookServer not constructed")
	}
}

// AndFrom let an already-set Option (i.e. a flag) win over the config file.
func TestFlagPrecedencePreserved(t *testing.T) {
	var c configv2.ProjectConfig
	if err := yaml.UnmarshalStrict([]byte(fullConfig), &c); err != nil {
		t.Fatal(err)
	}
	preset := ctrl.Options{HealthProbeBindAddress: ":7777", LeaderElectionID: "from-flag"}
	o := applyManagerConfig(preset, c.Manager)
	if o.HealthProbeBindAddress != ":7777" {
		t.Errorf("config file overrode a preset value: %q", o.HealthProbeBindAddress)
	}
	if o.LeaderElectionID != "from-flag" {
		t.Errorf("config file overrode a preset value: %q", o.LeaderElectionID)
	}
}

func TestEmptyManagerSpecIsHarmless(t *testing.T) {
	o := applyManagerConfig(ctrl.Options{}, configv2.ManagerSpec{})
	if o.LeaderElection {
		t.Error("leader election enabled from an empty spec")
	}
	if o.WebhookServer == nil {
		t.Error("webhook server should still be constructed, as AndFrom did")
	}
}
