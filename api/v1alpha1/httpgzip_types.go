/*
Copyright 2022.

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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type HttpGzipConditionType int
type ConditionStatus string
type HttpGzipPhase string

const (
	Ready            HttpGzipConditionType = 1
	Provisioned      HttpGzipConditionType = 2
	ConditionTrue    ConditionStatus       = "True"
	ConditionFalse   ConditionStatus       = "False"
	ConditionUnknown ConditionStatus       = "Unknown"
	Provisioning     HttpGzipPhase         = "Provisioning"
	Active           HttpGzipPhase         = "Active"
	Terminating      HttpGzipPhase         = "Terminating"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

type Kind string

const (
	Gateway Kind = "gateway"
	Pod     Kind = "pod"
)

type ApplyTo struct {
	Kind     Kind              `json:"kind"`
	Selector map[string]string `json:"selector"`
}

// HttpGzipSpec defines the desired state of HttpGzip
type HttpGzipSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	ApplyTo ApplyTo `json:"applyTo,omitempty"`
}

type HttpGzipCondition struct {
	Type   HttpGzipConditionType `json:"type"`
	Status ConditionStatus       `json:"status"`
	// +optional
	Reason string `json:"reason"`
	// +optional
	Message            string      `json:"message"`
	LastTransitionTime metav1.Time `json:"lastTransitionTime"`
}

// HttpGzipStatus defines the observed state of HttpGzip
type HttpGzipStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Phase       HttpGzipPhase       `json:"phase"`
	Conditions  []HttpGzipCondition `json:"conditions"`
	EnvoyFilter string              `json:"envoyFilter"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// HttpGzip is the Schema for the httpgzips API
type HttpGzip struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   HttpGzipSpec   `json:"spec,omitempty"`
	Status HttpGzipStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// HttpGzipList contains a list of HttpGzip
type HttpGzipList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []HttpGzip `json:"items"`
}

func init() {
	SchemeBuilder.Register(&HttpGzip{}, &HttpGzipList{})
}
