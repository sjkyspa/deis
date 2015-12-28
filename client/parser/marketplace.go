package parser

import (
	"github.com/deis/deis/client/cmd"
	docopt "github.com/docopt/docopt-go"
)

// Marketplace routes broker commands to their specific function.
func Marketplace(argv []string) error {
	usage := `
Valid commands for apps:

marketplace:list             list all the available service and plans
marketplace:info <service>   get the service detail

Use 'deis help [command]' to learn more.
`

	switch argv[0] {
	case "marketplace:list":
		return marketplaceList(argv)
	case "marketplace:info":
		return marketplaceInfo(argv)
	default:
		if printHelp(argv, usage) {
			return nil
		}

		PrintUsage()
		return nil
	}
}

func marketplaceList(argv []string) error {
	usage := `
Creates a new broker.

Usage: deis marketplace:list

`
	_, err := docopt.Parse(usage, argv, true, "", false, true)
	if err != nil {
		return err
	}

	return cmd.MarketplaceList()
}

func marketplaceInfo(argv []string) error {
	usage := `
Creates a new broker.

Usage: deis marketplace [-s <service>]

`
	args, err := docopt.Parse(usage, argv, true, "", false, true)
	if err != nil {
		return err
	}

	service := safeGetValue(args, "<service>")

	return cmd.MarketplaceInfo(service)
}
