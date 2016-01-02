package parser

import (
	"encoding/json"
	"fmt"
	"github.com/deis/deis/client/cmd"
	"github.com/deis/deis/client/controller/client"
	docopt "github.com/docopt/docopt-go"
	"io/ioutil"
	"strings"
)

// Services routes app commands to their specific function.
func Services(argv []string) error {
	usage := `
Valid commands for apps:

services:list          destroy an application
services:create        create a new application
services:update        list accessible applications
services:delete        view info about an application
services:rename        open the application in a browser
services:bind          view aggregated application logs
services:unbind        run a command in an ephemeral app container

Use 'deis help [command]' to learn more.
`

	switch argv[0] {
	case "services:list":
		return servicesList(argv)
	case "services:create":
		return serviceCreate(argv)
	case "services:update":
		return serviceUpdate(argv)
	case "services:delete":
		return serviceDelete(argv)
	case "services:rename":
		return serviceRename(argv)
	case "services:bind":
		return serviceBind(argv)
	case "services:unbind":
		return serviceUnbind(argv)
	default:
		if printHelp(argv, usage) {
			return nil
		}

		if argv[0] == "services" {
			argv[0] = "services:list"
			return servicesList(argv)
		}

		PrintUsage()
		return nil
	}
}

func servicesList(argv []string) error {
	usage := `
Lists services visible to the current user.

Usage: deis services:list [options]

Options:
  -l --limit=<num>
    the maximum number of results to display, defaults to config setting
`
	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	results, err := responseLimit(safeGetValue(args, "--limit"))

	if err != nil {
		return err
	}

	c, err := client.New()

	if err != nil {
		return err
	}

	return cmd.ServiceList(c, results)
}

func serviceCreate(argv []string) error {
	usage := `
Create service instance visible to the current user.

Usage: deis services:create <service-name> <plan-name> <service-instance-name> [options]

Options:
  -c, --config the config for the service when instantiation
`
	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	serviceName := safeGetValue(args, "<service-name>")
	planName := safeGetValue(args, "<plan-name>")
	serviceInstanceName := safeGetValue(args, "<service-instance-name>")
	var configJSON map[string]interface{}
	config := safeGetValue(args, "--config")

	if config != "" {
		if strings.HasPrefix(config, "@") {
			configFile := strings.TrimLeft(config, "@")
			content, err := ioutil.ReadFile(configFile)
			if err != nil {
				return fmt.Errorf("Please check the file realily exists")
			}

			err = json.Unmarshal(content, &configJSON)
			if err != nil {
				return fmt.Errorf("Please ensure the json in config file is valid")
			}
		} else {
			err = json.Unmarshal([]byte(config), &configJSON)

			if err != nil {
				return fmt.Errorf("Pleases provide the valid json as the additional parameters")
			}
		}
	}

	c, err := client.New()

	if err != nil {
		return err
	}

	return cmd.ServiceCreate(c, serviceName, planName, serviceInstanceName, configJSON)
}

func serviceUpdate(argv []string) error {
	return nil
}

func serviceDelete(argv []string) error {
	return nil
}

func serviceRename(argv []string) error {
	return nil
}

func serviceBind(argv []string) error {
	usage := `
Create service instance visible to the current user.

Usage: deis services:bind <app-name> <service-instance-name> [-c <config>]

Options:
  -c
    the config for the service when instantiation
`
	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}
	appName := safeGetValue(args, "<app-name>")
	serviceInstanceName := safeGetValue(args, "<service-instance-name>")

	c, err := client.New()

	if err != nil {
		return err
	}
	return cmd.ServiceBind(c, appName, serviceInstanceName)
}

func serviceUnbind(argv []string) error {
	return nil
}
