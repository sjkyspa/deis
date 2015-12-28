package brokers

import (
	"fmt"
	"github.com/deis/deis/client/controller/api"
	"github.com/deis/deis/client/controller/client"
	"github.com/deis/deis/version"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
)

type fakeHTTPServer struct {
	Name     string
	Username string
	Password string
	URL      string
}

const brokerCreateExpected string = `{"name":"mysql","username":"mysql","password":"mysql","url":"http://localhost"}`
const brokerFixture string = `
{
	"created": "2014-01-01T00:00:00UTC",
    "name": "mysql",
    "username": "mysql",
    "url": "http://localhost"
}`

const brokersFixture string = `
{
    "count": 1,
    "next": null,
    "previous": null,
    "results": [
        {
            "created": "2014-01-01T00:00:00UTC",
            "name": "mysql",
            "username": "mysql",
            "url": "http://localhost"
        }
    ]
}`

func (svr *fakeHTTPServer) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("DEIS_API_VERSION", version.APIVersion)
	if req.URL.Path == "/v1/brokers/" && req.Method == "POST" {
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			fmt.Println(err)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write(nil)
		}

		if string(body) == brokerCreateExpected {
			res.WriteHeader(http.StatusCreated)
			res.Write([]byte(brokerFixture))
			return
		}

		fmt.Printf("Unexpected Body: %s'\n", body)
		res.WriteHeader(http.StatusInternalServerError)
		res.Write(nil)
		return
	}

	if req.URL.Path == "/v1/brokers" && req.Method == "GET" {
		res.Write([]byte(brokersFixture))
		return
	}
}

func TestBrokerCreated(t *testing.T) {
	t.Parallel()

	expected := api.Broker{
		Created:  "2014-01-01T00:00:00UTC",
		Name:     "mysql",
		Username: "mysql",
		URL:      "http://localhost",
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

	for _, broker := range []api.BrokerCreateRequest{
		api.BrokerCreateRequest{
			Name:     "mysql",
			Username: "mysql",
			Password: "mysql",
			URL:      "http://localhost",
		},
	} {
		brokerURL, _ := url.Parse(broker.URL)

		actual, err := New(&client, broker.Name, broker.Username, broker.Password, *brokerURL)

		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(expected, actual) {
			t.Errorf("Expected %v, Got %v", expected, actual)
		}
	}
}

func TestBrokerList(t *testing.T) {
	t.Parallel()

	expected := []api.Broker{
		api.Broker{
			Created:  "2014-01-01T00:00:00UTC",
			Name:     "mysql",
			Username: "mysql",
			URL:      "http://localhost",
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

	brokers, _, err := List(&client, 100)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(expected, brokers) {
		t.Errorf("Expected %v, Got %v", expected, brokers)
	}
}
