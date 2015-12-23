package fleet

import (
	"github.com/deis/deis/launcher/model"
	"sync"
	"io"
)

func (*FleetClient) Scale(model.Container, int, *sync.WaitGroup, io.Writer, io.Writer) {

}
