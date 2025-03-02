package http

type Response struct {
	Body []byte
	Code int
	Headers
	FilePath string
}
