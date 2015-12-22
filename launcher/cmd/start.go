package cmd
import (
	"os"
	"path/filepath"
	"fmt"
)

func Start(filename, serviceDiscovery, backend string) error {
	if _, err :=os.Stat(filepath.Join(filepath.Dir(os.Args[0]), filename)); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "%s Not exists", filepath.Join(filepath.Dir(os.Args[0]), filename))
		return err
	}


	return nil
}
