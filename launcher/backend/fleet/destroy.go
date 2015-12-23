package fleet
import (
	"github.com/deis/deis/launcher/model"
	"sync"
	"io"
)

func (*FleetClient) Destroy(model.Container, *sync.WaitGroup, io.Writer, io.Writer) {

}

