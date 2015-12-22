package parser
import (
	"fmt"
	"os"
)


func Start(argv []string) error {
	fmt.Println("start")
	return nil
}

func Stop(argv []string) error {
	fmt.Println("stop")
	return nil
}

func PrintUsage() {
	fmt.Fprintln(os.Stderr, "Found no matching command, try 'launcher help'")
	fmt.Fprintln(os.Stderr, "Usage: launcher <command> [<args>...]")
}
