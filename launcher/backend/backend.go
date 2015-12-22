package backend
import (
"sync"
"io"
)

type Backend interface {
	Destroy([]string, *sync.WaitGroup, io.Writer, io.Writer)
	Start([]string, *sync.WaitGroup, io.Writer, io.Writer)
	Stop([]string, *sync.WaitGroup, io.Writer, io.Writer)
	Scale(string, int, *sync.WaitGroup, io.Writer, io.Writer)
}
