package serverinstance

import (
	"context"
	"fmt"

	hertznercloudclient "github.com/crossplane/provider-hertznercloud/internal/clients"

	"github.com/crossplane/provider-hertznercloud/apis/server/v1alpha1"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/hetznercloud/hcloud-go/hcloud"
)

type ServerClient interface {
	GetByID(ctx context.Context, id int) (*hcloud.Server, *hcloud.Response, error)
	GetByName(ctx context.Context, name string) (*hcloud.Server, *hcloud.Response, error)
	Get(ctx context.Context, idOrName string) (*hcloud.Server, *hcloud.Response, error)
	Create(ctx context.Context, opts hcloud.ServerCreateOpts) (hcloud.ServerCreateResult, *hcloud.Response, error)
	Delete(ctx context.Context, server *hcloud.Server) (*hcloud.Response, error)
	Update(ctx context.Context, server *hcloud.Server, opts hcloud.ServerUpdateOpts) (*hcloud.Server, *hcloud.Response, error)
}

func NewServerClient(creds []byte) (ServerClient, error) {
	c, err := hertznercloudclient.NewClientHertzner(creds)
	if err != nil {
		return nil, err
	}
	return &c.HertznerClient.Server, nil //, &c.Server
}

func FromServerSpecToServerRequestOpts(in v1alpha1.ServerInstanceParameters, c hertznercloudclient.Client, ctx context.Context) (*hcloud.ServerCreateOpts, error) {

	serverType, _, err := c.HertznerClient.ServerType.GetByName(ctx, in.ServerType)

	if err != nil {
		return nil, err
	}

	image, _, err := c.HertznerClient.Image.GetByName(ctx, in.Image)

	if err != nil {
		return nil, err
	}

	res := &hcloud.ServerCreateOpts{
		ServerType: serverType,
		Image:      image,
	}

	if in.SSHKeys != nil {
		var ssh_keys []*hcloud.SSHKey
		for _, key_id := range *in.SSHKeys {
			key, _, err := c.HertznerClient.SSHKey.GetByName(ctx, key_id)
			if err != nil {
				return nil, err
			}
			ssh_keys = append(ssh_keys, key)
		}
		res.SSHKeys = ssh_keys
		// fmt.Println("ssh_keys:", ssh_keys)
	}

	if in.Location != nil {
		location, _, err := c.HertznerClient.Location.GetByName(ctx, *in.Location)
		if err != nil {
			return nil, err
		}
		res.Location = location
		// fmt.Println("location:", location)
	}

	if in.Datacenter != nil {
		datacenter, _, err := c.HertznerClient.Datacenter.GetByName(ctx, *(in.Datacenter))
		if err != nil {
			return nil, err
		}
		res.Datacenter = datacenter
		// fmt.Println("datacenter:", datacenter)
	}

	if in.Volumes != nil {
		var volumes []*hcloud.Volume
		for _, volume := range *in.Volumes {
			volume, _, err := c.HertznerClient.Volume.GetByID(ctx, volume)
			if err != nil {
				return nil, err
			}
			volumes = append(volumes, volume)
		}
		res.Volumes = volumes
		// fmt.Println("volumes:", volumes)
	}

	if in.Networks != nil {
		var networks []*hcloud.Network
		for _, network := range *in.Networks {
			network, _, err := c.HertznerClient.Network.GetByName(ctx, network)
			if err != nil {
				return nil, err
			}
			networks = append(networks, network)
		}
		res.Networks = networks
		// fmt.Println("networks:", networks)
	}

	if in.Firewalls != nil {
		var firewalls []*hcloud.ServerCreateFirewall
		for _, firewall := range *in.Firewalls {
			firewall, _, err := c.HertznerClient.Firewall.GetByName(ctx, firewall)
			if err != nil {
				return nil, err
			}
			firewalls = append(firewalls, &hcloud.ServerCreateFirewall{Firewall: *firewall})
		}
		res.Firewalls = firewalls
		// fmt.Println("firewalls:", firewalls)
	}

	if in.PlacementGroup != nil {
		var placementGroup *hcloud.PlacementGroup
		placementGroup, _, err = c.HertznerClient.PlacementGroup.GetByID(ctx, *(in.PlacementGroup))
		if err != nil {
			return nil, err
		}
		res.PlacementGroup = placementGroup
		// fmt.Println("placementGroup:", placementGroup)
	}

	if in.UserData != nil {
		userData := *in.UserData
		// fmt.Println("userData:", userData)
		res.UserData = userData
	}

	if in.Labels != nil {
		labels := *in.Labels
		res.Labels = labels
		// fmt.Println("labels:", labels)
	}

	return res, nil

}

func IsUpToDate(in *v1alpha1.ServerInstanceParameters, o *hcloud.Server) (bool, error) { //TODO update ignore fields because there are only a few fields we can update

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
		// cmpopts.IgnoreFields(hcloud.Server{}, "ID", "Status", "Created", "PublicNet", "PrivateNet", "ServerType", "Datacenter", "IncludedTraffic", "OutgoingTraffic", "IngoingTraffic",
		// 	"BackupWindow", "RescueEnabled", "Locked", "ISO", "Image", "Protection", "Volumes", "PrimaryDiskSize", "PlacementGroup"),
		// cmpopts.IgnoreFields(v1alpha1.ServerInstanceParameters{}, "ServerType", "Image", "SSHKeys", "Location", "Datacenter", "UserData", "StartAfterCreate", "Labels", "Automount",
		// 	"Volumes", "Networks", "Firewalls", "PlacementGroup", "PublicNet"),
	))

	if diff != "" {
		fmt.Printf("\n\n%s\n\n", diff)
		return false, nil
	}

	return true, nil
}
