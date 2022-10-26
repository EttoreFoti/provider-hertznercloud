package hertznercloud

import (
	"errors"

	"github.com/hetznercloud/hcloud-go/hcloud"
)

func NewClientHertzner(creds []byte) (*hcloud.Client, error) {
	c := hcloud.NewClient(hcloud.WithToken(string(creds)))

	if c == nil {
		return nil, errors.New("unable to create client")
	}

	return c, nil
}
