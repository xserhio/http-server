package http

import (
	"regexp"
	"strings"
)

type Headers map[string]string

func isHeaderNameValid(inputString string) bool {
	pattern := `[^a-zA-Z0-9_-]`

	re := regexp.MustCompile(pattern)

	return !re.MatchString(inputString)
}

func parseHeaders(headersRaw string) Headers {
	headers := make(Headers, 10)

	for _, l := range strings.Split(headersRaw, "\r\n") {
		parts := strings.Split(l, ":")

		if len(parts) != 2 {
			continue
		}

		name := strings.TrimSpace(parts[0])
		name = strings.ToLower(name)
		value := strings.TrimSpace(parts[1])

		if len(name) > 15 {
			continue
		}

		if len(value) > 255 {
			continue
		}

		if !isHeaderNameValid(name) {
			continue
		}

		headers[name] = value
	}

	return headers
}
