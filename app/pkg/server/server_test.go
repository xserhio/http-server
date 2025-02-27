package server

import (
	"github.com/codecrafters-io/http-server-starter-go/app/pkg/http"
	"reflect"
	"testing"
)

var server = NewServer()

func init() {
	server.RegisterHandler("/echo/:arg", func(req *http.Request) *http.Response {
		resBody := req.PathParams["arg"]

		return &http.Response{Code: 200, Body: []byte(resBody), Headers: nil}
	})
}

func TestServer(t *testing.T) {
	tests := []struct {
		name               string
		reqPath            string
		expectedPathParams http.RoutePathParams
		ok                 bool
	}{
		{"Test server router 1", "/echo/test", http.RoutePathParams{"arg": "test"}, true},
		{"Test server router 2", "/echo", http.RoutePathParams{}, false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, gotReqPathParams, ok := server.router(test.reqPath)

			if ok != test.ok {
				t.Errorf("got ok %v, want %v", ok, test.ok)
			}

			if !reflect.DeepEqual(gotReqPathParams, test.expectedPathParams) {
				t.Errorf("got %#v, want %#v", gotReqPathParams, test.expectedPathParams)
			}
		})
	}
}
