package http

import (
	"reflect"
	"testing"
)

func isExpectedHeaders(got map[string]string, expected map[string]string) bool {
	return reflect.DeepEqual(got, expected)
}

func isExpectedReq(got Request, expected Request) bool {
	return reflect.DeepEqual(got, expected)
}

var testHeaders = map[string]string{"Host": "test.com", "User-Agent": "curl/7.64.1", "Accept": "*/*"}

func TestReq(t *testing.T) {
	tests := []struct {
		name        string
		arg         string
		expected    Request
		expectedErr bool
	}{
		{"Test parse req 1", "GET / HTTP/1.1\r\nHost: test.com\r\nUser-Agent: curl/7.64.1\r\nAccept: */*\r\n12345678", Request{
			Headers:    testHeaders,
			Path:       "/",
			Protocol:   "HTTP/1.1",
			Method:     "GET",
			Body:       []byte("12345678"),
			PathParams: RoutePathParams{},
		}, false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := ParseReq([]byte(test.arg))

			if (err != nil) != test.expectedErr {
				t.Errorf("Parse() error = %v, is error expected %v", err, test.expectedErr)
			}

			if !isExpectedReq(got, test.expected) {
				t.Errorf("ParseReq(`%s`) = \"%+v\"; want \"%+v\"", test.arg, got, test.expected)
			}
		})
	}
}

func TestParseHeaders(t *testing.T) {
	tests := []struct {
		name     string
		arg      string
		expected map[string]string
	}{
		{"Test header parser 1", "Host: test.com\r\n\nUser-Agent: curl/7.64.1\r\nAccept: */*\r\n ", testHeaders},
		{"Test header parser 2", "", map[string]string{}},
		{"Test header parser 3", "invalid-header adsdfsdfsdf, sdfksdfosdf:asdasdsad", map[string]string{}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := parseHeaders(test.arg)
			if !isExpectedHeaders(got, test.expected) {
				t.Errorf("ParseHeaders(`%s`) = \"%+q\"; want \"%+q\"", test.arg, got, test.expected)
			}
		})
	}
}
