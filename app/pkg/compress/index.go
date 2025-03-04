package compress

import (
	"strings"
)

type CompressionType string

const (
	GZIP CompressionType = "gzip"
)

type CompressionHandler = func([]byte) ([]byte, error)

var compressMethods = map[CompressionType]CompressionHandler{
	GZIP: gzip,
}

func GetCompressHandler(compressType string) CompressionHandler {
	handler, exists := compressMethods[CompressionType(strings.ToLower(compressType))]
	if !exists {
		return nil
	}

	return handler
}
