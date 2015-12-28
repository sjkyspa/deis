package cmd

import (
	"fmt"
	"github.com/deis/deis/client/controller/client"
	"github.com/deis/deis/client/controller/models/brokers"
	"net/url"
)

// BrokerCreate creates a broker.
func BrokerCreate(brokerName string, username string, password string, url url.URL) error {
	c, err := client.New()
	if err != nil {
		return err
	}

	fmt.Print("Creating Broker... ")
	quit := progress()
	quit <- true
	<-quit

	_, err = brokers.New(c, brokerName, username, password, url)
	if err != nil {
		return err
	}
	return nil
}
