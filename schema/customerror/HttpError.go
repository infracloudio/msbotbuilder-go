package customerror

import (
	"bytes"
	"fmt"
	"io"
)

// HTTPError wraps a raw HTTP error
type HTTPError struct {
	StatusCode int
	HtErr      error
	Body       io.ReadCloser
}

func (htErr HTTPError) Error() string {
	buf := new(bytes.Buffer)
	buf.ReadFrom(htErr.Body)
	return fmt.Sprintf("HTTP error %d: %s %s", htErr.StatusCode, htErr.HtErr, buf.String())
}
