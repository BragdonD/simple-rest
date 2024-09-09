package simplerest

import (
	"fmt"
	"net/http"
)

// RouteUnsupportedHttpMethod is an error that indicates that the
// http request method is not supported by the server route.
type RouteUnsupportedHttpMethod struct {
	// Route is the request route
	Route string
	// Method is the requested method
	Method string
	// SupportedMethods are the supported method by the route
	SupportedMethods []string
}

// Error is the interface method for [RouteUnsupportedHttpMethod]
// to be considered an error.
func (ruhm *RouteUnsupportedHttpMethod) Error() string {
	return fmt.Sprintf(
		"The route [%s] does not support the HTTP Method [%v]. It supports the following Methods [%s]",
		ruhm.Route, ruhm.Method, ruhm.SupportedMethods,
	)
}

// HttpWriteError writes an [error] to the given [http.ResponseWriter].
func HttpWriteError(w http.ResponseWriter, err error) error {
	switch err.(type) {
	case *RouteUnsupportedHttpMethod:
		w.WriteHeader(http.StatusMethodNotAllowed)
	default:
		w.WriteHeader(http.StatusBadRequest)
	}
	if _, err := w.Write([]byte(err.Error())); err != nil {
		return fmt.Errorf("unable to write error to response writer: %v", err)
	}
	return nil
}

// RouteAlreadyExists is an error that indicates that the new
// instanciated [Route] already exists on the [Server].
type RouteAlreadyExists struct {
	// Route is the route that should have been instantiated
	Route string
}

// Error is the interface method for [RouteAlreadyExists]
// to be considered an error.
func (rae *RouteAlreadyExists) Error() string {
	return fmt.Sprintf("the route [%s] already exists on this server.", rae.Route)
}
