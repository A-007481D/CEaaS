package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/chaos-engineering/controller/pkg/controller"
	"github.com/chaos-engineering/controller/pkg/chaos/apis/chaos/v1alpha1"
	clientset "github.com/chaos-engineering/controller/pkg/generated/clientset/versioned"
	informers "github.com/chaos-engineering/controller/pkg/generated/informers/externalversions"
	kubeinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"
)

var (
	masterURL  string
	kubeconfig string
)

func main() {
	klog.InitFlags(nil)
	flag.Parse()

	// Set up signals so we handle the first shutdown signal gracefully
	stopCh := setupSignalHandler()

	cfg, err := clientcmd.BuildConfigFromFlags(masterURL, kubeconfig)
	if err != nil {
		klog.Fatalf("Error building kubeconfig: %s", err.Error())
	}

	kubeClient, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		klog.Fatalf("Error building kubernetes clientset: %s", err.Error())
	}

	chaosClient, err := clientset.NewForConfig(cfg)
	if err != nil {
		klog.Fatalf("Error building chaos clientset: %s", err.Error())
	}

	kubeInformerFactory := kubeinformers.NewSharedInformerFactory(kubeClient, time.Second*30)
	chaosInformerFactory := informers.NewSharedInformerFactory(chaosClient, time.Second*30)

	controller := controller.NewController(
		kubeClient,
		chaosClient,
		cfg,
		chaosInformerFactory.Chaos().V1alpha1().ChaosExperiments(),
	)

	// Start the informer factories
	go kubeInformerFactory.Start(stopCh)
	go chaosInformerFactory.Start(stopCh)

	// Start the controller
	if err = controller.Run(2, stopCh); err != nil {
		klog.Fatalf("Error running controller: %s", err.Error())
	}
}

func init() {
	flag.StringVar(&kubeconfig, "kubeconfig", "", "Path to a kubeconfig. Only required if out-of-cluster.")
	flag.StringVar(&masterURL, "master", "", "The address of the Kubernetes API server. Overrides any value in kubeconfig. Only required if out-of-cluster.")
}

func setupSignalHandler() (stopCh <-chan struct{}) {
	stop := make(chan struct{})
	c := make(chan os.Signal, 2)
	
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C) or SIGTERM
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		close(stop)
		<-c
		os.Exit(1) // second signal. Exit directly.
	}()

	return stop
}
