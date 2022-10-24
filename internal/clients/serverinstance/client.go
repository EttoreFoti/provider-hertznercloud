package serverinstance

import (
	"context"

	"github.com/hetznercloud/hcloud-go/hcloud"
)

type Client interface {
	GetByID(ctx context.Context, id int) (*hcloud.Server, *hcloud.Response, error)
	GetByName(ctx context.Context, name string) (*hcloud.Server, *hcloud.Response, error)
	Get(ctx context.Context, idOrName string) (*hcloud.Server, *hcloud.Response, error)
	Create(ctx context.Context, opts hcloud.ServerCreateOpts) (hcloud.ServerCreateResult, *hcloud.Response, error)
	Delete(ctx context.Context, server *hcloud.Server) (*hcloud.Response, error)
	Update(ctx context.Context, server *hcloud.Server, opts hcloud.ServerUpdateOpts) (*hcloud.Server, *hcloud.Response, error)
}

func NewClient(creds []byte) Client {
	return &hcloud.NewClient(hcloud.WithToken(string(creds))).Server
}
