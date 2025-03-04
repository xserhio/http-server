package server

import (
	"fmt"
	"github.com/codecrafters-io/http-server-starter-go/app/pkg/http"
	"io"
	"net"
)

type Handler func(req *http.Request) *http.Response

type HandlerData struct {
	handler *Handler
	method  string
	route   string
}

type Server struct {
	handlers       []HandlerData
	defaultHandler *Handler
}

func (s *Server) internalServerError(conn *net.Conn, req *http.Request, err error) error {
	res := &http.Response{}
	res.Code = 500
	res.Body = []byte(err.Error())
	res.Headers = make(http.Headers, 2)
	res.Headers["Content-Type"] = "text/plain"

	err = s.sendResponse(conn, res, req)

	return err
}

func NewServer() *Server {
	handlers := make([]HandlerData, 10)

	return &Server{handlers: handlers, defaultHandler: nil}
}

func (s *Server) RegisterDefaultHandler(handler Handler) {
	(*s).defaultHandler = &handler
}

func (s *Server) RegisterHandler(route string, method string, handler Handler) {
	s.handlers = append(s.handlers, HandlerData{handler: &handler, method: method, route: route})
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

func (s *Server) handleConn(conn *net.Conn) {
	reqRaw, err := s.handleConnection(*conn)

	defer (*conn).Close()

	if err != nil {
		return
	}

	req, err := http.ParseReq(reqRaw)

	if err != nil {
		fmt.Println("Error parsing request: ", err.Error())
		return
	}

	handler, reqPathParams, ok := s.router(req.Path, req.Method)

	if !ok {
		handler = s.defaultHandler
		req.PathParams = nil
	} else {
		req.PathParams = reqPathParams
	}

	res := (*handler)(&req)

	err = s.sendResponse(conn, res, &req)

	if err != nil {
		_ = s.internalServerError(conn, &req, err)
	}
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
			conn.Close()
			return err
		}

		go s.handleConn(&conn)
	}
}
