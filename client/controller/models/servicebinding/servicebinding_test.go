package servicebinding

import (
	"fmt"
	"github.com/deis/deis/client/controller/client"
	"github.com/deis/deis/version"
	. "github.com/onsi/gomega"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

const serviceBindingFixture string = `
{
	"credentials": {
		"username": "mysql",
		"password": "password"
	}
}
`

type fakeHTTPServer struct {
}

func (svr *fakeHTTPServer) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("DEIS_API_VERSION", version.APIVersion)
	if req.URL.Path == "/v1/service_bindings" && req.Method == "POST" {
		res.Header().Add("Location", "/v1/service_bindings/serivce_binding_id")
		res.Write([]byte(serviceBindingFixture))
		return
	}

	fmt.Printf("Unrecongized URL %s\n", req.URL)
	res.WriteHeader(http.StatusNotFound)
	res.Write(nil)
}

func TestBindService(t *testing.T) {
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

	err = Bind(&client, "service_instance_id", "app_id", nil)
	Expect(err).To(BeNil())
}
