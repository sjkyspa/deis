package fleet

import (
	"io"
	"github.com/deis/deis/launcher/model"
	"github.com/coreos/fleet/schema"
	"fmt"
	"github.com/deis/deis/pkg/prettyprint"
	"time"
	"os"
)

func (f *FleetClient) Start(container model.Container) error {
	unit, err := f.createUnit(container)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when create unit %s", container)
		return err
	}

	f.Fleet.CreateUnit(unit)
	var name string = unit.Name

	lastSubState := "dead"
	requestState := "launched"
	desiredState := "running"
	err = f.Fleet.SetUnitTargetState(unit.Name, requestState)
	if err != nil {
		io.WriteString(os.Stderr, err.Error())
		return err
	}

	for {
		// poll for unit states
		states, err := f.Fleet.UnitStates()
		if err != nil {
			io.WriteString(os.Stderr, err.Error())
			return err
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
			fmt.Fprintf(os.Stderr, "Could not find unit: %v\n", name)
			return fmt.Errorf("Could not find unit: %v\n", name)
		}

		// if subState changed, send it across the output channel
		if lastSubState != currentState.SystemdSubState {
			l := prettyprint.Overwritef(prettyprint.Colorize("{{.Yellow}}%v:{{.Default}} %v/%v"), name, currentState.SystemdActiveState, currentState.SystemdSubState)
			fmt.Fprintf(os.Stdout, l)
		}

		// break when desired state is reached
		if currentState.SystemdSubState == desiredState {
			fmt.Fprintln(os.Stdout)
			return nil
		}

		lastSubState = currentState.SystemdSubState

		if lastSubState == "failed" {
			return fmt.Errorf("The service '%s' failed while starting", name)

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
