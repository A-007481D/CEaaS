package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Phase constants for experiment status
const (
	PhasePending   = "Pending"
	PhaseRunning   = "Running"
	PhaseCompleted = "Completed"
	PhaseFailed    = "Failed"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ChaosExperiment is the Schema for the chaosexperiments API
type ChaosExperiment struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ChaosExperimentSpec   `json:"spec,omitempty"`
	Status ChaosExperimentStatus `json:"status,omitempty"`
}

// ChaosExperimentSpec defines the desired state of ChaosExperiment
type ChaosExperimentSpec struct {
	// Target defines the target resource for the chaos experiment
	Target TargetResource `json:"target"`
	// ExperimentType is the type of chaos experiment to run
	ExperimentType string `json:"experimentType"`
	// Duration is how long the experiment should run
	Duration string `json:"duration"`
	// Parameters are the parameters for the experiment
	Parameters map[string]string `json:"parameters,omitempty"`
}

// TargetResource defines the target resource for the chaos experiment
type TargetResource struct {
	// API version of the target resource
	APIVersion string `json:"apiVersion"`
	// Kind of the target resource
	Kind string `json:"kind"`
	// Name of the target resource
	Name string `json:"name"`
	// Namespace of the target resource
	Namespace string `json:"namespace"`
}

// ChaosExperimentStatus defines the observed state of ChaosExperiment
type ChaosExperimentStatus struct {
	// Phase represents the current phase of the experiment
	Phase string `json:"phase,omitempty"`
	// StartTime is when the experiment started
	StartTime *metav1.Time `json:"startTime,omitempty"`
	// EndTime is when the experiment ended
	EndTime *metav1.Time `json:"endTime,omitempty"`
	// Message provides more details about the current phase
	Message string `json:"message,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ChaosExperimentList contains a list of ChaosExperiment
type ChaosExperimentList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ChaosExperiment `json:"items"`
}
