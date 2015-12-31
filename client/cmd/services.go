package cmd

import (
	"fmt"
	"github.com/deis/deis/client/controller/client"
	"github.com/deis/deis/client/controller/models/servicebinding"
	"github.com/deis/deis/client/controller/models/serviceinstance"
	"github.com/deis/deis/client/controller/models/services"
	"os"
	"strings"
	"text/tabwriter"
)

// ServiceList lists services
func ServiceList(c *client.Client, results int) error {
	if results == defaultLimit {
		results = c.ResponseLimit
	}

	services, count, err := services.List(c, results)

	if err != nil {
		return err
	}

	fmt.Printf("=== Services%s", limitCount(len(services), count))

	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 16, 8, 2, ' ', 0)
	fmt.Fprintln(w, fmt.Sprintf("%s\t%s", "Service Name", "Service Plans"))
	for _, service := range services {
		var plans []string
		for _, plan := range service.Plans {
			plans = append(plans, plan.Name)
		}

		fmt.Fprintln(w, fmt.Sprintf("%s\t%s", service.NAME, strings.Join(plans, ",")))
	}
	w.Flush()
	return nil
}

// ServiceCreate creates an service.
func ServiceCreate(c *client.Client, serviceName, planName, serviceInstanceName string) error {
	service, err := services.FindByName(c, serviceName)
	if err != nil {
		return err
	}

	plan, err := service.FindPlan(planName)
	if err != nil {
		return err
	}

	serviceInstance, err := services.New(c, serviceInstanceName, plan.ID)
	if err != nil {
		return err
	}

	fmt.Fprintf(os.Stdout, "Instance created %s", serviceInstance.Name)

	return nil
}

// ServiceUpdate update a service
func ServiceUpdate() error {
	return nil
}

// ServiceDelete delete a service
func ServiceDelete(c *client.Client, serviceInstanceName string) error {
	serviceInstance, err := serviceinstance.FindByName(c, serviceInstanceName)
	if err != nil {
		return err
	}

	return serviceinstance.Deprovision(c, serviceInstance.ID)
}

// ServiceRename rename the service
func ServiceRename() error {
	return nil
}

// ServiceBind bind the service to the app
func ServiceBind(c *client.Client, appName, serviceInstanceName string) error {
	serviceInstance, err := serviceinstance.FindByName(c, serviceInstanceName)
	if err != nil {
		return err
	}
	err = servicebinding.Bind(c, serviceInstance.ID, appName, nil)
	if err != nil {
		return err
	}

	fmt.Fprint(os.Stdout, "Binding created")

	ConfigSet(appName, []string{"DATABASE=test"})

	return nil
}

// ServiceUnbind unbind the service from the app
func ServiceUnbind(c *client.Client, appID, serviceInstanceName string) error {
	serviceInstance, err := serviceinstance.FindByName(c, serviceInstanceName)

	err = servicebinding.Unbind(c, serviceInstance.ID, appID)
	if err != nil {
		return err
	}
	return nil
}
