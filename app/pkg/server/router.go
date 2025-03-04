package server

import (
	"github.com/codecrafters-io/http-server-starter-go/app/pkg/http"
	"slices"
	"strings"
)

func (s *Server) router(reqPath string, method string) (*Handler, http.RoutePathParams, bool) {
	reqPathParts := strings.Split(reqPath, "/")

	var routeVariants []string

	routes := make([]string, len(s.handlers))

	for _, handlerData := range s.handlers {
		route := handlerData.route

		exist := slices.Contains(routes, route)

		if !exist {
			routes = append(routes, route)
		}
	}

	for _, route := range routes {
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
	
	for _, handlerData := range s.handlers {
		isFind := handlerData.method == method && handlerData.route == route

		if isFind {
			return handlerData.handler, routePathParams, true
		}
	}

	return nil, nil, false
}
