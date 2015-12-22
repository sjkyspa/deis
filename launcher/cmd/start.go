package cmd
import (
	"os"
	"path/filepath"
	"fmt"
	"github.com/deis/deis/launcher/config/etcd"
	"net/url"
	"github.com/deis/deis/launcher/backend/fleet"
)

func Start(filename string, configUrl *url.URL, backendUrl *url.URL) error {
	if _, err :=os.Stat(filepath.Join(filepath.Dir(os.Args[0]), filename)); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "%s Not exists", filepath.Join(filepath.Dir(os.Args[0]), filename))
		return err
	}

	configUrl.Scheme = "http"
	eb, err := config.NewEtcdBackend(*configUrl)
	if err != nil {
		return err
	}
	backendUrl.Scheme = "http"
	backend, err := fleet.NewClient(*backendUrl)
	if err != nil {
		return err
	}
	eb.Get("key")
	backend.Start(nil, nil, nil, nil)
	return nil
}
