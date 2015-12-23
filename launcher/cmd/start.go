package cmd
import (
	"os"
	"path/filepath"
	"fmt"
	"github.com/deis/deis/launcher/config/etcd"
	"net/url"
	"github.com/deis/deis/launcher/backend/fleet"
	"gopkg.in/yaml.v2"
	"io/ioutil"
"github.com/deis/deis/launcher/model"
	"sync"
)



func Start(filename string, configURL *url.URL, backendURL *url.URL) error {
	manifest := filepath.Join(filepath.Dir(os.Args[0]), filename)
	if _, err := os.Stat(manifest); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "%s Not exists", manifest)
		return err
	}

	configURL.Scheme = "http"
	_, err := config.NewEtcdBackend(*configURL)
	if err != nil {
		return err
	}
	backendURL.Scheme = "http"
	backend, err := fleet.NewClient(*backendURL)
	if err != nil {
		return err
	}

	content, err := ioutil.ReadFile(manifest)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can not read file %s", manifest)
		return err
	}
	m := make(map[string]map[string]interface{})

	err = yaml.Unmarshal(content, &m)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid format %s", manifest)
	}

	fmt.Println(m)
	containers := make([]model.Container, 0);
	for key, value := range m {
		containers = append(containers, model.Container{
			Name: key,
			Desc: model.ContainerDesc{
				Image: value["image"].(string),
				Links: tranform(value["links"]),
				Ports: tranform(value["ports"]),
			},
		})
	}

	for _, item := range containers {
		fmt.Println(item)
		wg := &sync.WaitGroup{}
		backend.Start(item, wg, os.Stdout, os.Stderr)
		wg.Wait()
	}

	return nil
}

func tranform(src interface{}) []string {
	dst := make([]string, 0)
	if src == nil {
		return dst
	}
	srcSlice := src.([]interface{})
	for _, item := range srcSlice {
		dst = append(dst, item.(string))
	}
	return dst
}
