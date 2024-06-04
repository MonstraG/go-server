package png

import (
	"bytes"
	"fmt"
	"go-server/helpers"
	"io"
	"net/http"
)

func Validate(w helpers.MyWriter, r *http.Request) {
	ok, err := isValidPNGHeader(r.Body)
	if !ok || err != nil {
	}
	if err != nil {
		w.WriteBadRequest(err.Error())
		return
	}
	if !ok {
		w.WriteBadRequest("png is invalid")
		return
	}

	w.WriteOk("png is valid")
}

const validPNGHeader = "\x89PNG\r\n\x1a\n"

func isValidPNGHeader(body io.ReadCloser) (bool, error) {
	expectedLength := len(validPNGHeader)
	givenHeader := make([]byte, expectedLength)
	readBytes, err := body.Read(givenHeader)
	if err != nil {
		return false, err
	}
	if readBytes != expectedLength {
		return false, fmt.Errorf("expected to read %d bytes for header, but read %d bytes", expectedLength, readBytes)
	}
	return bytes.HasPrefix(givenHeader, []byte(validPNGHeader)), nil
}
