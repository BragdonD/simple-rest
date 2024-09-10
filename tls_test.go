package simplerest_test

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"

	simplerest "github.com/BragdonD/simple-rest"
)

func TestServeTls(t *testing.T) {
	// create a server
	server := simplerest.NewServer(
		"localhost",
		9000,
		simplerest.WithHttps("./test/certs/server.crt", "./test/certs/server.key"))
	defer server.Close()
	const data = "Hello world!"
	server.HandleFunc("/hello", nil, func(w http.ResponseWriter, r *http.Request, p simplerest.Parameters) error {
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte(data))
		return nil
	}, http.MethodGet)
	// launch server
	go func() {
		server.Serve()
	}()

	// create a https client
	certPool, err := x509.SystemCertPool()
	if err != nil {
		t.Error(err)
	}
	// load server ca
	if caCertPEM, err := os.ReadFile("./test/ca/ca.crt"); err != nil {
		t.Error(err)
	} else if ok := certPool.AppendCertsFromPEM(caCertPEM); !ok {
		panic("invalid cert in CA PEM")
	}
	tlsConfig := &tls.Config{
		RootCAs: certPool,
	}
	tr := &http.Transport{
		TLSClientConfig: tlsConfig,
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Get("https://localhost:9000/hello")
	if err != nil {
		t.Error(err)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	if strings.Compare(string(body), data) != 0 {
		t.Error(fmt.Errorf("the data sent and received do not match, data: [%s], body: [%s]", data, string(body)))
	}
}
