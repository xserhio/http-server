package server

import (
	"fmt"
	"github.com/codecrafters-io/http-server-starter-go/app/pkg/http"
	"io"
	"net"
	"strings"
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

func (s *Server) router(reqPath string) (Handler, http.RoutePathParams, bool) {
	reqPathParts := strings.Split(reqPath, "/")

	var routeVariants []string

	for route, _ := range *s.handlers {
		routeParts := strings.Split(route, "/")

		if len(routeParts) == len(reqPathParts) {
			routeVariants = append(routeVariants, route)
		}
	}

	if len(routeVariants) == 0 {
		return nil, http.RoutePathParams{}, false
	}

	var route string
	routePathParams := make(http.RoutePathParams)

	for _, routeVariant := range routeVariants {
		i := 0

		routeVariantParts := strings.Split(routeVariant, "/")

		for ; i < len(reqPathParts); i++ {
			reqPathPart := reqPathParts[i]
			routeVariantPart := routeVariantParts[i]

			if reqPathPart == routeVariantPart {
				continue
			}

			isRouteParam := strings.HasPrefix(routeVariantPart, ":")

			if isRouteParam {
				routeParamName := strings.TrimPrefix(routeVariantPart, ":")
				routePathParams[routeParamName] = reqPathPart
			} else {
				break
			}
		}

		if i == len(reqPathParts) {
			route = routeVariant
			break
		}
	}

	if route == "" {
		return nil, http.RoutePathParams{}, false
	}

	handler, ok := (*s.handlers)[route]

	return handler, routePathParams, ok
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

		go func() {
			if err != nil {
				fmt.Println("Error accepting connection: ", err.Error())
				conn.Close()
				return
			}

			reqRaw, err := s.handleConnection(conn)

			if err != nil {
				conn.Close()
				return
			}

			req, err := http.ParseReq(reqRaw)

			if err != nil {
				fmt.Println("Error parsing request: ", err.Error())
				conn.Close()
				return
			}

			handler, reqPathParams, ok := s.router(req.Path)

			if !ok {
				handler = *s.defaultHandler
				req.PathParams = nil
			} else {
				req.PathParams = reqPathParams
			}

			res := handler(&req)

			err = s.sendResponse(conn, *res)

			if err != nil {
				fmt.Println("Error handling response: ", err.Error())
			}

			conn.Close()
		}()
	}
}
