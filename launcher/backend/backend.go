package backend
import (
	"github.com/deis/deis/launcher/model"
)

type Backend interface {
	Destroy(model.Container) error
	Start(model.Container) error
	Stop(model.Container) error
	Scale(model.Container) error
}
