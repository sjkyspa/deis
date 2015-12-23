package fleet

import (
	"github.com/coreos/fleet/client"
	"net/url"
	"net/http"
)
type FleetClient struct {
	Fleet client.API
}

func NewClient(httpClient *http.Client, ep url.URL) (*FleetClient, error) {
	client, err := client.NewHTTPClient(httpClient, ep)
	if err != nil {
		return nil, err
	}
	return &FleetClient{
		Fleet: client,
	}, nil
}
