package fleet

import (
	"github.com/deis/deis/launcher/model"
)

func (*FleetClient) Scale(model.Container) error {
	return nil
}
