package parser
import (
	"fmt"
	"os"
	"github.com/docopt/docopt-go"
	"github.com/cde/version"
	"github.com/deis/deis/launcher/cmd"
)


func Start(argv []string) error {
	usage := `
Start a stack
Usage:
	lanncher start <file> <config-url> <backend>

<file>					the stack definiation file
<config-url> 			the service discovery url (consul://localhost:4000 | etcd://localhost:4000)
<backend> 				the backend (fleet://localhost:4000)
`
	args, err := docopt.Parse(usage, argv, true, version.Version, false, true)

	if err != nil {
		return err
	}

	return cmd.Start(args["<file>"].(string), args["<config-url>"].(string), args["<backend>"].(string))
}

func Stop(argv []string) error {
	usage := `
Stop a stack
Usage:
	lanncher stop <stackid>

<stackid>					the stack runtime id to stop
`
	args, err := docopt.Parse(usage, argv, true, version.Version, false, true)

	if err != nil {
		return err
	}

	return cmd.Stop(args["<stackid>"].(string))
}

func PrintUsage() {
	fmt.Fprintln(os.Stderr, "Found no matching command, try 'launcher help'")
	fmt.Fprintln(os.Stderr, "Usage: launcher <command> [<args>...]")
}
