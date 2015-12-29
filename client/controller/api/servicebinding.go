package api

// ServiceBindingCreateRequest is the request to create service bindings
type ServiceBindingCreateRequest struct {
	ServiceInstanceID string                 `json:"service_instance_id"`
	AppID             string                 `json:"app_id"`
	Params            map[string]interface{} `json:"parameters,omitempty"`
}
