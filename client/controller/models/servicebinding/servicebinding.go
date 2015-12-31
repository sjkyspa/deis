package servicebinding

import (
	"encoding/json"
	"fmt"
	"github.com/deis/deis/client/controller/api"
	"github.com/deis/deis/client/controller/client"
)

// Bind bind the service instance to the app
func Bind(c *client.Client, serviceInstanceID, appID string, params map[string]interface{}) error {
	body := []byte{}

	bindingCreateRequest := api.ServiceBindingCreateRequest{
		ServiceInstanceID: serviceInstanceID,
		AppID:             appID,
		Params:            params,
	}

	body, err := json.Marshal(bindingCreateRequest)
	if err != nil {
		return err
	}

	_, err = c.BasicRequest("POST", "/v1/service_bindings", body)
	if err != nil {
		return err
	}

	return nil
}

// Unbind remove the service binding
func Unbind(c *client.Client, serviceInstanceID, appID string) error {
	sb, err := findServiceBinding(c, serviceInstanceID, appID)
	if err != nil {
		return err
	}

	_, err = c.BasicRequest("DELETE", fmt.Sprintf("/v1/service_bindings/%s", sb.ID), nil)
	if err != nil {
		return err
	}

	return nil
}

func findServiceBinding(c *client.Client, serviceInstanceID, appID string) (api.ServiceBindingFields, error) {
	var findServiceBindingInner func(path string) (api.ServiceBindingFields, error)
	findServiceBindingInner = func(path string) (api.ServiceBindingFields, error) {
		body, err := c.BasicRequest("GET", path, nil)
		if err != nil {
			return api.ServiceBindingFields{}, err
		}

		res := make(map[string]interface{})

		err = json.Unmarshal([]byte(body), &res)
		if err != nil {
			return api.ServiceBindingFields{}, err
		}

		if int(res["count"].(float64)) <= 0 {
			return api.ServiceBindingFields{}, fmt.Errorf("Can not find service by service instance: %s and appID: %s", serviceInstanceID, appID)
		}

		var serviceBindings []api.ServiceBindingFields
		results, err := json.Marshal(res["results"].([]interface{}))
		if err != nil {
			return api.ServiceBindingFields{}, err
		}
		err = json.Unmarshal([]byte(results), &serviceBindings)
		if err != nil {
			return api.ServiceBindingFields{}, err
		}

		for _, sb := range serviceBindings {
			if sb.AppID == appID && sb.ServiceInstanceID == serviceInstanceID {
				return sb, nil
			}
		}

		next := res["next"]
		if next == nil {
			return api.ServiceBindingFields{}, fmt.Errorf("Can Not find service by service instance: %s and appID: %s", serviceInstanceID, appID)
		}

		return findServiceBindingInner(next.(string))
	}

	return findServiceBindingInner("/v1/service_bindings")
}
