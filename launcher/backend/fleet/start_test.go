package fleet
import (
	"testing"
	"net/url"
	"github.com/deis/deis/launcher/model"
	. "github.com/onsi/gomega"
	"net/http"
)


func TestCreateUnit(t *testing.T) {
	t.Parallel()
	RegisterTestingT(t)

	fleetURL, err := url.Parse("http://localhost:4001")
	if err != nil {
		t.Fatal(err)
		return
	}
	c, err := NewClient(http.DefaultClient, *fleetURL)
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

	Expect(unit.Name).To(Equal(container.Name))
	Expect(len(unit.Options)).To(Equal(3), "Expected Options size is 3 (ExecStartPre, ExecStart, ExecStartPost)")
}
