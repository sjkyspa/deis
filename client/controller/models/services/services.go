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

// New create service on a Deis controller.
func New(c *client.Client, serviceName, planName string) (api.ServiceOffering, error) {
	body := []byte{}
	var err error
	req := api.ServiceInstanceCreateRequest{ServiceName: serviceName, PlanName: planName}
	body, err = json.Marshal(req)

	if err != nil {
		return api.ServiceOffering{}, err
	}

	res, err := c.BasicRequest("POST", "/v1/service-instances/", body)

	if err != nil {
		return api.ServiceOffering{}, err
	}

	var services api.ServiceOffering
	if err = json.Unmarshal([]byte(res), &services); err != nil {
		return api.ServiceOffering{}, err
	}

	return services, nil
}
