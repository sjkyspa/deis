package api

// Broker is the definition of the broker object.
type Broker struct {
	Created  string `json:"created"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
	URL      string `json:"url"`
}

// BrokerCreateRequest is the definition of POST /v1/brokers/.
type BrokerCreateRequest struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
	URL      string `json:"url"`
}
