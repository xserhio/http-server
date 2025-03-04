package compress

import (
	"bytes"
	compressGzip "compress/gzip"
)

func gzip(data []byte) ([]byte, error) {
	var b bytes.Buffer
	gz := compressGzip.NewWriter(&b)

	if _, err := gz.Write(data); err != nil {
		return nil, err
	}

	if err := gz.Close(); err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}
