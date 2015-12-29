package services

import (
	"fmt"
	"github.com/deis/deis/client/controller/api"
	"github.com/deis/deis/client/controller/client"
	"github.com/deis/deis/version"
	. "github.com/onsi/gomega"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
)

const servicesFixture string = `
{
    "count": 1,
    "next": null,
    "previous": null,
    "results": [
        {
            "id": "id",
            "label": "label",
			"plans": [
				{
					"name": "small"
				}
			]
        }
    ]
}`

const serviceInstanceCreated string = `
{
	"id": "service_instance_id",
	"name": "mysql",
	"service_id": "mysql_service_id",
	"plan_id": "mysql_small_plan_id"
}`

type fakeHTTPServer struct {
}

func (svr *fakeHTTPServer) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("DEIS_API_VERSION", version.APIVersion)
	if req.URL.Path == "/v1/services/" && req.Method == "GET" {
		res.Write([]byte(servicesFixture))
		return
	}

	if req.URL.Path == "/v1/service-instances/" && req.Method == "POST" {
		res.WriteHeader(http.StatusCreated)
		res.Header().Add("Location", "/v1/service-instances/id")
		res.Write([]byte(serviceInstanceCreated))
		return
	}

	fmt.Printf("Unrecongized URL %s\n", req.URL)
	res.WriteHeader(http.StatusNotFound)
	res.Write(nil)
}

func TestListService(t *testing.T) {
	RegisterTestingT(t)
	t.Parallel()

	handler := fakeHTTPServer{}
	server := httptest.NewServer(&handler)
	defer server.Close()

	u, err := url.Parse(server.URL)

	if err != nil {
		t.Fatal(err)
	}

	httpClient := client.CreateHTTPClient(false)

	client := client.Client{HTTPClient: httpClient, ControllerURL: *u, Token: "abc"}

	services, _, err := List(&client, 10)

	Expect(err).To(BeNil())

	Expect(len(services)).To(Equal(1))

	expected := []api.ServiceOffering{
		api.ServiceOffering{
			ServiceOfferingFields: api.ServiceOfferingFields{
				ID:    "id",
				Label: "label",
			},
			Plans: []api.ServicePlanFields{
				api.ServicePlanFields{
					Name: "small",
				},
			},
		},
	}

	if !reflect.DeepEqual(expected, services) {
		t.Errorf("Expected %v, Got %v", expected, services)
	}
}

func TestCreateService(t *testing.T) {
	RegisterTestingT(t)
	t.Parallel()

	handler := fakeHTTPServer{}
	server := httptest.NewServer(&handler)
	defer server.Close()

	u, err := url.Parse(server.URL)

	if err != nil {
		t.Fatal(err)
	}

	httpClient := client.CreateHTTPClient(false)

	client := client.Client{HTTPClient: httpClient, ControllerURL: *u, Token: "abc"}

	service, err := New(&client, "service_instance_name", "service_name", "plan_name")

	Expect(err).To(BeNil())

	Expect(service).ToNot(BeNil())
}
