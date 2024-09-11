package simplerest_test

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"

	simplerest "github.com/bragdond/simple-rest"
	"github.com/go-test/deep"
)

func TestRetrievePathParameters(t *testing.T) {
	route := simplerest.Route{
		Path: "/test/{id}/test/{name}",
	}
	route.ParseDynamicPathParameters()
	result := simplerest.RetrievePathParameters(&route, &url.URL{
		Path: "/test/123/test/thomas",
	})
	expected := map[string]string{
		"id":   "123",
		"name": "thomas",
	}
	if diff := deep.Equal(expected, result); diff != nil {
		t.Fatal(diff)
	}
}

func TestServe(t *testing.T) {
	server := simplerest.NewServer("localhost", 9000)
	defer server.Close()
	const data = "Hello world!"
	server.HandleFunc("/hello", nil, func(w http.ResponseWriter, r *http.Request, p simplerest.Parameters) error {
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte(data))
		return nil
	}, http.MethodGet)
	go func() {
		server.Serve()
	}()
	client := &http.Client{}
	resp, err := client.Get("http://localhost:9000/hello")
	if err != nil {
		t.Fatal(err)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Compare(string(body), data) != 0 {
		t.Fatal(fmt.Errorf("the data sent and received do not match, data: [%s], body: [%s]", data, string(body)))
	}
}
