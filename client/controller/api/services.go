package api

import "fmt"

// ServiceOfferingFields is the definition of the service meta.
type ServiceOfferingFields struct {
	ID               string `json:"uuid"`
	NAME             string `json:"name"`
	BrokerID         string `json:"broker_id"`
	Label            string `json:"label"`
	Provider         string `json:"provider"`
	Version          string `json:"version"`
	Description      string `json:"description"`
	DocumentationURL string `json:"doc_url"`
}

// ServiceOffering is the definition of the service offering by the broker
type ServiceOffering struct {
	ServiceOfferingFields
	Plans []ServicePlanFields `json:"plans"`
}

// FindPlan find plan in the service
func (so *ServiceOffering) FindPlan(planName string) (ServicePlanFields, error) {
	for _, plan := range so.Plans {
		if planName == plan.Name {
			return plan, nil
		}
	}

	return ServicePlanFields{}, fmt.Errorf("%s is not a valid plan for serivce %s", planName, so.ServiceOfferingFields.NAME)
}

// ServicePlanFields is the definition of the service plan meta
type ServicePlanFields struct {
	ID                string `json:"uuid"`
	Name              string `json:"name"`
	Free              bool   `json:"free"`
	Public            bool   `json:"public"`
	Description       string `json:"description"`
	Active            bool   `json:"active"`
	ServiceOfferingID string `json:"guid"`
}

// ServicePlan is the definition of the service plan the broker provided
type ServicePlan struct {
	ServicePlanFields
	ServiceOffering ServiceOfferingFields
}

// ServiceInstanceCreateRequest is the request to create service instance
type ServiceInstanceCreateRequest struct {
	Name   string `json:"name"`
	PlanID string `json:"plan_uuid"`
}

// ServiceInstanceFields is the definition of the service instance
type ServiceInstanceFields struct {
	ID               string
	Name             string
	SysLogDrainURL   string
	ApplicationNames []string
	Params           map[string]interface{}
	DashboardURL     string
}

// ServiceInstance is the definition the service instance
type ServiceInstance struct {
	ServiceInstanceFields
	ServiceBindings []ServiceBindingFields
	ServicePlan     ServicePlanFields
	ServiceOffering ServiceOfferingFields
}
