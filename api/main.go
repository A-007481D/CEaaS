package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gorilla/mux"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"

	chaosv1alpha1 "github.com/chaos-engineering/controller/pkg/chaos/apis/chaos/v1alpha1"
	chaosclientset "github.com/chaos-engineering/controller/pkg/generated/clientset/versioned"
)

// Server represents the API server
type Server struct {
	KubeClient  kubernetes.Interface
	ChaosClient chaosclientset.Interface
}

// ExperimentRequest represents a request to create a new experiment
type ExperimentRequest struct {
	Name          string            `json:"name"`
	Namespace     string            `json:"namespace"`
	TargetName    string            `json:"targetName"`
	TargetKind    string            `json:"targetKind"`
	ExperimentType string            `json:"experimentType"`
	Duration      string            `json:"duration"`
	Parameters    map[string]string `json:"parameters"`
}

// ExperimentResponse represents an experiment response
type ExperimentResponse struct {
	Name          string            `json:"name"`
	Namespace     string            `json:"namespace"`
	ExperimentType string            `json:"experimentType"`
	Status        string            `json:"status"`
	StartTime     *metav1.Time      `json:"startTime,omitempty"`
	EndTime       *metav1.Time      `json:"endTime,omitempty"`
	Message       string            `json:"message,omitempty"`
}

func main() {
	// Create the Kubernetes client
	var kubeconfig string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = filepath.Join(home, ".kube", "config")
	}

	var config *rest.Config
	var err error

	// Try to use in-cluster config
	config, err = rest.InClusterConfig()
	if err != nil {
		// Fall back to kubeconfig
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			log.Fatalf("Error building kubeconfig: %s", err.Error())
		}
	}

	kubeClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Error building kubernetes clientset: %s", err.Error())
	}

	chaosClient, err := chaosclientset.NewForConfig(config)
	if err != nil {
		log.Fatalf("Error building chaos clientset: %s", err.Error())
	}

	server := &Server{
		KubeClient:  kubeClient,
		ChaosClient: chaosClient,
	}

	// Create the router
	r := mux.NewRouter()

	// API routes
	r.HandleFunc("/api/experiments", server.listExperiments).Methods("GET")
	r.HandleFunc("/api/experiments", server.createExperiment).Methods("POST")
	r.HandleFunc("/api/experiments/{namespace}/{name}", server.getExperiment).Methods("GET")
	r.HandleFunc("/api/experiments/{namespace}/{name}", server.deleteExperiment).Methods("DELETE")

	// Serve static files for the React app
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./dashboard/build")))

	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting server on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), r))
}

// listExperiments lists all chaos experiments
func (s *Server) listExperiments(w http.ResponseWriter, r *http.Request) {
	// Get all experiments from all namespaces
	experiments, err := s.ChaosClient.ChaosV1alpha1().ChaosExperiments("").List(r.Context(), metav1.ListOptions{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Convert to response format
	var response []ExperimentResponse
	for _, exp := range experiments.Items {
		response = append(response, ExperimentResponse{
			Name:          exp.Name,
			Namespace:     exp.Namespace,
			ExperimentType: exp.Spec.ExperimentType,
			Status:        exp.Status.Phase,
			StartTime:     exp.Status.StartTime,
			EndTime:       exp.Status.EndTime,
			Message:       exp.Status.Message,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// getExperiment gets a specific chaos experiment
func (s *Server) getExperiment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	namespace := vars["namespace"]
	name := vars["name"]

	experiment, err := s.ChaosClient.ChaosV1alpha1().ChaosExperiments(namespace).Get(r.Context(), name, metav1.GetOptions{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := ExperimentResponse{
		Name:          experiment.Name,
		Namespace:     experiment.Namespace,
		ExperimentType: experiment.Spec.ExperimentType,
		Status:        experiment.Status.Phase,
		StartTime:     experiment.Status.StartTime,
		EndTime:       experiment.Status.EndTime,
		Message:       experiment.Status.Message,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// createExperiment creates a new chaos experiment
func (s *Server) createExperiment(w http.ResponseWriter, r *http.Request) {
	var req ExperimentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Create the experiment
	experiment := &chaosv1alpha1.ChaosExperiment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      req.Name,
			Namespace: req.Namespace,
		},
		Spec: chaosv1alpha1.ChaosExperimentSpec{
			Target: chaosv1alpha1.TargetResource{
				APIVersion: "v1",
				Kind:       req.TargetKind,
				Name:       req.TargetName,
				Namespace:  req.Namespace,
			},
			ExperimentType: req.ExperimentType,
			Duration:      req.Duration,
			Parameters:    req.Parameters,
		},
	}

	result, err := s.ChaosClient.ChaosV1alpha1().ChaosExperiments(req.Namespace).Create(r.Context(), experiment, metav1.CreateOptions{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := ExperimentResponse{
		Name:          result.Name,
		Namespace:     result.Namespace,
		ExperimentType: result.Spec.ExperimentType,
		Status:        result.Status.Phase,
		StartTime:     result.Status.StartTime,
		EndTime:       result.Status.EndTime,
		Message:       result.Status.Message,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// deleteExperiment deletes a chaos experiment
func (s *Server) deleteExperiment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	namespace := vars["namespace"]
	name := vars["name"]

	err := s.ChaosClient.ChaosV1alpha1().ChaosExperiments(namespace).Delete(r.Context(), name, metav1.DeleteOptions{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
