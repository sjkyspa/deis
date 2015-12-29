package servicebinding

import (
	"encoding/json"
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
