package fleet
import (
	"testing"
	"net/url"
	"github.com/deis/deis/launcher/model"
)


func TestCreateUnit(t *testing.T) {
	t.Parallel()

	fleetURL, err := url.Parse("http://localhost:4001")
	if err != nil {
		t.Fatal(err)
		return
	}
	c, err := NewClient(*fleetURL)
	container := model.Container{
		Name: "mysql",
		Desc: model.ContainerDesc{
			Image: "mysql",
			Ports: []string{"80:80"},
		},
	}
	unit, err := c.createUnit(container)
	if err != nil {
		t.Fatal(err)
	}
	if (unit.DesiredState != "launched") {
		t.Error("The status should be launched")
	}
}
