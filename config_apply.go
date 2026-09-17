package main

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/webhook"

	configv2 "github.com/xco-sk/eck-custom-resources/apis/config/v2"
)

// applyManagerConfig maps the config file's manager section onto ctrl.Options.
//
// It replaces Options.AndFrom, removed with ComponentConfig in controller-runtime
// v0.18.0, and preserves that function's semantics: a value already set on Options
// (i.e. from a command-line flag) wins over the config file.
//
// Two fields are applied here that upstream AndFrom accepted but never used:
// gracefulShutDown and controller.recoverPanic.
func applyManagerConfig(o ctrl.Options, m configv2.ManagerSpec) ctrl.Options {
	o = applyLeaderElectionConfig(o, m.LeaderElection)

	if o.Cache.SyncPeriod == nil && m.SyncPeriod != nil {
		o.Cache.SyncPeriod = &m.SyncPeriod.Duration
	}

	if len(o.Cache.DefaultNamespaces) == 0 && m.CacheNamespace != "" {
		o.Cache.DefaultNamespaces = map[string]cache.Config{m.CacheNamespace: {}}
	}

	if o.GracefulShutdownTimeout == nil && m.GracefulShutdownTimeout != nil {
		o.GracefulShutdownTimeout = &m.GracefulShutdownTimeout.Duration
	}

	if o.Metrics.BindAddress == "" && m.Metrics.BindAddress != "" {
		o.Metrics.BindAddress = m.Metrics.BindAddress
	}

	if o.HealthProbeBindAddress == "" && m.Health.HealthProbeBindAddress != "" {
		o.HealthProbeBindAddress = m.Health.HealthProbeBindAddress
	}

	if o.ReadinessEndpointName == "" && m.Health.ReadinessEndpointName != "" {
		o.ReadinessEndpointName = m.Health.ReadinessEndpointName
	}

	if o.LivenessEndpointName == "" && m.Health.LivenessEndpointName != "" {
		o.LivenessEndpointName = m.Health.LivenessEndpointName
	}

	if o.WebhookServer == nil {
		port := 0
		if m.Webhook.Port != nil {
			port = *m.Webhook.Port
		}
		o.WebhookServer = webhook.NewServer(webhook.Options{
			Port:    port,
			Host:    m.Webhook.Host,
			CertDir: m.Webhook.CertDir,
		})
	}

	if m.Controller != nil {
		if o.Controller.CacheSyncTimeout == 0 && m.Controller.CacheSyncTimeout != nil {
			o.Controller.CacheSyncTimeout = *m.Controller.CacheSyncTimeout
		}

		if len(o.Controller.GroupKindConcurrency) == 0 && len(m.Controller.GroupKindConcurrency) > 0 {
			o.Controller.GroupKindConcurrency = m.Controller.GroupKindConcurrency
		}

		if o.Controller.RecoverPanic == nil && m.Controller.RecoverPanic != nil {
			o.Controller.RecoverPanic = m.Controller.RecoverPanic
		}
	}

	return o
}

func applyLeaderElectionConfig(o ctrl.Options, le *configv2.ManagerLeaderElectionSpec) ctrl.Options {
	if le == nil {
		return o
	}

	if !o.LeaderElection && le.LeaderElect != nil {
		o.LeaderElection = *le.LeaderElect
	}

	if o.LeaderElectionResourceLock == "" && le.ResourceLock != "" {
		o.LeaderElectionResourceLock = le.ResourceLock
	}

	if o.LeaderElectionNamespace == "" && le.ResourceNamespace != "" {
		o.LeaderElectionNamespace = le.ResourceNamespace
	}

	if o.LeaderElectionID == "" && le.ResourceName != "" {
		o.LeaderElectionID = le.ResourceName
	}

	if o.LeaseDuration == nil && le.LeaseDuration != (metav1.Duration{}) {
		o.LeaseDuration = &le.LeaseDuration.Duration
	}

	if o.RenewDeadline == nil && le.RenewDeadline != (metav1.Duration{}) {
		o.RenewDeadline = &le.RenewDeadline.Duration
	}

	if o.RetryPeriod == nil && le.RetryPeriod != (metav1.Duration{}) {
		o.RetryPeriod = &le.RetryPeriod.Duration
	}

	return o
}
