package hertznercloud

import (
	"errors"

	"github.com/hetznercloud/hcloud-go/hcloud"
)

type Client struct {
	HertznerClient *hcloud.Client
}

func NewClientHertzner(creds []byte) (*Client, error) {
	c := hcloud.NewClient(hcloud.WithToken(string(creds)))

	if c == nil {
		return nil, errors.New("unable to create client")
	}

	return &Client{
		HertznerClient: c,
	}, nil
}
