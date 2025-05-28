package controller

import (
	"context"
	"fmt"
	"time"

	"github.com/chaos-engineering/controller/pkg/chaos/apis/chaos/v1alpha1"
	"github.com/chaos-engineering/controller/pkg/chaos/experiments"
	clientset "github.com/chaos-engineering/controller/pkg/generated/clientset/versioned"
	informers "github.com/chaos-engineering/controller/pkg/generated/informers/externalversions/chaos/v1alpha1"
	listers "github.com/chaos-engineering/controller/pkg/generated/listers/chaos/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog/v2"
)

type Controller struct {
	kubeclientset    kubernetes.Interface
	chaosclientset   clientset.Interface
	restConfig      *rest.Config

	experimentsLister listers.ChaosExperimentLister
	experimentsSynced cache.InformerSynced

	workqueue workqueue.RateLimitingInterface

	// activeExperiments keeps track of running experiments
	activeExperiments map[string]experiments.ChaosExperiment
}

func NewController(
	kubeclientset kubernetes.Interface,
	chaosclientset clientset.Interface,
	restConfig *rest.Config,
	experimentInformer informers.ChaosExperimentInformer) *Controller {

	controller := &Controller{
		kubeclientset:    kubeclientset,
		chaosclientset:   chaosclientset,
		restConfig:      restConfig,
		experimentsLister: experimentInformer.Lister(),
		experimentsSynced: experimentInformer.Informer().HasSynced,
		workqueue:        workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "ChaosExperiments"),
		activeExperiments: make(map[string]experiments.ChaosExperiment),
	}

	klog.Info("Setting up event handlers")
	experimentInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: controller.enqueueChaosExperiment,
		UpdateFunc: func(old, new interface{}) {
			controller.enqueueChaosExperiment(new)
		},
	})

	return controller
}

func (c *Controller) enqueueChaosExperiment(obj interface{}) {
	var key string
	var err error
	if key, err = cache.MetaNamespaceKeyFunc(obj); err != nil {
		runtime.HandleError(err)
		return
	}
	c.workqueue.Add(key)
}

func (c *Controller) Run(threadiness int, stopCh <-chan struct{}) error {
	defer runtime.HandleCrash()
	defer c.workqueue.ShutDown()

	klog.Info("Starting Chaos Controller")

	if ok := cache.WaitForCacheSync(stopCh, c.experimentsSynced); !ok {
		return fmt.Errorf("failed to wait for caches to sync")
	}

	klog.Info("Starting workers")
	for i := 0; i < threadiness; i++ {
		go wait.Until(c.runWorker, time.Second, stopCh)
	}

	klog.Info("Started workers")
	<-stopCh
	klog.Info("Shutting down workers")

	return nil
}

func (c *Controller) runWorker() {
	for c.processNextWorkItem() {
	}
}

func (c *Controller) processNextWorkItem() bool {
	obj, shutdown := c.workqueue.Get()

	if shutdown {
		return false
	}

	err := func(obj interface{}) error {
		defer c.workqueue.Done(obj)
		var key string
		var ok bool

		if key, ok = obj.(string); !ok {
			c.workqueue.Forget(obj)
			runtime.HandleError(fmt.Errorf("expected string in workqueue but got %#v", obj))
			return nil
		}

		if err := c.syncHandler(key); err != nil {
			return fmt.Errorf("error syncing '%s': %s", key, err.Error())
		}

		c.workqueue.Forget(obj)
		klog.Infof("Successfully synced '%s'", key)
		return nil
	}(obj)

	if err != nil {
		runtime.HandleError(err)
		return true
	}

	return true
}

func (c *Controller) syncHandler(key string) error {
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		runtime.HandleError(fmt.Errorf("invalid resource key: %s", key))
		return nil
	}

	experiment, err := c.experimentsLister.ChaosExperiments(namespace).Get(name)
	if err != nil {
		if errors.IsNotFound(err) {
			runtime.HandleError(fmt.Errorf("chaosexperiment '%s' in work queue no longer exists", key))
			return nil
		}
		return err
	}

	switch experiment.Status.Phase {
	case "", v1alpha1.PhasePending:
		return c.handlePendingExperiment(experiment)
	case v1alpha1.PhaseRunning:
		return c.handleRunningExperiment(experiment)
	case v1alpha1.PhaseCompleted, v1alpha1.PhaseFailed:
		// Nothing to do for completed/failed experiments
		return nil
	default:
		return fmt.Errorf("unknown experiment phase: %s", experiment.Status.Phase)
	}
}

func (c *Controller) handlePendingExperiment(experiment *v1alpha1.ChaosExperiment) error {
	experiment = experiment.DeepCopy()
	experiment.Status.Phase = v1alpha1.PhaseRunning
	experiment.Status.StartTime = &metav1.Time{Time: time.Now()}
	experiment.Status.Message = "Experiment started"

	_, err := c.chaosclientset.ChaosV1alpha1().ChaosExperiments(experiment.Namespace).UpdateStatus(context.TODO(), experiment, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("failed to update experiment status: %v", err)
	}

	// Create and start the chaos experiment
	klog.Infof("Starting chaos experiment %s/%s of type %s", experiment.Namespace, experiment.Name, experiment.Spec.ExperimentType)
	
	// Create the experiment
	experimentImpl := experiments.ExperimentFactory(c.kubeclientset, c.restConfig, experiment.Spec.ExperimentType)
	if experimentImpl == nil {
		experiment.Status.Phase = v1alpha1.PhaseFailed
		experiment.Status.Message = fmt.Sprintf("Unknown experiment type: %s", experiment.Spec.ExperimentType)
		_, err := c.chaosclientset.ChaosV1alpha1().ChaosExperiments(experiment.Namespace).UpdateStatus(context.TODO(), experiment, metav1.UpdateOptions{})
		return fmt.Errorf("unknown experiment type: %s", experiment.Spec.ExperimentType)
	}

	// Start the experiment
	err = experimentImpl.Start(context.TODO(), experiment)
	if err != nil {
		experiment.Status.Phase = v1alpha1.PhaseFailed
		experiment.Status.Message = fmt.Sprintf("Failed to start experiment: %v", err)
		_, updateErr := c.chaosclientset.ChaosV1alpha1().ChaosExperiments(experiment.Namespace).UpdateStatus(context.TODO(), experiment, metav1.UpdateOptions{})
		if updateErr != nil {
			klog.Errorf("Failed to update experiment status: %v", updateErr)
		}
		return fmt.Errorf("failed to start experiment: %v", err)
	}

	// Store the experiment in the active experiments map
	key := fmt.Sprintf("%s/%s", experiment.Namespace, experiment.Name)
	c.activeExperiments[key] = experimentImpl

	return nil
}

func (c *Controller) handleRunningExperiment(experiment *v1alpha1.ChaosExperiment) error {
	duration, err := time.ParseDuration(experiment.Spec.Duration)
	if err != nil {
		return fmt.Errorf("invalid duration: %v", err)
	}

	// Check if the experiment has completed
	if time.Since(experiment.Status.StartTime.Time) >= duration {
		experiment = experiment.DeepCopy()
		experiment.Status.Phase = v1alpha1.PhaseCompleted
		now := metav1.Now()
		experiment.Status.EndTime = &now
		experiment.Status.Message = "Experiment completed successfully"

		// Get the experiment from the active experiments map
		key := fmt.Sprintf("%s/%s", experiment.Namespace, experiment.Name)
		experimentImpl, exists := c.activeExperiments[key]
		if exists {
			// Stop the experiment
			klog.Infof("Stopping chaos experiment %s/%s", experiment.Namespace, experiment.Name)
			err = experimentImpl.Stop(context.TODO(), experiment)
			if err != nil {
				klog.Errorf("Failed to stop experiment %s/%s: %v", experiment.Namespace, experiment.Name, err)
				experiment.Status.Message = fmt.Sprintf("Experiment completed with errors: %v", err)
			}
			
			// Remove the experiment from the active experiments map
			delete(c.activeExperiments, key)
		} else {
			klog.Warningf("Experiment %s/%s not found in active experiments map", experiment.Namespace, experiment.Name)
		}

		_, err = c.chaosclientset.ChaosV1alpha1().ChaosExperiments(experiment.Namespace).UpdateStatus(context.TODO(), experiment, metav1.UpdateOptions{})
		if err != nil {
			return fmt.Errorf("failed to update experiment status: %v", err)
		}

		klog.Infof("Completed chaos experiment %s/%s", experiment.Namespace, experiment.Name)
	}

	return nil
}
