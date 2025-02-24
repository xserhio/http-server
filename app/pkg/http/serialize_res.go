package http

import (
	"fmt"
	netHttp "net/http"
)

func serializeHeaders(headers Headers) string {
	if headers == nil || len(headers) == 0 {
		return "\r\n"
	}

	resStr := ""

	for h, v := range headers {
		resStr += fmt.Sprintf("%s: %s\r\n", h, v)
	}

	return resStr
}

func SerializeRes(res Response) ([]byte, error) {
	resRaw := fmt.Sprintf("HTTP/1.1 %d %s\r\n", res.Code, netHttp.StatusText(res.Code))
	resRaw += serializeHeaders(res.Headers)
	resRaw += string(res.Body)

	return []byte(resRaw), nil
}
