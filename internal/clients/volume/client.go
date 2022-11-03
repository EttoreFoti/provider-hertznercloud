package volume

import (
	"context"
	"fmt"

	"github.com/crossplane/provider-hertznercloud/apis/volumes/v1alpha1"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/hetznercloud/hcloud-go/hcloud"
)

func FromVolumeSpecToVolumeCreateOpts(in *v1alpha1.VolumeParameters, c *hcloud.Client, ctx context.Context) (*hcloud.VolumeCreateOpts, error) {

	res := &hcloud.VolumeCreateOpts{
		Size: in.Size,
	}

	if in.Location != nil {
		location, _, err := c.Location.GetByName(ctx, *in.Location)
		if err != nil {
			return nil, err
		}
		res.Location = location
		// fmt.Println("location:", location)
	}

	if in.Server != nil {
		server, _, err := c.Server.GetByID(ctx, *(in.Server))
		if err != nil {
			return nil, err
		}
		res.Server = server
		// fmt.Println("server:", server)
	}

	if in.Automount != nil {
		res.Automount = in.Automount
		// fmt.Println("server:", server)
	}

	if in.Format != nil {
		res.Format = in.Format
		// fmt.Println("server:", server)
	}

	if in.Labels != nil {
		labels := *in.Labels
		res.Labels = labels
		// fmt.Println("labels:", labels)
	}

	return res, nil

}

func IsVolumeUpToDate(in *v1alpha1.VolumeParameters, o *hcloud.Server) (bool, error) {

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
