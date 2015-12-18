package parser

import (
	"github.com/deis/deis/client/cmd"
	"github.com/docopt/docopt-go"
)

func Stacks(argv []string) error {
	switch argv[0] {
	case "stacks:init":
		return stackCreate(argv)
	default:
		PrintUsage()
		return nil
	}

}

func stackCreate(argv []string) error {
	usage := `
Init the app with the stack template

Usage: deis stacks:init [options]

Options:
  -a --app=<app>
    the uniquely identifiable name of the application.
  -s --stack=<stack>
    the stack in the stack repo.
`

	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}
	return cmd.StackCreate(safeGetValue(args, "--app"), safeGetValue(args, "--stack"))
}


