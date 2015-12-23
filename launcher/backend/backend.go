package backend
import (
"sync"
"io"
	"github.com/deis/deis/launcher/model"
)

type Backend interface {
	Destroy(model.Container, *sync.WaitGroup, io.Writer, io.Writer)
	Start(model.Container, *sync.WaitGroup, io.Writer, io.Writer)
	Stop(model.Container, *sync.WaitGroup, io.Writer, io.Writer)
	Scale(model.Container, int, *sync.WaitGroup, io.Writer, io.Writer)
}
