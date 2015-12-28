package services

import (
	"encoding/json"
	"github.com/deis/deis/client/controller/api"
	"github.com/deis/deis/client/controller/client"
)

// List lists services on a Deis controller.
func List(c *client.Client, results int) ([]api.ServiceOffering, int, error) {
	body, count, err := c.LimitedRequest("/v1/services/", results)

	if err != nil {
		return []api.ServiceOffering{}, -1, err
	}

	var services []api.ServiceOffering
	if err = json.Unmarshal([]byte(body), &services); err != nil {
		return []api.ServiceOffering{}, -1, err
	}

	return services, count, nil
}
