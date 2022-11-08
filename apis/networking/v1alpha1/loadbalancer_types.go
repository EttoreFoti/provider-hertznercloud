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

// LoadBalancerParameters are the configurable fields of a LoadBalancer.
type LoadBalancerParameters struct {
	Name string `json:"name"`

	// Algorithm of the Load Balancer (round_robin or least_connections)
	Algorithm LoadBalancerAlgorithm `json:"algorithm"`

	//+optional
	// User-defined labels (key-value pairs)
	Labels *map[string]string `json:"labels,omitempty"`

	// ID or name of the Load Balancer type this Load Balancer should be created with
	LoadBalancerType string `json:"load_balancer_type"`

	//+optional
	// ID or name of Location to create Load Balancer in
	Location *string `json:"location,omitempty"`

	//+optional
	// ID of the network the Load Balancer should be attached to on creation
	Network *uint64 `json:"network,omitempty"`

	//+optional
	// Name of network zone
	NetworkZone *string `json:"network_zone,omitempty"`

	//+optional
	// Enable or disable the public interface of the Load Balancer
	PublicInterface *bool `json:"public_interface,omitempty"`

	//+optional
	// Array of services
	Services *[]LoadBalancerService `json:"services,omitempty"`

	//+optional
	// Array of targets
	Targets *[]LoadBalancerTarget `json:"targets,omitempty"`
}

type LoadBalancerAlgorithm struct {
	// Type of the algorithm (round_robin or least connection)
	Type string `json:"type"`
}

type LoadBalancerTargetServer struct {
	// ID of the Server
	ID int `json:"id"`
}

type LoadBalancerTargetLabelSelector struct {
	// Label selector
	Selector string `json:"selector"`
}

type LoadBalancerTargetIP struct {
	//IP of a server that belongs to the same customer (public IPv4/IPv6) or private IP in a Subnetwork type vswitch.
	IP string `json:"ip"`
}

type LoadBalancerTargetHealthStatus struct {
	//+optional
	// Port to check
	ListenPort *int `json:"listen_port,omitempty"`

	//+optional
	// Possible enum values:
	Status *string `json:"status,omitempty"`
}

type LoadBalancerTarget struct {

	// Type of the resource (server, label_selector or ip)
	Type string `json:"type"`

	//+optional
	// Server where the traffic should be routed through
	Server *LoadBalancerTargetServer `json:"server,omitempty"`

	//+optional
	// Label selector and a list of selected targets
	LabelSelector *LoadBalancerTargetLabelSelector `json:"label_selector,omitempty"`

	//+optional
	// IP targets where the traffic should be routed through. It is only possible to use the (Public or vSwitch) IPs of Hetzner Online Root Servers belonging to the project owner. IPs belonging to other users are blocked. Additionally IPs belonging to services provided by Hetzner Cloud (Servers, Load Balancers, ...) are blocked as well.
	IP *LoadBalancerTargetIP `json:"ip,omitempty"`

	//+optional
	// List of health statuses of the services on this target
	HealthStatus *[]LoadBalancerTargetHealthStatus `json:"health_status,omitempty"`

	//+optional
	// Use the private network IP instead of the public IP. Default value is false.
	UsePrivateIP *bool `json:"use_private_ip,omitempty"`

	//+optional
	// List of selected targets
	Targets *[]LoadBalancerTarget `json:"targets,omitempty"`
}

type LoadBalancerServiceHTTP struct {

	//+optional
	// Name of the cookie used for sticky sessions
	CookieName *string `json:"cookie_name"`

	//+optional
	// Lifetime of the cookie used for sticky sessions
	CookieLifetime *int `json:"cookie_lifetime"`

	//+optional
	// IDs of the Certificates to use for TLS/SSL termination by the Load Balancer; empty for TLS/SSL passthrough or if protocol is "http"
	Certificates *[]int `json:"certificates"`

	//+optional
	// Redirect HTTP requests to HTTPS. Only available if protocol is "https". Default false
	RedirectHTTP *bool `json:"redirect_http"`

	//+optional
	// Use sticky sessions. Only available if protocol is "http" or "https". Default false
	StickySessions *bool `json:"sticky_sessions"`
}

type LoadBalancerServiceHealthCheckHTTP struct {
	// Host header to send in the HTTP request. May not contain spaces, percent or backslash symbols. Can be null, in that case no host header is sent.
	Domain string `json:"domain"`

	// HTTP path to use for health checks. May not contain literal spaces, use percent-encoding instead.
	Path string `json:"path"`

	//+optional
	// String that must be contained in HTTP response in order to pass the health check
	Response *string `json:"response,omitempty"`

	//+optional
	// List of returned HTTP status codes in order to pass the health check. Supports the wildcards ? for exactly one character and * for multiple ones. The default is to pass the health check for any status code between 2?? and 3??.
	StatusCodes *[]string `json:"status_codes,omitempty"`

	//+optional
	// Use HTTPS for health check
	TLS *bool `json:"tls,omitempty"`
}

type LoadBalancerServiceHealthCheck struct {

	// Type of the health check (tcp or http)
	Protocol string `json:"protocol"`

	// Port the health check will be performed on
	Port int `json:"port"`

	// Time interval in seconds health checks are performed
	Interval int `json:"interval"`

	// Time in seconds after an attempt is considered a timeout
	Timeout int `json:"timeout"`

	// Unsuccessful retries needed until a target is considered unhealthy; an unhealthy target needs the same number of successful retries to become healthy again
	Retries int `json:"retries"`

	//+optional
	// Additional configuration for protocol http
	HTTP *LoadBalancerServiceHealthCheckHTTP `json:"http,omitempty"`
}

// Refers to: hcloud.schema.LoadBalancerService
type LoadBalancerService struct {

	// Protocol of the Load Balancer (http, https, tcp)
	Protocol string `json:"protocol"`

	// Port the Load Balancer listens on
	ListenPort int `json:"listen_port"`

	// Port the Load Balancer will balance to
	DestinationPort int `json:"destination_port"`

	// Is Proxyprotocol enabled or not
	Proxyprotocol bool `json:"proxy_protocol"`

	// Configuration option for protocols http and https
	HTTP LoadBalancerServiceHTTP `json:"http"`

	// Service health check
	HealthCheck LoadBalancerServiceHealthCheck `json:"health_check"`
}

// LoadBalancerObservation are the observable fields of a LoadBalancer.
type LoadBalancerObservation struct {
	ObservableField string `json:"observableField,omitempty"`
}

// A LoadBalancerSpec defines the desired state of a LoadBalancer.
type LoadBalancerSpec struct {
	xpv1.ResourceSpec `json:",inline"`
	ForProvider       LoadBalancerParameters `json:"forProvider"`
}

// A LoadBalancerStatus represents the observed state of a LoadBalancer.
type LoadBalancerStatus struct {
	xpv1.ResourceStatus `json:",inline"`
	AtProvider          LoadBalancerObservation `json:"atProvider,omitempty"`
}

// +kubebuilder:object:root=true

// A LoadBalancer is an example API type.
// +kubebuilder:printcolumn:name="READY",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="SYNCED",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="EXTERNAL-NAME",type="string",JSONPath=".metadata.annotations.crossplane\\.io/external-name"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster,categories={crossplane,managed,hertznercloud}
type LoadBalancer struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   LoadBalancerSpec   `json:"spec"`
	Status LoadBalancerStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// LoadBalancerList contains a list of LoadBalancer
type LoadBalancerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []LoadBalancer `json:"items"`
}

// LoadBalancer type metadata.
var (
	LoadBalancerKind             = reflect.TypeOf(LoadBalancer{}).Name()
	LoadBalancerGroupKind        = schema.GroupKind{Group: Group, Kind: LoadBalancerKind}.String()
	LoadBalancerKindAPIVersion   = LoadBalancerKind + "." + SchemeGroupVersion.String()
	LoadBalancerGroupVersionKind = SchemeGroupVersion.WithKind(LoadBalancerKind)
)

func init() {
	SchemeBuilder.Register(&LoadBalancer{}, &LoadBalancerList{})
}
