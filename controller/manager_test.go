package controller

import (
	"testing"

	"github.com/labstack/echo"
)

func TestNewManager(t *testing.T) {
	counter := 0
	matchFunc := func(methods []string, path string, handler echo.HandlerFunc, middleware ...echo.MiddlewareFunc) []*echo.Route {
		counter++
		return nil
	}

	data := &Data{
		Controllers: []*Controller{
			{Path: "1"},
			{Path: "2"},
			{Path: "3"},
		},
		MatchFunc: matchFunc,
	}
	data.Match()

	if counter != 3 {
		t.Fatalf("expected counter 3, got %v", counter)
	}
}
