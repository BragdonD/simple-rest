package simplerest

import (
	"net/http"
)

// Parameters is a map of path parameters.
type Parameters map[string]string

// Handler is a function that handles an HTTP request.
type Handler func(http.ResponseWriter, *http.Request, Parameters) error
