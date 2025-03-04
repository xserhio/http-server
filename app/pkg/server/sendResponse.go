package server

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/codecrafters-io/http-server-starter-go/app/pkg/compress"
	"github.com/codecrafters-io/http-server-starter-go/app/pkg/http"
	"net"
	"os"
	"strconv"
)

func handleSendFile(res *http.Response) error {
	if res.FilePath == "" {
		return nil
	}

	f, err := os.Open(res.FilePath)

	if err != nil {
		return err
	}

	defer f.Close()

	stat, err := f.Stat()
	if err != nil {
		return err
	}

	bs := make([]byte, stat.Size())
	_, err = bufio.NewReader(f).Read(bs)

	if err != nil {
		return err
	}

	res.Body = bs

	res.Headers["Content-Length"] = fmt.Sprint(stat.Size())
	res.Headers["Content-Type"] = "application/octet-stream"

	return err
}

func compressResponse(req *http.Request, res *http.Response) error {
	encoding, ok := req.Headers["accept-encoding"]

	if !ok {
		return nil
	}

	compressHandler := compress.GetCompressHandler(encoding)

	if compressHandler == nil {
		return nil
	}

	compressedBody, err := compressHandler(res.Body)

	if err != nil {
		return err
	}

	res.Body = compressedBody
	res.Headers["Content-Encoding"] = encoding

	return nil
}

func (s *Server) sendResponse(conn *net.Conn, res *http.Response, req *http.Request) error {
	if res.Headers == nil {
		res.Headers = make(http.Headers, 2)
	}

	err := handleSendFile(res)

	if errors.Is(err, os.ErrNotExist) {
		res = (*s.defaultHandler)(req)

		if res.Headers == nil {
			res.Headers = make(http.Headers, 2)
		}
	} else if err != nil {
		return err
	}

	if res.FilePath == "" {
		res.Headers["Content-Length"] = strconv.Itoa(len(res.Body))
	}

	err = compressResponse(req, res)

	if err != nil {
		s.SendErr(500, map[string]string{"error": "failed to compress response"})
	}

	serRes, err := http.SerializeRes(*res)

	if err != nil {
		return err
	}

	_, err = (*conn).Write(serRes)

	defer (*conn).Close()

	return err
}
