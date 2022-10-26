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

	// The name of the server; generated from the CR name
	// +immutable
	// +optional
	Name string `json:"name,omitempty"` // TODO check

	// The type of the instance used to create the server (ex: CX11)
	// +immutable
	ServerType string `json:"server_type"`

	// The image used as OS of the server
	// +immutable
	Image string `json:"image"` // name not ID

	// The SSH Keys configured to be used to ssh login to the server
	// +immutable
	// +optional
	SSHKeys *[]string `json:"ssh_keys,omitempty"`

	// The location (Region/Zone) where to put the server (configure either Location or Datacenter)
	// +immutable
	// +optional
	Location *string `json:"location,omitempty"` // name

	// The Datacenter where to put the server (configure either Location or Datacenter)
	// +immutable
	// +optional
	Datacenter *string `json:"datacenter,omitempty"` // name

	// The UserData script to launch on first boot of the server
	// +immutable
	// +optional
	UserData *string `json:"user_data,omitempty"`

	// Whether to start or not the server after creation (default: true)
	// +immutable
	// +optional
	StartAfterCreate *bool `json:"start_after_create,omitempty"`

	// Simple labels to tag the server
	// +optional
	Labels *map[string]string `json:"labels,omitempty"`

	// Whether to mount the Volumes automatically (configurable only if volumes is configured)
	// +optional
	// +immutable
	Automount *bool `json:"automount,omitempty"`

	// Volumes to attach to the server (configurable using Volumes ID)
	// +optional
	// +immutable
	Volumes *[]int `json:"volumes,omitempty"`

	// Networks to attach to the server (use Networks name)
	// +optional
	// +immutable
	Networks *[]string `json:"networks,omitempty"` // name

	// Firewalls to attach to the server (use Firewalls name)
	// +immutable
	// +optional
	Firewalls *[]string `json:"firewalls,omitempty"` // name

	// PlacementGroup to put the server in (Use PlacementGroup ID)
	// +immutable
	// +optional
	PlacementGroup *int `json:"placement_group,omitempty"`

	// Bool to create or not a Public IPv4
	// +immutable
	// +optional
	PublicNetIPv4 *bool `json:"public_net_ipv4,omitempty"` // to fix json definition

	// Bool to create or not a Public IPv4
	// +immutable
	// +optional
	PublicNetIPv6 *bool `json:"public_net_ipv6,omitempty"`
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

// A ServerInstance is a Managed resource that defines a Server on HertznerCloud
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
