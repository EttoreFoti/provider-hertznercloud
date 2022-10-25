/*
Copyright 2022 The Crossplane Authors.

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
	"reflect"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"

	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
)

// ServerInstanceParameters are the configurable fields of a ServerInstance.
type ServerInstanceParameters struct {
	Name             string             `json:"name,omitempty"` // TODO check
	ServerType       string             `json:"server_type"`
	Image            string             `json:"image"` // name not ID
	SSHKeys          *[]string          `json:"ssh_keys,omitempty"`
	Location         *string            `json:"location,omitempty"`   // name
	Datacenter       *string            `json:"datacenter,omitempty"` // name
	UserData         *string            `json:"user_data,omitempty"`
	StartAfterCreate *bool              `json:"start_after_create,omitempty"`
	Labels           *map[string]string `json:"labels,omitempty"`
	Automount        *bool              `json:"automount,omitempty"`
	Volumes          *[]int             `json:"volumes,omitempty"`
	Networks         *[]string          `json:"networks,omitempty"`  // nmae
	Firewalls        *[]string          `json:"firewalls,omitempty"` // name
	PlacementGroup   *int               `json:"placement_group,omitempty"`
	PublicNetIPv4    *bool              `json:"public_net_ipv4,omitempty"` // to fix json definition
	PublicNetIPv6    *bool              `json:"public_net_ipv6,omitempty"`
}

// ServerInstanceObservation are the observable fields of a ServerInstance.
type ServerInstanceObservation struct {
	State string `json:"status,omitempty"`
}

// A ServerInstanceSpec defines the desired state of a ServerInstance.
type ServerInstanceSpec struct {
	xpv1.ResourceSpec `json:",inline"`
	ForProvider       ServerInstanceParameters `json:"forProvider"`
}

// A ServerInstanceStatus represents the observed state of a ServerInstance.
type ServerInstanceStatus struct {
	xpv1.ResourceStatus `json:",inline"`
	AtProvider          ServerInstanceObservation `json:"atProvider,omitempty"`
}

// +kubebuilder:object:root=true

// A ServerInstance is an example API type.
// +kubebuilder:printcolumn:name="READY",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="SYNCED",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="EXTERNAL-NAME",type="string",JSONPath=".metadata.annotations.crossplane\\.io/external-name"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster,categories={crossplane,managed,hertznercloud}
type ServerInstance struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ServerInstanceSpec   `json:"spec"`
	Status ServerInstanceStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ServerInstanceList contains a list of ServerInstance
type ServerInstanceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ServerInstance `json:"items"`
}

// ServerInstance type metadata.
var (
	ServerInstanceKind             = reflect.TypeOf(ServerInstance{}).Name()
	ServerInstanceGroupKind        = schema.GroupKind{Group: Group, Kind: ServerInstanceKind}.String()
	ServerInstanceKindAPIVersion   = ServerInstanceKind + "." + SchemeGroupVersion.String()
	ServerInstanceGroupVersionKind = SchemeGroupVersion.WithKind(ServerInstanceKind)
)

func init() {
	SchemeBuilder.Register(&ServerInstance{}, &ServerInstanceList{})
}
