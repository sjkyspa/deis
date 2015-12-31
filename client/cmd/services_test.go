package cmd

import (
	"github.com/deis/deis/client/controller/client"
	"github.com/deis/deis/version"
	. "github.com/onsi/gomega"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
)

type fakeHTTPServer struct {
}

const serviceInstanceCreated string = `
{
	"id": "service_instance_id",
	"name": "mysql",
	"service_id": "mysql_service_id",
	"plan_id": "mysql_small_plan_id"
}`

const servicesFixture string = `
{
    "count": 1,
    "next": null,
    "previous": null,
    "results": [
        {
            "uuid": "uuid",
            "label": "label",
			"plans": [
				{
					"name": "small",
					"uuid": "uuid",
					"free": true,
					"description": "desc"
				}
			]
        }
    ]
}`

const serviceBindingFixture string = `
{
	"credentials": {
		"username": "mysql",
		"password": "password"
	}
}
`

func (c *fakeHTTPServer) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("DEIS_API_VERSION", version.APIVersion)
	if req.URL.Path == "/v1/service_instances" && req.Method == "POST" {
		res.WriteHeader(http.StatusCreated)
		res.Header().Add("Location", "/v1/service_instances/id")
		res.Write([]byte(serviceInstanceCreated))
		return
	}

	if req.URL.Path == "/v1/services" && reflect.DeepEqual(req.URL.Query()["name"], []string{"service_name"}) && req.Method == "GET" {
		res.WriteHeader(http.StatusOK)
		res.Write([]byte(servicesFixture))
		return
	}

	if req.URL.Path == "/v1/service_instances" && reflect.DeepEqual(req.URL.Query()["name"], []string{"service_instance_name"}) && req.Method == "GET" {
		res.WriteHeader(http.StatusOK)
		res.Write([]byte(servicesFixture))
		return
	}

	if req.URL.Path == "/v1/service_bindings" && req.Method == "POST" {
		res.Header().Add("Location", "/v1/service_bindings/serivce_binding_id")
		res.Write([]byte(serviceBindingFixture))
		return
	}

	res.WriteHeader(http.StatusNotFound)
}

func TestCreateServiceSuccess(t *testing.T) {
	RegisterTestingT(t)
	t.Parallel()

	server := httptest.NewServer(&fakeHTTPServer{})
	defer server.Close()

	u, err := url.Parse(server.URL)

	if err != nil {
		t.Fatal(err)
	}

	httpClient := client.CreateHTTPClient(false)
	c := client.Client{HTTPClient: httpClient, ControllerURL: *u, Token: "abc"}

	err = ServiceCreate(&c, "service_name", "plan_name", "service_instance_name")

	Expect(err).To(BeNil())
}

func TestBindServiceSuccess(t *testing.T) {
	RegisterTestingT(t)
	t.Parallel()

	server := httptest.NewServer(&fakeHTTPServer{})
	defer server.Close()

	u, err := url.Parse(server.URL)

	if err != nil {
		t.Fatal(err)
	}

	httpClient := client.CreateHTTPClient(false)
	c := client.Client{HTTPClient: httpClient, ControllerURL: *u, Token: "abc"}

	err = ServiceBind(&c, "app_name", "service_instance_name")

	Expect(err).To(BeNil())
}
