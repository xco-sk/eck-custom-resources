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

package v2

import (
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ProjectConfigStatus defines the observed state of ProjectConfig
type ProjectConfigStatus struct {
}

//+kubebuilder:object:root=true

// The types below reproduce the full surface of
// sigs.k8s.io/controller-runtime/pkg/config/v1alpha1.ControllerManagerConfigurationSpec,
// which controller-runtime removed in v0.18.0. ComponentConfig was retired with no
// in-tree successor - upstream's guidance is that each project owns its config type
// (kubernetes-sigs/controller-runtime#895). Field names, JSON tags and types are kept
// identical to the originals so existing config files stay valid verbatim.

// ManagerMetricsSpec configures the metrics endpoint.
type ManagerMetricsSpec struct {
	BindAddress string `json:"bindAddress,omitempty"`
}

// ManagerHealthSpec configures the health probe endpoints.
type ManagerHealthSpec struct {
	HealthProbeBindAddress string `json:"healthProbeBindAddress,omitempty"`
	ReadinessEndpointName  string `json:"readinessEndpointName,omitempty"`
	LivenessEndpointName   string `json:"livenessEndpointName,omitempty"`
}

// ManagerWebhookSpec configures the webhook server.
type ManagerWebhookSpec struct {
	Port    *int   `json:"port,omitempty"`
	Host    string `json:"host,omitempty"`
	CertDir string `json:"certDir,omitempty"`
}

// ManagerLeaderElectionSpec mirrors
// k8s.io/component-base/config/v1alpha1.LeaderElectionConfiguration.
type ManagerLeaderElectionSpec struct {
	LeaderElect       *bool           `json:"leaderElect,omitempty"`
	LeaseDuration     metav1.Duration `json:"leaseDuration,omitempty"`
	RenewDeadline     metav1.Duration `json:"renewDeadline,omitempty"`
	RetryPeriod       metav1.Duration `json:"retryPeriod,omitempty"`
	ResourceLock      string          `json:"resourceLock,omitempty"`
	ResourceName      string          `json:"resourceName,omitempty"`
	ResourceNamespace string          `json:"resourceNamespace,omitempty"`
}

// ManagerControllerSpec mirrors the original ControllerConfigurationSpec.
//
// CacheSyncTimeout is a bare time.Duration rather than metav1.Duration, matching
// upstream: it therefore decodes from an integer nanosecond count, not "30s".
// The metav1.Duration fields above do accept strings like "10m".
type ManagerControllerSpec struct {
	GroupKindConcurrency map[string]int `json:"groupKindConcurrency,omitempty"`
	CacheSyncTimeout     *time.Duration `json:"cacheSyncTimeout,omitempty"`
	RecoverPanic         *bool          `json:"recoverPanic,omitempty"`
}

// ManagerSpec replaces ControllerManagerConfigurationSpec in full.
type ManagerSpec struct {
	SyncPeriod              *metav1.Duration           `json:"syncPeriod,omitempty"`
	LeaderElection          *ManagerLeaderElectionSpec `json:"leaderElection,omitempty"`
	CacheNamespace          string                     `json:"cacheNamespace,omitempty"`
	GracefulShutdownTimeout *metav1.Duration           `json:"gracefulShutDown,omitempty"`
	Controller              *ManagerControllerSpec     `json:"controller,omitempty"`
	Metrics                 ManagerMetricsSpec         `json:"metrics,omitempty"`
	Health                  ManagerHealthSpec          `json:"health,omitempty"`
	Webhook                 ManagerWebhookSpec         `json:"webhook,omitempty"`
}

//+kubebuilder:object:root=true

type ProjectConfig struct {
	metav1.TypeMeta `json:"type"`

	Manager ManagerSpec `json:"manager,omitempty"`

	Elasticsearch ElasticsearchSpec `json:"elasticsearch,omitempty"`
	Kibana        KibanaSpec        `json:"kibana,omitempty"`
}

//+kubebuilder:object:root=true

func init() {
	SchemeBuilder.Register(&ProjectConfig{})
}
