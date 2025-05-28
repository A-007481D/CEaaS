package memoryhog

import (
	"context"
	"fmt"
	"strings"

	"github.com/chaos-engineering/controller/pkg/chaos/apis/chaos/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
	"k8s.io/klog/v2"
)

// MemoryHogExperiment implements the memory hog chaos experiment
type MemoryHogExperiment struct {
	client kubernetes.Interface
	config *rest.Config
}

// NewMemoryHogExperiment creates a new memory hog experiment
func NewMemoryHogExperiment(client kubernetes.Interface, config *rest.Config) *MemoryHogExperiment {
	return &MemoryHogExperiment{
		client: client,
		config: config,
	}
}

// Start starts the memory hog experiment
func (e *MemoryHogExperiment) Start(ctx context.Context, experiment *v1alpha1.ChaosExperiment) error {
	// Get the target pod
	pod, err := e.client.CoreV1().Pods(experiment.Spec.Target.Namespace).Get(ctx, experiment.Spec.Target.Name, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("failed to get target pod: %v", err)
	}

	klog.Infof("Starting memory hog experiment on pod %s/%s", pod.Namespace, pod.Name)

	// Get memory parameter
	memoryMB := "256" // default
	if val, ok := experiment.Spec.Parameters["memoryMB"]; ok {
		memoryMB = val
	}

	// Run stress command to hog memory
	cmd := []string{
		"sh",
		"-c",
		fmt.Sprintf("apt-get update && apt-get install -y stress && stress --vm 1 --vm-bytes %sM --timeout %s", 
			memoryMB, 
			experiment.Spec.Duration),
	}

	req := e.client.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(pod.Name).
		Namespace(pod.Namespace).
		SubResource("exec").
		VersionedParams(&corev1.PodExecOptions{
			Command: cmd,
			Stdin:   false,
			Stdout:  true,
			Stderr:  true,
			TTY:     false,
		}, scheme.ParameterCodec)

	exec, err := remotecommand.NewSPDYExecutor(e.config, "POST", req.URL())
	if err != nil {
		return fmt.Errorf("failed to create executor: %v", err)
	}

	var stdout, stderr strings.Builder
	err = exec.Stream(remotecommand.StreamOptions{
		Stdout: &stdout,
		Stderr: &stderr,
	})
	if err != nil {
		return fmt.Errorf("failed to execute command: %v, stderr: %s", err, stderr.String())
	}

	klog.Infof("Successfully started memory hog on pod %s/%s: %s", pod.Namespace, pod.Name, stdout.String())
	return nil
}

// Stop stops the memory hog experiment
func (e *MemoryHogExperiment) Stop(ctx context.Context, experiment *v1alpha1.ChaosExperiment) error {
	// Get the target pod
	pod, err := e.client.CoreV1().Pods(experiment.Spec.Target.Namespace).Get(ctx, experiment.Spec.Target.Name, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("failed to get target pod: %v", err)
	}

	klog.Infof("Stopping memory hog experiment on pod %s/%s", pod.Namespace, pod.Name)

	// Kill stress process
	cmd := []string{
		"sh",
		"-c",
		"pkill stress || true",
	}

	req := e.client.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(pod.Name).
		Namespace(pod.Namespace).
		SubResource("exec").
		VersionedParams(&corev1.PodExecOptions{
			Command: cmd,
			Stdin:   false,
			Stdout:  true,
			Stderr:  true,
			TTY:     false,
		}, scheme.ParameterCodec)

	exec, err := remotecommand.NewSPDYExecutor(e.config, "POST", req.URL())
	if err != nil {
		return fmt.Errorf("failed to create executor: %v", err)
	}

	var stdout, stderr strings.Builder
	err = exec.Stream(remotecommand.StreamOptions{
		Stdout: &stdout,
		Stderr: &stderr,
	})
	if err != nil {
		return fmt.Errorf("failed to execute command: %v, stderr: %s", err, stderr.String())
	}

	klog.Infof("Successfully stopped memory hog on pod %s/%s: %s", pod.Namespace, pod.Name, stdout.String())
	return nil
}
