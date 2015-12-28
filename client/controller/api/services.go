package api

// ServiceOfferingFields is the definition of the service meta.
type ServiceOfferingFields struct {
	ID               string `json:"id"`
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

// ServicePlanFields is the definition of the service plan meta
type ServicePlanFields struct {
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
