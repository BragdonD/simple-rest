package simplerest

import (
	"regexp"
	"strings"
)

// route represents a route in a RESTful server.
type Route struct {
	// The path to match.
	Path string
	// parametersPosition is the number of "/" minus 1 before
	// finding the string corresponding to the parameter in the [url.URL]
	parametersPosition map[string]int
}

const (
	// DynamicPathParametersRegex is the regex to match dynamic path parameters.
	DynamicPathParametersRegex = `{([^/]+)}`
)

// ParseDynamicPathParameters parses the dynamic path parameters from the
// [route] path.
func (r *Route) ParseDynamicPathParameters() {
	r.parametersPosition = make(map[string]int)
	re := regexp.MustCompile(DynamicPathParametersRegex)
	pathParts := strings.Split(r.Path, "/")
	for i, pathPart := range pathParts {
		if re.Match([]byte(pathPart)) {
			cleanParam := strings.Replace(pathPart, "{", "", 1)
			cleanParam = strings.Replace(cleanParam, "}", "", 1)
			r.parametersPosition[cleanParam] = i - 1
		}
	}
}

// GetDynamicPathParameters returns the dynamic parameters position inside
// the [route] path.
func (r *Route) GetDynamicPathParameters() map[string]int {
	return r.parametersPosition
}
