package cmd

import (
	"fmt"
	"github.com/deis/deis/client/controller/client"
	"github.com/deis/deis/client/controller/models/marketplace"
)

// MarketplaceList list all the available service and plans
func MarketplaceList() error {
	c, err := client.New()
	if err != nil {
		return err
	}

	fmt.Print("Listing marketplace... ")
	quit := progress()
	quit <- true
	<-quit

	_, err = marketplace.List(c)
	if err != nil {
		return err
	}
	return nil
}

// MarketplaceInfo get detail of the service
func MarketplaceInfo(service string) error {
	c, err := client.New()
	if err != nil {
		return err
	}

	fmt.Print("Listing marketplace... ")
	quit := progress()
	quit <- true
	<-quit

	_, err = marketplace.List(c)
	if err != nil {
		return err
	}
	return nil
}
