package fleet
import (
	"testing"
	"net/url"
	"net/http"
	"sync"
	"os"
	"fmt"
	"github.com/coreos/fleet/schema"
	"github.com/deis/deis/launcher/model"
	. "github.com/onsi/gomega"
)


func TestCreateUnit(t *testing.T) {
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


func TestStart(t *testing.T) {
	RegisterTestingT(t)

	testFleetClient := stubFleetClient{testUnits: []*schema.Unit{}, unitsMutex: &sync.Mutex{},
		unitStatesMutex: &sync.Mutex{}}

	c := &FleetClient{Fleet: &testFleetClient}

	container := model.Container{
		Name: "mysql",
		Desc: model.ContainerDesc{
			Image: "mysql",
			Ports: []string{"80:80"},
		},
	}

	var errOutput string
	var wg sync.WaitGroup

	c.Start(container, &wg, os.Stdout, os.Stderr)
	wg.Wait()

	logMutex := sync.Mutex{}
	logMutex.Lock()
	if errOutput != "" {
		t.Fatal(errOutput)
	}
	logMutex.Unlock()

	var found bool
	for _, unit := range testFleetClient.testUnitStates {
		if unit.Name == container.Name {
			found = true
			Expect(unit.SystemdSubState).To(Equal("running"), fmt.Sprintf("Unit %s is %s, expected running", unit.Name, unit.SystemdSubState))
			break
		}
	}
	if !found {
		t.Fatal("Not found services")
	}
}
