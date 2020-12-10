/*


Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1

//go:generate controller-gen object paths=$GOFILE

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ManagedOperatorSpec defines the desired state of ManagedOperator
type ManagedOperatorSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of ManagedOperator. Edit ManagedOperator_types.go to remove/update
	Name             string    `json:"name"`
	Version          string    `json:"version"`
	DeploymentName   string    `json:"deploymentname"`
	ServiceType      string    `json:"servicetype"`
	ServiceLabel     string    `json:"servicelabel"`
	DeploymentInputs    []SpecUIGroup `json:"deploymentinputs,omitempty"`
}

type SpecUIGroup struct {
	ControlName string       `json:"controlName"`
	ControlType string       `json:"controlType"`
	ValueType   string       `json:"valueType,omitempty"`
	Placeholder string       `json:"placeholder"`
	Options     []Options    `json:"options,omitempty"`
	Validators  Validators `json:"validators,omitempty"`
}

type Options struct {
	OptionName string `json:"optionName"`
	Value      string `json:"value"`
}

type Validators struct {
	Required  bool `json:"required"`
	Minlength int  `json:"minlength,omitempty"`
	Maxlength int  `json:"maxlength,omitempty"`
}


// ManagedOperatorStatus defines the observed state of ManagedOperator
type ManagedOperatorStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true

// ManagedOperator is the Schema for the managedoperators API
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type ManagedOperator struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ManagedOperatorSpec   `json:"spec,omitempty"`
	Status ManagedOperatorStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ManagedOperatorList contains a list of ManagedOperator
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type ManagedOperatorList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ManagedOperator `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ManagedOperator{}, &ManagedOperatorList{})
}
