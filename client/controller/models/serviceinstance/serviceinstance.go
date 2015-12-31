package serviceinstance

import (
	"encoding/json"
	"fmt"
	"github.com/deis/deis/client/controller/api"
	"github.com/deis/deis/client/controller/client"
)

// FindByName find the service instance by the service name
func FindByName(c *client.Client, name string) (api.ServiceInstance, error) {
	body, err := c.BasicRequest("GET", fmt.Sprintf("/v1/service_instances?name=%s", name), nil)
	if err != nil {
		return api.ServiceInstance{}, err
	}

	res, count, err := extractResult(body)

	if count > 1 || count <= 0 {
		return api.ServiceInstance{}, fmt.Errorf("The service is name is not unique or service not found")
	}

	var serviceInstance []api.ServiceInstance

	err = json.Unmarshal([]byte(res), &serviceInstance)
	if err != nil {
		return api.ServiceInstance{}, err
	}

	return serviceInstance[0], nil
}

func extractResult(body string) (string, int, error) {
	res := make(map[string]interface{})
	if err := json.Unmarshal([]byte(body), &res); err != nil {
		return "", -1, err
	}

	out, err := json.Marshal(res["results"].([]interface{}))

	if err != nil {
		return "", -1, err
	}

	return string(out), int(res["count"].(float64)), nil
}
