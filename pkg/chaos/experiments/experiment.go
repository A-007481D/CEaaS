package experiments

import (
	"context"

	"github.com/chaos-engineering/controller/pkg/chaos/apis/chaos/v1alpha1"
	"github.com/chaos-engineering/controller/pkg/chaos/experiments/cpu-hog"
	"github.com/chaos-engineering/controller/pkg/chaos/experiments/memory-hog"
	"github.com/chaos-engineering/controller/pkg/chaos/experiments/network-latency"
	"github.com/chaos-engineering/controller/pkg/chaos/experiments/pod-failure"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// ChaosExperiment is the interface that all chaos experiments must implement
type ChaosExperiment interface {
	// Start starts the chaos experiment
	Start(ctx context.Context, experiment *v1alpha1.ChaosExperiment) error
	
	// Stop stops the chaos experiment
	Stop(ctx context.Context, experiment *v1alpha1.ChaosExperiment) error
}

// ExperimentFactory creates a new chaos experiment based on the experiment type
func ExperimentFactory(client kubernetes.Interface, config *rest.Config, experimentType string) ChaosExperiment {
	switch experimentType {
	case "pod-failure":
		return podfailure.NewPodFailureExperiment(client)
	case "network-latency":
		return networklatency.NewNetworkLatencyExperiment(client, config)
	case "cpu-hog":
		return cpuhog.NewCPUHogExperiment(client, config)
	case "memory-hog":
		return memoryhog.NewMemoryHogExperiment(client, config)
	default:
		return nil
	}
}
