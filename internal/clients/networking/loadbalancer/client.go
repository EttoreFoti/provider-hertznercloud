package loadbalancer

import (
	"context"
	"fmt"
	"time"

	"github.com/crossplane/provider-hertznercloud/apis/networking/v1alpha1"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/hetznercloud/hcloud-go/hcloud"
)

// type LoadBalancerCreateOpts struct {
//     Name             string
//     LoadBalancerType *LoadBalancerType
//     Algorithm        *LoadBalancerAlgorithm
//     Location         *Location
//     NetworkZone      NetworkZone
//     Labels           map[string]string
//     Targets          []LoadBalancerCreateOptsTarget
//     Services         []LoadBalancerCreateOptsService
//     PublicInterface  *bool
//     Network          *Network
// }

func FromLoadBalancerParametersToLoadBalancerCreateOpts(in *v1alpha1.LoadBalancerParameters, c *hcloud.Client, ctx context.Context) (*hcloud.LoadBalancerCreateOpts, error) {

	res := &hcloud.LoadBalancerCreateOpts{
		Name: in.Name,
	}

	t, _, err := c.LoadBalancerType.Get(ctx, in.LoadBalancerType)

	if err != nil {
		return nil, err
	}

	res.LoadBalancerType = t

	res.Algorithm.Type = hcloud.LoadBalancerAlgorithmType(in.Algorithm.Type)

	if in.Location != nil {
		l, _, err := c.Location.Get(ctx, *in.Location)
		if err != nil {
			return nil, err
		}
		res.Location = l
	}

	if in.NetworkZone != nil {
		res.NetworkZone = hcloud.NetworkZone(*in.NetworkZone)
	}

	if in.Labels != nil {
		labels := *in.Labels
		res.Labels = labels
	}

	if in.Services != nil {
		var services []hcloud.LoadBalancerCreateOptsService
		for _, service := range *in.Services {
			s := hcloud.LoadBalancerCreateOptsService{
				Protocol:        hcloud.LoadBalancerServiceProtocol(service.Protocol),
				ListenPort:      &service.ListenPort,
				DestinationPort: &service.DestinationPort,
				Proxyprotocol:   &service.Proxyprotocol,
			}
			var certificates []*hcloud.Certificate
			for _, ct := range *service.HTTP.Certificates {
				crt, _, err := c.Certificate.GetByID(ctx, ct)
				if err != nil {
					return nil, err
				}
				certificates = append(certificates, crt)
			}
			t := time.Duration(*service.HTTP.CookieLifetime)
			s.HTTP = &hcloud.LoadBalancerCreateOptsServiceHTTP{
				Certificates:   certificates,
				CookieName:     service.HTTP.CookieName,
				CookieLifetime: &t,
				RedirectHTTP:   service.HTTP.RedirectHTTP,
				StickySessions: service.HTTP.StickySessions,
			}
			services = append(services, s)
		}
		res.Services = services
	}

	if in.PublicInterface != nil {
		res.PublicInterface = in.PublicInterface
	}

	if in.Network != nil {
		n, _, err := c.Network.GetByID(ctx, int(*in.Network))
		if err != nil {
			return nil, err
		}
		res.Network = n
	}

	return res, nil

}
func IsLoadBalancerUpToDate(in *v1alpha1.LoadBalancerParameters, o *hcloud.LoadBalancer) (bool, error) {

	// Labels is the only updatable fields

	if in.Labels == nil {
		return true, nil
	}

	diff := (cmp.Diff(*in.Labels, o.Labels,
		cmpopts.EquateEmpty(),
		cmp.Comparer(func(a, b *bool) bool {
			if a == nil {
				return (b == nil) || (!*b)
			}

			if b == nil {
				return (a == nil) || (!*a)
			}

			return *a == *b
		}),
		cmpopts.SortSlices(func(a, b string) bool {
			return a < b
		}),
	))

	if diff != "" {
		fmt.Printf("\n\n%s\n\n", diff)
		return false, nil
	}

	return true, nil
}
