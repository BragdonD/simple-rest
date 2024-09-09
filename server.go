package simplerest

import (
	"fmt"
	"net/http"
	"net/url"
	"slices"
	"strings"
)

// Server represents a RESTful server.
type Server struct {
	// The address to listen on.
	Address string
	// The port to listen on.
	Port int
	// The http server.
	server *http.Server
	// The http serve mux.
	mux *http.ServeMux
	// The routes handled by the server.
	routes []*Route
}

// Sopts is a function that sets server options.
type Sopts func(*Server)

// NewServer creates a new RESTful server.
func NewServer(address string, port int, opts ...Sopts) *Server {
	server := &Server{
		Address: address,
		Port:    port,
		mux:     http.NewServeMux(),
	}
	for _, opt := range opts {
		opt(server)
	}
	server.server = &http.Server{
		Addr:    fmt.Sprintf("%s:%d", address, port),
		Handler: server.mux,
	}
	return server
}

// retrievePathParameters retrieves the dynamic parameter from the defined [Route]
// and associate a value extracted from the [url.URL]
func RetrievePathParameters(route Route, url *url.URL) map[string]string {
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
	if slices.ContainsFunc(s.routes, func(r *Route) bool {
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
	s.routes = append(s.routes, route)
	// handle the route on the server
	s.mux.HandleFunc(path, func(w http.ResponseWriter, req *http.Request) {
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
		}
		handler(w, req, Parameters{})
	})
	return nil
}

// Serve launches the listening process of the [http.Server]
func (s *Server) Serve() error {
	return s.server.ListenAndServe()
}

// Close stops the listening process of the [http.Server]
func (s *Server) Close() error {
	return s.server.Close()
}
