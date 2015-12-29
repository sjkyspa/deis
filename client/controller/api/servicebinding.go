package api

// ServiceBindingCreateRequest is the request to create service bindings
type ServiceBindingCreateRequest struct {
	ServiceInstanceID string                 `json:"service_instance_id"`
	AppID             string                 `json:"app_id"`
	Params            map[string]interface{} `json:"parameters,omitempty"`
}

// ServiceBindingFields is the definition of the service binding
type ServiceBindingFields struct {
	ID    string `json:"id"`
	URL   string `json:"url"`
	AppID string `json:"app_id"`
}
