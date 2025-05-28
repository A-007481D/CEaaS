package v1alpha1

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/kubernetes/scheme"
)

const (
	GroupName = "chaos.engineering"
	Version   = "v1alpha1"
)

var (
	// SchemeGroupVersion is the group version used to register these objects
	SchemeGroupVersion = schema.GroupVersion{Group: GroupName, Version: Version}

	// SchemeBuilder is used to add go types to the GroupVersionKind scheme
	SchemeBuilder = &scheme.Builder{GroupVersion: SchemeGroupVersion}

	// AddToScheme adds the types in this group-version to the given scheme.
	AddToScheme = SchemeBuilder.AddToScheme
)

// Codecs returns a new codec factory for this scheme
func Codecs() serializer.CodecFactory {
	return serializer.NewCodecFactory(scheme.Scheme)
}

func init() {
	// Register the types with the Scheme so the components can map objects to GroupVersionKinds and back
	SchemeBuilder.Register(&ChaosExperiment{}, &ChaosExperimentList{})
}
