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

type HttpGzipConditionType string
type ConditionStatus string
type HttpGzipPhase string

const (
	Ready            HttpGzipConditionType = "Ready"
	ConditionTrue    ConditionStatus       = "True"
	ConditionFalse   ConditionStatus       = "False"
	ConditionUnknown ConditionStatus       = "Unknown"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

type Kind string

const (
	Gateway Kind = "gateway"
	Pod     Kind = "pod"
)

type ApplyTo struct {
	// +kubebuilder:validation:Enum={gateway,pod}
	// Kubernetes Kind for which you want ot enable gzip compression.
	// Only supports: (pod, gateway)
	// pod: Kubernetes Pod,
	// gateway: Istio Gateway
	Kind Kind `json:"kind"`
	// Label selector to select pods or gateways for which
	// you want to enable gzip compression
	Selector map[string]string `json:"selector"`
}

// HttpGzipSpec defines the desired state of HttpGzip
type HttpGzipSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// ApplyTo defines the workloads you want to
	// enable gzip compression for
	ApplyTo ApplyTo `json:"applyTo,omitempty"`
}

type HttpGzipCondition struct {
	// Type of condition (only `Ready` is supported as of now)
	Type HttpGzipConditionType `json:"type"`
	// Either of (True, False, Unknown)
	Status ConditionStatus `json:"status"`
	// Last time the resource was synced
	LastTransitionTime metav1.Time `json:"lastTransitionTime"`
}

// HttpGzipStatus defines the observed state of HttpGzip
type HttpGzipStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Shows the current state of the resource
	Conditions []HttpGzipCondition `json:"conditions"`
	// Name of the EnvoyFilter resource created
	EnvoyFilter string `json:"envoyFilter,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// HttpGzip is used to enable http gzip compression
// for your services (uses Istio's EnvoyFilter custom resource under the hood)
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
