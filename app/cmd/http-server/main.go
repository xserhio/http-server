package main

import (
	"github.com/codecrafters-io/http-server-starter-go/app/pkg/http"
	server "github.com/codecrafters-io/http-server-starter-go/app/pkg/server"
	"strconv"
)

func main() {
	s := server.NewServer()

	s.RegisterHandler("/", func(req *http.Request) *http.Response {
		return &http.Response{
			Code:    200,
			Body:    nil,
			Headers: nil,
		}
	})

	s.RegisterHandler("/echo/:str", func(req *http.Request) *http.Response {
		str, ok := req.PathParams["str"]

		if !ok {
			str = ""
		}

		contentLength := len(str)

		return &http.Response{
			Code:    200,
			Body:    []byte(str),
			Headers: http.Headers{"Content-Type": "text/plain", "Content-Length": strconv.Itoa(contentLength)},
		}
	})

	s.RegisterDefaultHandler(func(req *http.Request) *http.Response {
		return &http.Response{
			Code:    404,
			Body:    nil,
			Headers: nil,
		}
	})

	err := s.Listen(4221)

	if err != nil {
		panic(err)
	}
}
