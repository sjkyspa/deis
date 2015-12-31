package services

import (
	"encoding/json"
	"fmt"
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
func New(c *client.Client, serviceInstanceName, planID string) (api.ServiceInstance, error) {
	body := []byte{}
	var err error
	req := api.ServiceInstanceCreateRequest{
		Name:   serviceInstanceName,
		PlanID: planID,
	}
	body, err = json.Marshal(req)

	if err != nil {
		return api.ServiceInstance{}, err
	}

	res, err := c.BasicRequest("POST", "/v1/service_instances", body)

	if err != nil {
		return api.ServiceInstance{}, err
	}

	var serviceInstance api.ServiceInstance
	if err = json.Unmarshal([]byte(res), &serviceInstance); err != nil {
		return api.ServiceInstance{}, err
	}

	return serviceInstance, nil
}

// FindByName find service by the service name
func FindByName(c *client.Client, serviceName string) (api.ServiceOffering, error) {

	res, err := c.BasicRequest("GET", fmt.Sprintf("/v1/services?name=%s", serviceName), nil)

	if err != nil {
		return api.ServiceOffering{}, err
	}

	var service api.ServiceOffering
	if err = json.Unmarshal([]byte(res), &service); err != nil {
		return api.ServiceOffering{}, err
	}

	return service, nil
}
