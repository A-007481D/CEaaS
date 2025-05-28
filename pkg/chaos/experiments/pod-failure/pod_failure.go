package podfailure

import (
	"context"
	"fmt"
	"time"

	"github.com/chaos-engineering/controller/pkg/chaos/apis/chaos/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/klog/v2"
)

// PodFailureExperiment implements the pod failure chaos experiment
type PodFailureExperiment struct {
	client kubernetes.Interface
}

// NewPodFailureExperiment creates a new pod failure experiment
func NewPodFailureExperiment(client kubernetes.Interface) *PodFailureExperiment {
	return &PodFailureExperiment{
		client: client,
	}
}

// Start starts the pod failure experiment
func (e *PodFailureExperiment) Start(ctx context.Context, experiment *v1alpha1.ChaosExperiment) error {
	// Get the target pod
	pod, err := e.client.CoreV1().Pods(experiment.Spec.Target.Namespace).Get(ctx, experiment.Spec.Target.Name, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("failed to get target pod: %v", err)
	}

	klog.Infof("Starting pod failure experiment on pod %s/%s", pod.Namespace, pod.Name)

	// Delete the pod to simulate failure
	err = e.client.CoreV1().Pods(pod.Namespace).Delete(ctx, pod.Name, metav1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete pod: %v", err)
	}

	klog.Infof("Successfully deleted pod %s/%s", pod.Namespace, pod.Name)
	return nil
}

// Stop stops the pod failure experiment
func (e *PodFailureExperiment) Stop(ctx context.Context, experiment *v1alpha1.ChaosExperiment) error {
	// Nothing to do here, the pod will be recreated by its controller
	klog.Infof("Pod failure experiment completed for %s/%s", experiment.Spec.Target.Namespace, experiment.Spec.Target.Name)
	return nil
}
