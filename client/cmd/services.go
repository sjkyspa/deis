package cmd

import (
	"fmt"
	"github.com/deis/deis/client/controller/client"
	"github.com/deis/deis/client/controller/models/services"
	"os"
)

// ServiceList lists services
func ServiceList(results int) error {
	c, err := client.New()

	if err != nil {
		return err
	}

	if results == defaultLimit {
		results = c.ResponseLimit
	}

	services, count, err := services.List(c, results)

	if err != nil {
		return err
	}

	fmt.Printf("=== Apps%s", limitCount(len(services), count))

	for _, app := range services {
		fmt.Println(app.ID)
	}
	return nil
}

// ServiceCreate creates an service.
func ServiceCreate(serviceName, planName, serviceInstanceName string) error {
	c, err := client.New()

	if err != nil {
		return err
	}

	serviceInstance, err := services.New(c, serviceName, planName, serviceInstanceName)
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
func ServiceDelete() error {
	return nil
}

// ServiceRename rename the service
func ServiceRename() error {
	return nil
}

// ServiceBind bind the service to the app
func ServiceBind() error {
	return nil
}

// ServiceUnbind unbind the service from the app
func ServiceUnbind() error {
	return nil
}
