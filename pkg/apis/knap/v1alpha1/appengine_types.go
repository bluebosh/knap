package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	//pipelinev1alpha1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1alpha1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// AppengineSpec defines the desired state of Appengine
// +k8s:openapi-gen=true
type AppengineSpec struct {
	// The application name
	AppName string `json:"appName"`

	// The git repo of the application
	GitRepo string `json:"gitRepo"`

	// The git revision of the application
	GitRevision string `json:"gitRevision"`

	// The application instance count
	Size int32 `json:"size"`

	// The pipeline template of the application
	PipelineTemplate string `json:"pipelineTemplate"`
}

// AppengineStatus defines the observed state of Appengine
// +k8s:openapi-gen=true
type AppengineStatus struct {
	// The application status
	Status string `json:"status"`

	// The application ready
	Ready string `json:"ready"`

	// The pipeline of the application
	// PipelineRun *pipelinev1alpha1.PipelineRun `json:"pipelineRun"`
	Domain string `json:"domain"`

	Instance int32 `json:"instance"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Appengine is the Schema for the appengines API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
type Appengine struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AppengineSpec   `json:"spec,omitempty"`
	Status AppengineStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AppengineList contains a list of Appengine
type AppengineList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Appengine `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Appengine{}, &AppengineList{})
}
