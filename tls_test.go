package simplerest_test

import (
	"crypto/tls"
	"crypto/x509"
	"io"
	"net/http"
	"os"
	"testing"

	simplerest "github.com/BragdonD/simple-rest"
)

func TestWithHttps(t *testing.T) {
	// Create server
	server := simplerest.NewServer(
		"localhost",
		9000,
		simplerest.WithHttps("./test/certs/server.crt", "./test/certs/server.key"),
	)
	defer server.Close()

	const data = "Hello world!"
	server.HandleFunc("/hello", nil, func(w http.ResponseWriter, r *http.Request, p simplerest.Parameters) error {
		w.WriteHeader(http.StatusAccepted)
		_, err := w.Write([]byte(data))
		return err
	}, http.MethodGet)

	// Launch server
	go server.Serve()

	// Create HTTPS client
	certPool, err := x509.SystemCertPool()
	if err != nil {
		t.Fatal(err)
	}

	// Load server CA
	caCertPEM, err := os.ReadFile("./test/ca/ca.crt")
	if err != nil {
		t.Fatal(err)
	}
	if ok := certPool.AppendCertsFromPEM(caCertPEM); !ok {
		t.Fatal("invalid cert in CA PEM")
	}

	tlsConfig := &tls.Config{RootCAs: certPool}
	client := &http.Client{
		Transport: &http.Transport{TLSClientConfig: tlsConfig},
	}

	// Send GET request
	resp, err := client.Get("https://localhost:9000/hello")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	if string(body) != data {
		t.Fatalf("mismatch: expected %q, got %q", data, string(body))
	}
}

func TestWithMtls(t *testing.T) {
	// Create server
	server := simplerest.NewServer(
		"localhost",
		9000,
		simplerest.WithMtls("./test/certs/server.crt", "./test/certs/server.key", "./test/ca/ca.crt"),
	)
	defer server.Close()

	const data = "Hello world!"
	server.HandleFunc("/hello", nil, func(w http.ResponseWriter, r *http.Request, p simplerest.Parameters) error {
		w.WriteHeader(http.StatusAccepted)
		_, err := w.Write([]byte(data))
		return err
	}, http.MethodGet)

	// Launch server
	go server.Serve()

	// Load client certificate
	clientTLSCert, err := tls.LoadX509KeyPair("./test/certs/client.crt", "./test/certs/client.key")
	if err != nil {
		t.Fatalf("error loading client certificate and key: %v", err)
	}

	// Create mTLS client
	certPool, err := x509.SystemCertPool()
	if err != nil {
		t.Fatal(err)
	}

	// Load server CA
	caCertPEM, err := os.ReadFile("./test/ca/ca.crt")
	if err != nil {
		t.Fatal(err)
	}
	if ok := certPool.AppendCertsFromPEM(caCertPEM); !ok {
		t.Fatal("invalid cert in CA PEM")
	}

	tlsConfig := &tls.Config{
		RootCAs:      certPool,
		Certificates: []tls.Certificate{clientTLSCert},
	}
	client := &http.Client{
		Transport: &http.Transport{TLSClientConfig: tlsConfig},
	}

	// Send GET request
	resp, err := client.Get("https://localhost:9000/hello")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	if string(body) != data {
		t.Fatalf("mismatch: expected %q, got %q", data, string(body))
	}
}
