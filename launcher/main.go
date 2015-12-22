package main
import (
	"os"
	"github.com/docopt/docopt-go"
	"github.com/deis/deis/version"
//	"github.com/deis/deis/launcher/parser"
	"fmt"
	"github.com/deis/deis/launcher/parser"
)


func main() {
	os.Exit(Command(os.Args[1:]))
}

func Command(argv[]string) int {
	usage := `
The launcher command-line client.
Usage:
	launcher <command> [<args>...]
	launcher -h | --help

start 		launch a stack
destroy 	destroy a stack
stop 		a stack
`
	args, err := docopt.Parse(usage, argv, false, version.Version, true, false)
	if err != nil {
		fmt.Println(args)
		return 1
	}


	switch args["<command>"] {
	case "start":
		err = parser.Start(append([]string{args["<command>"].(string)}, (args["<args>"].([]string))...))
	case "stop":
		err = parser.Stop(append([]string{args["<command>"].(string)}, (args["<args>"].([]string))...))
	case "help":
		parser.PrintUsage()
	default:
		parser.PrintUsage()
	}

	if err != nil {
		return 1
	}

	return 0
}


