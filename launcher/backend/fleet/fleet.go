package fleet

import (
	"github.com/coreos/fleet/client"
	"sync"
	"io"
	"net/url"
	"net/http"
	"github.com/deis/deis/launcher/model"
	"github.com/coreos/fleet/schema"
	"fmt"
	"github.com/deis/deis/pkg/prettyprint"
	"time"
)
type FleetClient struct {
	Fleet client.API
}

func NewClient(ep url.URL) (*FleetClient, error) {
	client, err := client.NewHTTPClient(http.DefaultClient, ep)
	if err != nil {
		return nil, err
	}
	return &FleetClient{
		Fleet: client,
	}, nil
}

func (*FleetClient) Destroy(model.Container, *sync.WaitGroup, io.Writer, io.Writer) {

}

func (f *FleetClient) Start(container model.Container, wg *sync.WaitGroup, out, ew io.Writer) {
	defer wg.Done()
	unit, err := f.createUnit(container)
	if err != nil {
		fmt.Fprintf(ew, "Error when create unit %s", container)
		return
	}

	f.Fleet.CreateUnit(unit)
	var name string = unit.Name

	lastSubState := "dead"
	requestState := "launched"
	desiredState := "running"
	err = f.Fleet.SetUnitTargetState(unit.Name, requestState)
	if err != nil {
		io.WriteString(ew, err.Error())
		return
	}

	for {
		// poll for unit states
		states, err := f.Fleet.UnitStates()
		if err != nil {
			io.WriteString(ew, err.Error())
			return
		}

		// FIXME: fleet UnitStates API forces us to iterate for now
		var currentState *schema.UnitState
		for _, s := range states {
			if name == s.Name {
				currentState = s
				break
			}
		}
		if currentState == nil {
			fmt.Fprintf(ew, "Could not find unit: %v\n", name)
			return
		}

		// if subState changed, send it across the output channel
		if lastSubState != currentState.SystemdSubState {
			l := prettyprint.Overwritef(prettyprint.Colorize("{{.Yellow}}%v:{{.Default}} %v/%v"), name, currentState.SystemdActiveState, currentState.SystemdSubState)
			fmt.Fprintf(out, l)
		}

		// break when desired state is reached
		if currentState.SystemdSubState == desiredState {
			fmt.Fprintln(out)
			return
		}

		lastSubState = currentState.SystemdSubState

		if lastSubState == "failed" {
			o := prettyprint.Colorize("{{.Red}}The service '%s' failed while starting.{{.Default}}\n")
			fmt.Fprintf(ew, o, name)
			return
		}
		time.Sleep(250 * time.Millisecond)
	}
}

func (f *FleetClient) createUnit(container model.Container) (*schema.Unit, error) {
	options := make([]*schema.UnitOption, 0)
	options = append(options, &schema.UnitOption{
		Section: "Service",
		Name: "ExecStartPre",
		Value: fmt.Sprintf("docker pull %s", container.Desc.Image),
	})
	options = append(options, &schema.UnitOption{
		Section: "Service",
		Name: "ExecStart",
		Value: fmt.Sprintf("docker run --rm --name=%s %s ", container.Name, container.Desc.Image),
	})
	options = append(options, &schema.UnitOption{
		Section: "Service",
		Name: "ExecStop",
		Value: fmt.Sprintf("docker stop %s", container.Name),
	})
	unit := schema.Unit{
		DesiredState: "launched",
		Name: container.Name,
		Options: options,
	}

	return &unit, nil
}

func (*FleetClient) Stop(model.Container, *sync.WaitGroup, io.Writer, io.Writer) {

}

func (*FleetClient) Scale(model.Container, int, *sync.WaitGroup, io.Writer, io.Writer) {

}
