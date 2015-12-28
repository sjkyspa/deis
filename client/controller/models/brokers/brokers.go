package brokers

import (
	"encoding/json"
	"github.com/deis/deis/client/controller/api"
	"github.com/deis/deis/client/controller/client"
	"net/url"
)

// New creates a new broker.
func New(c *client.Client, brokerName, username, password string, url url.URL) (api.Broker, error) {
	body := []byte{}

	var err error
	req := api.BrokerCreateRequest{
		Name:     brokerName,
		Username: username,
		Password: password,
		URL:      url.String(),
	}
	body, err = json.Marshal(req)

	if err != nil {
		return api.Broker{}, err
	}

	resBody, err := c.BasicRequest("POST", "/v1/brokers/", body)

	if err != nil {
		return api.Broker{}, err
	}

	broker := api.Broker{}
	if err = json.Unmarshal([]byte(resBody), &broker); err != nil {
		return api.Broker{}, err
	}

	return broker, nil
}
