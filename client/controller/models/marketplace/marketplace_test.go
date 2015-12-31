package marketplace

import (
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

type fakeHTTPServer struct {
}

const servicesFixture string = `
{
    "count": 1,
    "next": null,
    "previous": null,
    "results": [
        {
            "uuid": "id",
            "label": "label",
			"plans": [
				{
					"name": "small"
				}
			]
        }
    ]
}`

func (svr *fakeHTTPServer) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("DEIS_API_VERSION", version.APIVersion)

	if req.URL.Path == "/v1/services" && req.Method == "GET" {
		res.Write([]byte(servicesFixture))
		return
	}
}

func TestMarketplace(t *testing.T) {
	t.Parallel()
	RegisterTestingT(t)

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

	handler := fakeHTTPServer{}
	server := httptest.NewServer(&handler)
	defer server.Close()

	u, err := url.Parse(server.URL)

	if err != nil {
		t.Fatal(err)
	}

	httpClient := client.CreateHTTPClient(false)

	client := client.Client{HTTPClient: httpClient, ControllerURL: *u, Token: "abc"}

	serviceOfferings, err := List(&client)
	Expect(err).To(BeNil())
	Expect(len(serviceOfferings)).To(Equal(1))

	if !reflect.DeepEqual(expected, serviceOfferings) {
		t.Errorf("Expected %v, Got %v", expected, serviceOfferings)
	}
}
