package fleet
import (
	"github.com/deis/deis/launcher/model"
)

func (*FleetClient) Destroy(model.Container) error {
	return nil
}

