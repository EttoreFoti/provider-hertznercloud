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

package loadbalancer

import (
	"context"
	"fmt"

	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	"github.com/crossplane/crossplane-runtime/pkg/connection"
	"github.com/crossplane/crossplane-runtime/pkg/controller"
	"github.com/crossplane/crossplane-runtime/pkg/event"
	"github.com/crossplane/crossplane-runtime/pkg/meta"
	"github.com/crossplane/crossplane-runtime/pkg/ratelimiter"
	"github.com/crossplane/crossplane-runtime/pkg/reconciler/managed"
	"github.com/crossplane/crossplane-runtime/pkg/resource"

	"github.com/crossplane/provider-hertznercloud/apis/networking/v1alpha1"
	apisv1alpha1 "github.com/crossplane/provider-hertznercloud/apis/v1alpha1"
	hertznercloud "github.com/crossplane/provider-hertznercloud/internal/clients"
	"github.com/crossplane/provider-hertznercloud/internal/clients/networking/loadbalancer"
	"github.com/crossplane/provider-hertznercloud/internal/controller/features"
)

const (
	errNotLoadBalancer = "managed resource is not a LoadBalancer custom resource"
	errTrackPCUsage    = "cannot track ProviderConfig usage"
	errGetPC           = "cannot get ProviderConfig"
	errGetCreds        = "cannot get credentials"

	errNewClient = "cannot create new Service"
)

// A LoadBalancerService does nothing.
type LoadBalancerService struct {
	client *hcloud.Client
}

var (
	newLoadBalancerService = func(creds []byte) (*LoadBalancerService, error) {
		c, err := hertznercloud.NewClientHertzner(creds)

		if err != nil {
			return nil, err
		}

		return &LoadBalancerService{
			client: c,
		}, nil
	}
)

// Setup adds a controller that reconciles LoadBalancer managed resources.
func Setup(mgr ctrl.Manager, o controller.Options) error {
	name := managed.ControllerName(v1alpha1.LoadBalancerGroupKind)

	cps := []managed.ConnectionPublisher{managed.NewAPISecretPublisher(mgr.GetClient(), mgr.GetScheme())}
	if o.Features.Enabled(features.EnableAlphaExternalSecretStores) {
		cps = append(cps, connection.NewDetailsManager(mgr.GetClient(), apisv1alpha1.StoreConfigGroupVersionKind))
	}

	r := managed.NewReconciler(mgr,
		resource.ManagedKind(v1alpha1.LoadBalancerGroupVersionKind),
		managed.WithExternalConnecter(&connector{
			kube:         mgr.GetClient(),
			usage:        resource.NewProviderConfigUsageTracker(mgr.GetClient(), &apisv1alpha1.ProviderConfigUsage{}),
			newServiceFn: newLoadBalancerService}),
		managed.WithLogger(o.Logger.WithValues("controller", name)),
		managed.WithRecorder(event.NewAPIRecorder(mgr.GetEventRecorderFor(name))),
		managed.WithConnectionPublishers(cps...))

	return ctrl.NewControllerManagedBy(mgr).
		Named(name).
		WithOptions(o.ForControllerRuntime()).
		For(&v1alpha1.LoadBalancer{}).
		Complete(ratelimiter.NewReconciler(name, r, o.GlobalRateLimiter))
}

// A connector is expected to produce an ExternalClient when its Connect method
// is called.
type connector struct {
	kube         client.Client
	usage        resource.Tracker
	newServiceFn func(creds []byte) (*LoadBalancerService, error)
}

// Connect typically produces an ExternalClient by:
// 1. Tracking that the managed resource is using a ProviderConfig.
// 2. Getting the managed resource's ProviderConfig.
// 3. Getting the credentials specified by the ProviderConfig.
// 4. Using the credentials to form a client.
func (c *connector) Connect(ctx context.Context, mg resource.Managed) (managed.ExternalClient, error) {
	cr, ok := mg.(*v1alpha1.LoadBalancer)
	if !ok {
		return nil, errors.New(errNotLoadBalancer)
	}

	if err := c.usage.Track(ctx, mg); err != nil {
		return nil, errors.Wrap(err, errTrackPCUsage)
	}

	pc := &apisv1alpha1.ProviderConfig{}
	if err := c.kube.Get(ctx, types.NamespacedName{Name: cr.GetProviderConfigReference().Name}, pc); err != nil {
		return nil, errors.Wrap(err, errGetPC)
	}

	cd := pc.Spec.Credentials
	data, err := resource.CommonCredentialExtractor(ctx, cd.Source, c.kube, cd.CommonCredentialSelectors)
	if err != nil {
		return nil, errors.Wrap(err, errGetCreds)
	}

	svc, err := c.newServiceFn(data)
	if err != nil {
		return nil, errors.Wrap(err, errNewClient)
	}

	return &external{service: svc}, nil
}

// An ExternalClient observes, then either creates, updates, or deletes an
// external resource to ensure it reflects the managed resource's desired state.
type external struct {
	// A 'client' used to connect to the external resource API. In practice this
	// would be something like an AWS SDK client.
	service *LoadBalancerService
}

func (c *external) Observe(ctx context.Context, mg resource.Managed) (managed.ExternalObservation, error) {
	cr, ok := mg.(*v1alpha1.LoadBalancer)
	if !ok {
		return managed.ExternalObservation{}, errors.New(errNotLoadBalancer)
	}

	// These fmt statements should be removed in the real implementation.
	fmt.Printf("Observing: %+v", cr)

	lb, _, err := c.service.client.LoadBalancer.Get(ctx, meta.GetExternalName(cr))

	if lb == nil && err == nil {
		return managed.ExternalObservation{
			ResourceExists: false,
		}, nil
	}

	if err != nil {
		if hErr, ok := err.(*hcloud.Error); ok && hErr.Code == hcloud.ErrorCodeNotFound { // this might need improving!!
			return managed.ExternalObservation{
				ResourceExists: false,
			}, nil
		}
	}

	if lb != nil {

		cr.Status.SetConditions(xpv1.Available())

		isUpdated, e := loadbalancer.IsLoadBalancerUpToDate(cr.Spec.ForProvider.DeepCopy(), lb)

		if e != nil {
			return managed.ExternalObservation{}, e
		}

		if !isUpdated && e == nil {
			return managed.ExternalObservation{
				ResourceExists:   true,
				ResourceUpToDate: false,
			}, nil

		}
	}

	return managed.ExternalObservation{
		// Return false when the external resource does not exist. This lets
		// the managed resource reconciler know that it needs to call Create to
		// (re)create the resource, or that it has successfully been deleted.
		ResourceExists: true,

		// Return false when the external resource exists, but it not up to date
		// with the desired managed resource state. This lets the managed
		// resource reconciler know that it needs to call Update.
		ResourceUpToDate: true,

		// Return any details that may be required to connect to the external
		// resource. These will be stored as the connection secret.
		ConnectionDetails: managed.ConnectionDetails{},
	}, nil
}

func (c *external) Create(ctx context.Context, mg resource.Managed) (managed.ExternalCreation, error) {
	cr, ok := mg.(*v1alpha1.LoadBalancer)
	if !ok {
		return managed.ExternalCreation{}, errors.New(errNotLoadBalancer)
	}

	fmt.Println("\n\nCreating:")

	loadBalancerCreateOpts, err := loadbalancer.FromLoadBalancerParametersToLoadBalancerCreateOpts(cr.Spec.ForProvider.DeepCopy(), c.service.client, ctx)

	if err != nil {
		fmt.Println("\n\nError:", err)
	}

	loadBalancerCreateOpts.Name = meta.GetExternalName(cr)

	fmt.Println("\n\nDebug1\n\n", loadBalancerCreateOpts)
	lb, _, err := c.service.client.LoadBalancer.Create(ctx, *loadBalancerCreateOpts)
	fmt.Println("\n\nDebug2", lb)
	if err != nil {
		return managed.ExternalCreation{
			// Optionally return any details that may be required to connect to the
			// external resource. These will be stored as the connection secret.
			ConnectionDetails: managed.ConnectionDetails{},
		}, err
	}

	// TODO
	// cr.Status.AtProvider.State = string(lb.LoadBalancer.Created)

	fmt.Printf("Creating: %+v", lb)

	return managed.ExternalCreation{
		// Optionally return any details that may be required to connect to the
		// external resource. These will be stored as the connection secret.
		ConnectionDetails: managed.ConnectionDetails{},
	}, nil
}

func (c *external) Update(ctx context.Context, mg resource.Managed) (managed.ExternalUpdate, error) {
	cr, ok := mg.(*v1alpha1.LoadBalancer)
	if !ok {
		return managed.ExternalUpdate{}, errors.New(errNotLoadBalancer)
	}

	fmt.Printf("Updating: %+v", cr)

	lb, _, err := c.service.client.LoadBalancer.Get(ctx, meta.GetExternalName(cr))

	if err != nil {
		return managed.ExternalUpdate{}, err
	}

	opts := hcloud.LoadBalancerUpdateOpts{
		Labels: *cr.Spec.ForProvider.DeepCopy().Labels,
	}

	c.service.client.LoadBalancer.Update(ctx, lb, opts)

	return managed.ExternalUpdate{
		// Optionally return any details that may be required to connect to the
		// external resource. These will be stored as the connection secret.
		ConnectionDetails: managed.ConnectionDetails{},
	}, nil
}

func (c *external) Delete(ctx context.Context, mg resource.Managed) error {
	cr, ok := mg.(*v1alpha1.LoadBalancer)
	if !ok {
		return errors.New(errNotLoadBalancer)
	}

	lb, _, err := c.service.client.LoadBalancer.Get(ctx, meta.GetExternalName(cr))

	if err != nil {
		return err
	}

	_, e := c.service.client.LoadBalancer.Delete(ctx, lb)

	if e != nil {
		return e
	}
	fmt.Printf("Deleting: %+v", cr)

	return nil
}
