package fleet

import (
	"github.com/coreos/fleet/client"
	"sync"
	"io"
	"net/url"
	"net/http"
)
type FleetClient struct {
	fleet         client.API
}

func NewClient(ep url.URL) (*FleetClient, error){
	client, err := client.NewHTTPClient(http.DefaultClient, ep)
	if err != nil {
		return nil, err
	}
	return &FleetClient{
		fleet: client,
	}, nil
}

func (*FleetClient)Destroy([]string, *sync.WaitGroup, io.Writer, io.Writer) {

}

func (*FleetClient) Start([]string, *sync.WaitGroup, io.Writer, io.Writer) {

}

func (*FleetClient) Stop([]string, *sync.WaitGroup, io.Writer, io.Writer) {

}

func (*FleetClient) Scale(string, int, *sync.WaitGroup, io.Writer, io.Writer) {

}
