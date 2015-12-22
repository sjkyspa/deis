package fleet

import (
	"github.com/coreos/fleet/client"
	"sync"
	"io"
"github.com/deis/deis/launcher/config"
)
type FleetClient struct {
	Fleet         client.API
	config config.Backend
}

func NewClient() {

}

func (*FleetClient)Destroy([]string, *sync.WaitGroup, io.Writer, io.Writer) {

}

func (*FleetClient) Start([]string, *sync.WaitGroup, io.Writer, io.Writer) {

}

func (*FleetClient) Stop([]string, *sync.WaitGroup, io.Writer, io.Writer) {

}

func (*FleetClient) Scale(string, int, *sync.WaitGroup, io.Writer, io.Writer) {

}
