package networklatency

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

// NetworkLatencyExperiment implements the network latency chaos experiment
type NetworkLatencyExperiment struct {
	client kubernetes.Interface
	config *rest.Config
}

// NewNetworkLatencyExperiment creates a new network latency experiment
func NewNetworkLatencyExperiment(client kubernetes.Interface, config *rest.Config) *NetworkLatencyExperiment {
	return &NetworkLatencyExperiment{
		client: client,
		config: config,
	}
}

// Start starts the network latency experiment
func (e *NetworkLatencyExperiment) Start(ctx context.Context, experiment *v1alpha1.ChaosExperiment) error {
	// Get the target pod
	pod, err := e.client.CoreV1().Pods(experiment.Spec.Target.Namespace).Get(ctx, experiment.Spec.Target.Name, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("failed to get target pod: %v", err)
	}

	klog.Infof("Starting network latency experiment on pod %s/%s", pod.Namespace, pod.Name)

	// Get latency parameter
	latency := "100ms" // default
	if val, ok := experiment.Spec.Parameters["latency"]; ok {
		latency = val
	}

	// Add network latency using tc
	cmd := []string{
		"sh",
		"-c",
		fmt.Sprintf("tc qdisc add dev eth0 root netem delay %s", latency),
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

	klog.Infof("Successfully added network latency to pod %s/%s: %s", pod.Namespace, pod.Name, stdout.String())
	return nil
}

// Stop stops the network latency experiment
func (e *NetworkLatencyExperiment) Stop(ctx context.Context, experiment *v1alpha1.ChaosExperiment) error {
	// Get the target pod
	pod, err := e.client.CoreV1().Pods(experiment.Spec.Target.Namespace).Get(ctx, experiment.Spec.Target.Name, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("failed to get target pod: %v", err)
	}

	klog.Infof("Stopping network latency experiment on pod %s/%s", pod.Namespace, pod.Name)

	// Remove network latency using tc
	cmd := []string{
		"sh",
		"-c",
		"tc qdisc del dev eth0 root",
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

	klog.Infof("Successfully removed network latency from pod %s/%s: %s", pod.Namespace, pod.Name, stdout.String())
	return nil
}
