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

// List broker installed.
func List(c *client.Client, count int) ([]api.Broker, int, error) {

	var err error
	resBody, count, err := c.LimitedRequest("/v1/brokers", count)

	if err != nil {
		return []api.Broker{}, -1, err
	}

	var brokers []api.Broker
	if err = json.Unmarshal([]byte(resBody), &brokers); err != nil {
		return []api.Broker{}, -1, err
	}

	return brokers, count, nil
}
