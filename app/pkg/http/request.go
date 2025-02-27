package http

type RoutePathParams map[string]string

type Request struct {
	Path       string
	PathParams RoutePathParams
	Method     string
	Headers    map[string]string
	Protocol   string
}
