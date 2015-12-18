package parser

import (
	"github.com/deis/deis/client/cmd"
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
	return cmd.StackCreate(argv[1])
}
