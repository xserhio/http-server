package http

type Request struct {
	Path     string
	Method   string
	Headers  map[string]string
	Protocol string
}
