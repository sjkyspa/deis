package marketplace

import (
	"encoding/json"
	"fmt"
	"github.com/deis/deis/client/controller/api"
	"github.com/deis/deis/client/controller/client"
)

// List all the available service and plans
func List(c *client.Client) ([]api.ServiceOffering, error) {
	resBody, _, err := c.LimitedRequest("/v1/services", -1)
	if err != nil {
		return []api.ServiceOffering{}, err
	}

	fmt.Println("Service in marketplace ===> marketplace")
	var services []api.ServiceOffering
	if err = json.Unmarshal([]byte(resBody), &services); err != nil {
		return []api.ServiceOffering{}, err
	}

	return services, nil
}

// Info get the detail of the service
func Info(*client.Client) (api.ServiceOffering, error) {
	return api.ServiceOffering{}, nil
}
