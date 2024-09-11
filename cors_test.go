package simplerest_test

import (
	"io"
	"net/http"
	"strings"
	"testing"

	simplerest "github.com/bragdond/simple-rest"
)

func TestWithCors(t *testing.T) {
	// Allowed origins and headers for CORS
	allowedOrigins := []string{"http://example.com", "http://localhost"}
	allowedHeaders := []string{"Content-Type", "Authorization"}

	// Create server with CORS enabled
	server := simplerest.NewServer(
		"localhost",
		9000,
		simplerest.WithCors(allowedOrigins, allowedHeaders),
	)
	defer server.Close()

	const data = "CORS Test"
	server.HandleFunc("/cors", nil, func(w http.ResponseWriter, r *http.Request, p simplerest.Parameters) error {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(data))
		return err
	}, http.MethodGet)

	// Launch the server
	go server.Serve()

	// Test client making a request to the CORS-enabled server
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, "http://localhost:9000/cors", nil)
	if err != nil {
		t.Fatal(err)
	}
	// Set Origin header to simulate a cross-origin request
	req.Header.Set("Origin", "http://example.com")

	// Perform the request
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	// Check if response body matches
	if string(body) != data {
		t.Fatalf("expected %q, got %q", data, string(body))
	}

	// Check CORS headers
	origin := resp.Header.Get("Access-Control-Allow-Origin")
	if !strings.Contains(strings.Join(allowedOrigins, ", "), origin) {
		t.Fatalf("expected Allow-Origin to be one of %v, got %v", allowedOrigins, origin)
	}

	methods := resp.Header.Get("Access-Control-Allow-Methods")
	if methods != "GET" {
		t.Fatalf("expected Allow-Methods to be GET, got %v", methods)
	}

	headers := resp.Header.Get("Access-Control-Allow-Headers")
	for _, hdr := range allowedHeaders {
		if !strings.Contains(headers, hdr) {
			t.Fatalf("expected header %v in Allow-Headers, got %v", hdr, headers)
		}
	}

	credentials := resp.Header.Get("Access-Control-Allow-Credentials")
	if credentials != "true" {
		t.Fatalf("expected Allow-Credentials to be true, got %v", credentials)
	}
}
