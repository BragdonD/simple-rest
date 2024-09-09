package simplerest_test

import (
	"testing"

	simplerest "github.com/BragdonD/simple-rest"
	"github.com/go-test/deep"
)

func TestParseDynamicPathParameters(t *testing.T) {
	route := simplerest.Route{
		Path: "/users/{id}/test/{name}/test",
	}
	route.ParseDynamicPathParameters()
	params := route.GetDynamicPathParameters()
	expected := map[string]int{
		"id":   1,
		"name": 3,
	}
	if diff := deep.Equal(params, expected); diff != nil {
		t.Error(diff)
	}
}
