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

// BrokerList list all the brokers
func BrokerList(results int) error {
	c, err := client.New()
	if err != nil {
		return err
	}

	fmt.Print("Listing Broker... ")
	quit := progress()
	quit <- true
	<-quit

	if results == defaultLimit {
		results = c.ResponseLimit
	}

	brokers, count, err := brokers.List(c, results)
	if err != nil {
		return err
	}
	fmt.Printf("=== Apps%s", limitCount(len(brokers), count))

	for _, broker := range brokers {
		fmt.Printf("%s %s %s\n", broker.Name, broker.Username, broker.URL)
	}
	return nil
}
