/*
Copyright 2021.

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

import (
	"github.com/flanksource/konfig-manager/pkg"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Output defines where and how properties file need to be created
type Output struct {
	Name      string `yaml:"name,omitempty" json:"name,omitempty"`
	Namespace string `yaml:"namespace,omitempty" json:"namespace,omitempty"`
	Kind      string `yaml:"kind,omitempty" json:"kind,omitempty"`
	Type      string `yaml:"type,omitempty" json:"type,omitempty"`
	FileName  string `yaml:"fileName,omitempty" json:"fileName,omitempty"`
}

// KonfigSpec defines the desired state of Konfig
type KonfigSpec struct {
	Hierarchy []pkg.Item `yaml:"hierarchy" json:"hierarchy"`
	Output    Output     `yaml:"output,omitempty" json:"output,omitempty"`
}

// KonfigStatus defines the observed state of Konfig
type KonfigStatus struct {
	Hierarchy []pkg.Item `yaml:"hierarchy" json:"hierarchy"`
	Output    Output     `yaml:"output,omitempty" json:"output,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Konfig is the Schema for the hierarchyconfigs API
type Konfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   KonfigSpec   `json:"spec,omitempty"`
	Status KonfigStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// KonfigList contains a list of Konfig
type KonfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Konfig `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Konfig{}, &KonfigList{})
}
