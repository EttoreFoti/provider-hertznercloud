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

package certificate

import (
	"context"
	"fmt"

	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	"github.com/crossplane/crossplane-runtime/pkg/connection"
	"github.com/crossplane/crossplane-runtime/pkg/controller"
	"github.com/crossplane/crossplane-runtime/pkg/event"
	"github.com/crossplane/crossplane-runtime/pkg/meta"
	"github.com/crossplane/crossplane-runtime/pkg/ratelimiter"
	"github.com/crossplane/crossplane-runtime/pkg/reconciler/managed"
	"github.com/crossplane/crossplane-runtime/pkg/resource"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/crossplane/provider-hertznercloud/apis/certificates/v1alpha1"
	apisv1alpha1 "github.com/crossplane/provider-hertznercloud/apis/v1alpha1"
	hertznercloud "github.com/crossplane/provider-hertznercloud/internal/clients"
	"github.com/crossplane/provider-hertznercloud/internal/clients/certificate"
	"github.com/crossplane/provider-hertznercloud/internal/controller/features"
)

const (
	errNotCertificate = "managed resource is not a Certificate custom resource"
	errTrackPCUsage   = "cannot track ProviderConfig usage"
	errGetPC          = "cannot get ProviderConfig"
	errGetCreds       = "cannot get credentials"

	errNewClient = "cannot create new Service"
)

// A CertificateService does nothing.
type CertificateService struct {
	client *hcloud.Client
}

var (
	newCertificateService = func(creds []byte) (*CertificateService, error) {
		c, err := hertznercloud.NewClientHertzner(creds)

		if err != nil {
			return nil, err
		}

		return &CertificateService{
			client: c,
		}, nil
	}
)

// Setup adds a controller that reconciles Certificate managed resources.
func Setup(mgr ctrl.Manager, o controller.Options) error {
	name := managed.ControllerName(v1alpha1.CertificateGroupKind)

	cps := []managed.ConnectionPublisher{managed.NewAPISecretPublisher(mgr.GetClient(), mgr.GetScheme())}
	if o.Features.Enabled(features.EnableAlphaExternalSecretStores) {
		cps = append(cps, connection.NewDetailsManager(mgr.GetClient(), apisv1alpha1.StoreConfigGroupVersionKind))
	}

	r := managed.NewReconciler(mgr,
		resource.ManagedKind(v1alpha1.CertificateGroupVersionKind),
		managed.WithExternalConnecter(&connector{
			kube:         mgr.GetClient(),
			usage:        resource.NewProviderConfigUsageTracker(mgr.GetClient(), &apisv1alpha1.ProviderConfigUsage{}),
			newServiceFn: newCertificateService}),
		managed.WithLogger(o.Logger.WithValues("controller", name)),
		managed.WithRecorder(event.NewAPIRecorder(mgr.GetEventRecorderFor(name))),
		managed.WithConnectionPublishers(cps...))

	return ctrl.NewControllerManagedBy(mgr).
		Named(name).
		WithOptions(o.ForControllerRuntime()).
		For(&v1alpha1.Certificate{}).
		Complete(ratelimiter.NewReconciler(name, r, o.GlobalRateLimiter))
}

// A connector is expected to produce an ExternalClient when its Connect method
// is called.
type connector struct {
	kube         client.Client
	usage        resource.Tracker
	newServiceFn func(creds []byte) (*CertificateService, error)
}

// Connect typically produces an ExternalClient by:
// 1. Tracking that the managed resource is using a ProviderConfig.
// 2. Getting the managed resource's ProviderConfig.
// 3. Getting the credentials specified by the ProviderConfig.
// 4. Using the credentials to form a client.
func (c *connector) Connect(ctx context.Context, mg resource.Managed) (managed.ExternalClient, error) {
	cr, ok := mg.(*v1alpha1.Certificate)
	if !ok {
		return nil, errors.New(errNotCertificate)
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
	service *CertificateService
}

func (c *external) Observe(ctx context.Context, mg resource.Managed) (managed.ExternalObservation, error) {
	cr, ok := mg.(*v1alpha1.Certificate)
	if !ok {
		return managed.ExternalObservation{}, errors.New(errNotCertificate)
	}

	// These fmt statements should be removed in the real implementation.
	fmt.Printf("Observing: %+v", cr)

	cert, _, err := c.service.client.Certificate.Get(ctx, meta.GetExternalName(cr))

	if cert == nil && err == nil {
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

	if cert != nil {

		if cert.Type == hcloud.CertificateTypeManaged {
			if cert.Status.Issuance == hcloud.CertificateStatusTypeCompleted {
				cr.Status.SetConditions(xpv1.Available())

				isUpdated, e := certificate.IsCertificateUpToDate(cr.Spec.ForProvider.DeepCopy(), cert)

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
		}

		if cert.Type == hcloud.CertificateTypeUploaded {
			cr.Status.SetConditions(xpv1.Available())

			isUpdated, e := certificate.IsCertificateUpToDate(cr.Spec.ForProvider.DeepCopy(), cert)

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
	cr, ok := mg.(*v1alpha1.Certificate)
	if !ok {
		return managed.ExternalCreation{}, errors.New(errNotCertificate)
	}

	fmt.Printf("Creating: %+v", cr)

	certificateCreateOps, _ := certificate.FromCertificateSpecToCertificateCreateOpts(cr.Spec.ForProvider.DeepCopy(), c.service.client, ctx)
	certificateCreateOps.Name = meta.GetExternalName(cr)

	cert, _, err := c.service.client.Certificate.Create(ctx, *certificateCreateOps)

	if err != nil {
		return managed.ExternalCreation{
			// Optionally return any details that may be required to connect to the
			// external resource. These will be stored as the connection secret.
			ConnectionDetails: managed.ConnectionDetails{},
		}, err
	}

	cr.Status.AtProvider.State = string(cert.Status.Issuance)

	return managed.ExternalCreation{
		// Optionally return any details that may be required to connect to the
		// external resource. These will be stored as the connection secret.
		ConnectionDetails: managed.ConnectionDetails{},
	}, nil
}

func (c *external) Update(ctx context.Context, mg resource.Managed) (managed.ExternalUpdate, error) {
	cr, ok := mg.(*v1alpha1.Certificate)
	if !ok {
		return managed.ExternalUpdate{}, errors.New(errNotCertificate)
	}

	fmt.Printf("Updating: %+v", cr)

	cert, _, err := c.service.client.Certificate.Get(ctx, meta.GetExternalName(cr))

	if err != nil {
		return managed.ExternalUpdate{}, err
	}

	opts := hcloud.CertificateUpdateOpts{
		Labels: *cr.Spec.ForProvider.DeepCopy().Labels,
	}

	c.service.client.Certificate.Update(ctx, cert, opts)

	return managed.ExternalUpdate{
		// Optionally return any details that may be required to connect to the
		// external resource. These will be stored as the connection secret.
		ConnectionDetails: managed.ConnectionDetails{},
	}, nil
}

func (c *external) Delete(ctx context.Context, mg resource.Managed) error {
	cr, ok := mg.(*v1alpha1.Certificate)
	if !ok {
		return errors.New(errNotCertificate)
	}

	fmt.Printf("Deleting: %+v", cr)

	cert, _, err := c.service.client.Certificate.Get(ctx, meta.GetExternalName(cr))

	if err != nil {
		return err
	}

	_, e := c.service.client.Certificate.Delete(ctx, cert)

	if e != nil {
		return e
	}

	return nil
}
