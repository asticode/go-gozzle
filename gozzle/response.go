package gozzle
import (
	"bytes"
	"io"
	"net/http"
)

type Response struct {
	originalResponse *http.Response
	errorCode        Error
	contentLength    int
	body             *[]byte
}

type Error int

const (
	ERROR_NONE Error = iota
	ERROR_INVALID_STATUS_CODE
)

func (oResponse Response) IsError() bool {
	if oResponse.errorCode != ERROR_NONE {
		return true
	}
	return false
}

func (oResponse Response) ErrorCode() Error {
	return oResponse.errorCode
}

func (oResponse Response) ContentLength() int {
	return oResponse.contentLength
}

func (oResponse Response) Body() *[]byte {
	return oResponse.body
}

func (oResponse Response) BodyReader() io.Reader {
	return struct {
		io.Reader
	}{bytes.NewReader(*oResponse.body)}
}

func (oResponse Response) OriginalResponse() *http.Response {
	return oResponse.originalResponse
}

func (oResponse *Response) SetBody(oBody *[]byte) {
	oResponse.body = oBody
}
