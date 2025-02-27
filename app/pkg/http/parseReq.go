package http

import (
	"fmt"
	"strings"
)

func parseReqLine(reqLine string) (string, string, string, error) {
	parts := strings.Split(reqLine, " ")

	if len(parts) != 3 {
		return "", "", "", fmt.Errorf("invalid request line: %s", reqLine)
	}

	return parts[0], parts[1], parts[2], nil
}

func ParseReq(reqRaw []byte) (Request, error) {
	reqParts := strings.Split(string(reqRaw), "\r\n")

	if len(reqParts) < 2 {
		return Request{}, fmt.Errorf("invalid request")
	}

	reqLine := reqParts[0]

	method, target, protocol, err := parseReqLine(reqLine)

	if err != nil {
		return Request{}, err
	}

	headersRaw := reqParts[1 : len(reqParts)-1]

	headers := parseHeaders(strings.Join(headersRaw, "\r\n"))

	return Request{
		Method:     method,
		Protocol:   protocol,
		Path:       target,
		Headers:    headers,
		PathParams: RoutePathParams{},
	}, nil
}
