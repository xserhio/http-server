package main

import (
	"encoding/json"
	"flag"
	"github.com/codecrafters-io/http-server-starter-go/app/pkg/http"
	"github.com/codecrafters-io/http-server-starter-go/app/pkg/server"
	"os"
	"path"
	"strconv"
)

func main() {
	var fileDir string
	flag.StringVar(&fileDir, "directory", "", "dir with files that served by server")

	flag.Parse()

	s := server.NewServer()

	s.RegisterHandler("/", "GET", func(req *http.Request) *http.Response {
		return &http.Response{
			Code:    200,
			Body:    nil,
			Headers: nil,
		}
	})

	s.RegisterHandler("/user-agent", "GET", func(req *http.Request) *http.Response {
		userAgent := req.Headers["user-agent"]

		if userAgent == "" {
			return &http.Response{
				Code:    400,
				Headers: nil,
				Body:    nil,
			}
		}

		contentLength := len(userAgent)

		return &http.Response{
			Code:    200,
			Body:    []byte(userAgent),
			Headers: http.Headers{"Content-Type": "text/plain", "Content-Length": strconv.Itoa(contentLength)},
		}
	})

	s.RegisterHandler("/echo/:str", "GET", func(req *http.Request) *http.Response {
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

	s.RegisterHandler("/files/:fileName", "POST", func(req *http.Request) *http.Response {
		fileName, ok := req.PathParams["fileName"]

		if !ok {
			return s.SendErr(400, map[string]string{
				"error": "fileName must be string",
			})
		}

		var filePath string

		if fileDir != "" {
			filePath = path.Join(fileDir, fileName)
		}

		file, err := os.Create(filePath)

		if err != nil {
			return s.SendErr(500, map[string]string{"error": "failed to create file"})
		}

		_, err = file.Write(req.Body)

		if err != nil {
			return s.SendErr(500, map[string]string{"error": "failed to write to file"})
		}

		return &http.Response{
			Code: 201,
		}
	})

	s.RegisterHandler("/files/:fileName", "GET", func(req *http.Request) *http.Response {
		fileName, ok := req.PathParams["fileName"]

		if !ok {
			resBody, err := json.Marshal(map[string]string{
				"error": "fileName must be string",
			})

			if err != nil {
				return &http.Response{
					Code: 500,
				}
			}

			return &http.Response{
				Code:    400,
				Body:    resBody,
				Headers: http.Headers{"Content-Type": "application/json", "Content-Length": strconv.Itoa(len(resBody))},
			}
		}

		var filePath string

		if fileDir != "" {
			filePath = path.Join(fileDir, fileName)
		}

		return &http.Response{
			Code:     200,
			FilePath: filePath,
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
