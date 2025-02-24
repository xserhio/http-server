package server

import (
	"fmt"
	"github.com/codecrafters-io/http-server-starter-go/app/pkg/http"
	"io"
	"net"
)

type Handler func(req *http.Request) *http.Response

type Server struct {
	handlers       *map[string]Handler
	defaultHandler *Handler
}

func NewServer() *Server {
	handlers := make(map[string]Handler, 10)

	return &Server{handlers: &handlers, defaultHandler: nil}
}

func (s *Server) RegisterDefaultHandler(handler Handler) {
	(*s).defaultHandler = &handler
}

func (s *Server) RegisterHandler(route string, handler Handler) {
	(*s.handlers)[route] = handler
}

func (s *Server) handleConnection(conn net.Conn) ([]byte, error) {
	buffer := make([]byte, 1024)

	for {
		n, err := conn.Read(buffer)

		if err != nil && err != io.EOF {
			return nil, err
		}

		if n == 0 {
			return nil, fmt.Errorf("connection closed")
		}

		return buffer[:n], nil
	}
}

func (s *Server) sendResponse(conn net.Conn, res http.Response) error {
	serRes, err := http.SerializeRes(res)

	if err != nil {
		return err
	}

	_, err = conn.Write(serRes)

	defer conn.Close()

	return err
}

func (s *Server) Listen(port int) error {
	l, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", port))

	if err != nil {
		return err
	}

	for {
		conn, err := l.Accept()

		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			continue
		}

		reqRaw, err := s.handleConnection(conn)

		if err != nil {
			conn.Close()
			continue
		}

		req, err := http.Parse(reqRaw)

		if err != nil {
			fmt.Println("Error parsing request: ", err.Error())
			conn.Close()
			continue
		}

		handler, ok := (*s.handlers)[req.Path]

		if !ok {
			handler = *s.defaultHandler
		}

		res := handler(&req)

		err = s.sendResponse(conn, *res)

		if err != nil {
			fmt.Println("Error handling response: ", err.Error())
		}

		conn = nil
	}
}
