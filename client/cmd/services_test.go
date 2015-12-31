package cmd

import (
	"encoding/json"
	"github.com/deis/deis/client/controller/client"
	"github.com/deis/deis/version"
	. "github.com/onsi/gomega"
	"io/ioutil"
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

const serviceInstancesFixture string = `
{
    "count": 1,
    "next": null,
    "previous": null,
    "results": [
        {
        	"uuid": "service_instance_id",
            "plan_uuid": "service_plan_id",
            "name": "mysql",
            "type": "managed_service_instance"
        }
    ]
}`

const servicesFixture string = `
{
    "count": 1,
    "next": null,
    "previous": null,
    "results": [
        {
            "uuid": "service_instance_id",
            "label": "label",
			"plans": [
				{
					"name": "plan_name",
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

const serviceBindingsFixture string = `
{
    "count": 1,
    "next": null,
    "previous": null,
    "results": [
        {
        	"service_instance_id": "service_instance_id",
        	"app_id": "app_id",
   			"uuid": "uuid"
        }
    ]
}
`

func (c *fakeHTTPServer) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("DEIS_API_VERSION", version.APIVersion)
	if req.URL.Path == "/v1/service_instances" && req.Method == "POST" {
		var reqJSON map[string]interface{}
		data, err := ioutil.ReadAll(req.Body)
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		err = json.Unmarshal(data, &reqJSON)
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		if reqJSON["name"].(string) == "service_instance_name" {
			res.WriteHeader(http.StatusCreated)
			res.Header().Add("Location", "/v1/service_instances/id")
			res.Write([]byte(serviceInstanceCreated))
			return
		}

		res.WriteHeader(http.StatusBadRequest)
		return
	}

	if req.URL.Path == "/v1/services" && reflect.DeepEqual(req.URL.Query()["name"], []string{"service_name"}) && req.Method == "GET" {
		res.WriteHeader(http.StatusOK)
		res.Write([]byte(servicesFixture))
		return
	}

	if req.URL.Path == "/v1/service_instances" && reflect.DeepEqual(req.URL.Query()["name"], []string{"service_instance_name"}) && req.Method == "GET" {
		res.WriteHeader(http.StatusOK)
		res.Write([]byte(serviceInstancesFixture))
		return
	}

	if req.URL.Path == "/v1/service_bindings" && req.Method == "POST" {
		res.Header().Add("Location", "/v1/service_bindings/serivce_binding_id")
		res.Write([]byte(serviceBindingFixture))
		return
	}

	if req.URL.Path == "/v1/service_bindings" && req.Method == "GET" {
		res.Write([]byte(serviceBindingsFixture))
		return
	}

	if req.URL.Path == "/v1/service_bindings/uuid" && req.Method == "DELETE" {
		res.Write([]byte(serviceBindingsFixture))
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

func TestCreateServiceFailPlanNotFound(t *testing.T) {
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

	err = ServiceCreate(&c, "service_name", "not_existed_plan", "service_instance_name")

	Expect(err).NotTo(BeNil())
}

func TestCreateServiceFailServiceNameNotFound(t *testing.T) {
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

	err = ServiceCreate(&c, "not_existed_service_name", "plan_name", "service_instance_name")

	Expect(err).NotTo(BeNil())
}

func TestCreateServiceFailDuplicatedServiceInstanceName(t *testing.T) {
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

	err = ServiceCreate(&c, "service_name", "plan_name", "duplicated_service_instance_name")

	Expect(err).NotTo(BeNil())
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

func TestBindServiceFailServiceInstanceNotFound(t *testing.T) {
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

	err = ServiceBind(&c, "app_id", "not_existed_service_instance_name")

	Expect(err).NotTo(BeNil())
}

func TestUnbindServiceSuccess(t *testing.T) {
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

	err = ServiceUnbind(&c, "app_id", "service_instance_name")

	Expect(err).To(BeNil())
}

func TestUnbindServiceFailNoAppFound(t *testing.T) {
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

	err = ServiceUnbind(&c, "not_exists_app_id", "service_instance_name")

	Expect(err.Error()).To(MatchRegexp("Can Not find service by service instance.*"))
}

func TestUnbindServiceFailNoServiceInstanceNameFound(t *testing.T) {
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

	err = ServiceUnbind(&c, "app_id", "not_exists_service_instance_name")

	Expect(err.Error()).To(MatchRegexp("Can Not find service by service instance.*"))
}
