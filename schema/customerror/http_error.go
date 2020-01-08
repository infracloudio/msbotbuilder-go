package customerror

import (
	"bytes"
	"fmt"
	"io"
)

type HttpError struct {
	StatusCode int
	HtErr      error
	Body       io.ReadCloser
}

func (htErr HttpError) Error() string {
	buf := new(bytes.Buffer)
	buf.ReadFrom(htErr.Body)
	return fmt.Sprintf("HTTP error %d: %s %s", htErr.StatusCode, htErr.HtErr, buf.String())
}
