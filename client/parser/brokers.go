package parser

import (
	"github.com/deis/deis/client/cmd"
	docopt "github.com/docopt/docopt-go"
	"net/url"
)

// Brokers routes broker commands to their specific function.
func Brokers(argv []string) error {
	usage := `
Valid commands for apps:

brokers:create        create a service broker
brokers:list          list all service broker
brokers:info          get the broker info

Use 'deis help [command]' to learn more.
`

	switch argv[0] {
	case "brokers:create":
		return brokersCreate(argv)
	case "brokers:list":
		return brokersList(argv)
	case "brokers:info":
		return brokersInfo(argv)
	default:
		if printHelp(argv, usage) {
			return nil
		}

		if argv[0] == "brokers" {
			argv[0] = "apps:list"
			return brokersList(argv)
		}

		PrintUsage()
		return nil
	}
}

func brokersCreate(argv []string) error {
	usage := `
Creates a new broker.

Usage: deis brokers:create <name> <username> <password> <url>
`
	args, err := docopt.Parse(usage, argv, true, "", false, true)
	if err != nil {
		return err
	}
	name := safeGetValue(args, "<name>")
	username := safeGetValue(args, "<username>")
	password := safeGetValue(args, "<password>")

	brokerURL, err := url.Parse(safeGetValue(args, "<url>"))
	if err != nil {
		return err
	}

	return cmd.BrokerCreate(name, username, password, *brokerURL)
}

func brokersList(argv []string) error {
	usage := `
Lists brokders visible to the current user.

Usage: deis brokers:list [options]

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

	return cmd.BrokerList(results)
}

func brokersInfo(argv []string) error {
	return nil
}
