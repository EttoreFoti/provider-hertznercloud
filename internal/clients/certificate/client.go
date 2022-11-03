package certificate

import (
	"context"
	"fmt"

	"github.com/crossplane/provider-hertznercloud/apis/certificates/v1alpha1"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/hetznercloud/hcloud-go/hcloud"
)

func FromCertificateSpecToCertificateCreateOpts(in *v1alpha1.CertificateParameters, c *hcloud.Client, ctx context.Context) (*hcloud.CertificateCreateOpts, error) {

	res := &hcloud.CertificateCreateOpts{
		Name: in.Name, // ????
	}

	if in.Type != nil {
		res.Type = hcloud.CertificateType(*in.Type)
	}

	if in.Labels != nil {
		labels := *in.Labels
		res.Labels = labels
	}

	if in.Certificate != nil {
		res.Certificate = *in.Certificate
	}

	if in.DomainNames != nil {
		domainNames := *in.DomainNames
		res.DomainNames = domainNames
	}

	if in.PrivateKey != nil {
		res.PrivateKey = *in.PrivateKey
	}

	return res, nil

}

func IsCertificateUpToDate(in *v1alpha1.CertificateParameters, o *hcloud.Certificate) (bool, error) {

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
