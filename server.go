package simplerest

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"slices"
	"strings"
)

// Server represents a RESTful server.
type Server struct {
	Address string // The address to listen on.
	Port    int    // The port to listen on.

	// HTTP server and routing
	HTTP struct {
		Server *http.Server   // The HTTP server.
		Mux    *http.ServeMux // The HTTP serve mux.
		Routes []*Route       // The routes handled by the server.
	}

	// TLS/MTLS configuration
	TLS struct {
		Cert     string   // Server certificate filepath.
		Key      string   // Certificate associated Key file
		ClientCA []string // Client CA for mTLS communication.
	}
}

// sopts is a function that sets server options.
type sopts func(any)

// NewServer creates a new RESTful server.
func NewServer(address string, port int, opts ...sopts) *Server {
	server := &Server{
		Address: address,
		Port:    port,
	}
	for _, opt := range opts {
		opt(server)
	}
	server.HTTP.Mux = &http.ServeMux{}
	server.HTTP.Server = &http.Server{
		Addr:    fmt.Sprintf("%s:%d", address, port),
		Handler: server.HTTP.Mux,
	}
	return server
}

// retrievePathParameters retrieves the dynamic parameter from the defined [Route]
// and associate a value extracted from the [url.URL]
func RetrievePathParameters(route *Route, url *url.URL) map[string]string {
	parametersValue := make(map[string]string)
	pathParts := strings.Split(url.Path, "/")
	for paramName, paramPos := range route.parametersPosition {
		if len(pathParts) > paramPos {
			parametersValue[paramName] = pathParts[paramPos+1]
		}
	}
	return parametersValue
}

// HandleFunc adds a route to the server.
func (s *Server) HandleFunc(path string, middleware Middleware, handler Handler, methods ...string) error {
	// check if the route already exists on the server
	if slices.ContainsFunc(s.HTTP.Routes, func(r *Route) bool {
		return r.Path == path
	}) {
		return &RouteAlreadyExists{
			Route: path,
		}
	}
	route := &Route{
		Path: path,
	}
	route.ParseDynamicPathParameters()
	s.HTTP.Routes = append(s.HTTP.Routes, route)

	wrappedHandler := handler
	// apply middleware if any
	if middleware != nil {
		wrappedHandler = middleware(handler)
	}
	// handle the route on the server
	s.HTTP.Mux.HandleFunc(path, func(w http.ResponseWriter, req *http.Request) {
		if !slices.Contains(methods, req.Method) {
			if err := HttpWriteError(w, &RouteUnsupportedHttpMethod{
				Method:           req.Method,
				Route:            path,
				SupportedMethods: methods,
			}); err != nil {
				err := fmt.Errorf("unable to write error to response writer: %v", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				// TODO: log the error
			}
			return
		}
		parameters := RetrievePathParameters(route, req.URL)
		if err := wrappedHandler(w, req, parameters); err != nil {
			// TODO: log the error
			fmt.Println(err)
		}
	})
	return nil
}

// Serve launches the listening process of the [http.Server]
func (s *Server) Serve() error {
	// At least HTTPS server
	if s.TLS.Cert != "" {
		// check server certificates
		cert, err := tls.LoadX509KeyPair(s.TLS.Cert, s.TLS.Key)
		if err != nil {
			return fmt.Errorf("no certificate could be loaded from this path: %v", err)
		}

		tlsConf := &tls.Config{
			Certificates: []tls.Certificate{cert},
		}

		// MTLS
		if len(s.TLS.ClientCA) != 0 {
			// check client certificates
			// tlsConf.ClientCAs = s.TLS.ClientCA
		}

		s.HTTP.Server.TLSConfig = tlsConf
		return s.HTTP.Server.ListenAndServeTLS("", "")
	}
	// default HTTP server
	return s.HTTP.Server.ListenAndServe()
}

// Close stops the listening process of the [http.Server]
func (s *Server) Close() error {
	return s.HTTP.Server.Close()
}
