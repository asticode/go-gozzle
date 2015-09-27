package gozzle

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
)

type ResponseSet map[string]*Response

func (oResponseSet *ResponseSet) addResponse(oRequest *Request, oHttpResponse *http.Response, oGozzle Gozzle) {
	// Initialize.
	oResponse := Response{}
	oResponse.originalResponse = oHttpResponse
	oResponse.errorCode = ERROR_NONE

	// Check status code
	if oHttpResponse.StatusCode < 200 || oHttpResponse.StatusCode >= 300 {
		oResponse.errorCode = ERROR_INVALID_STATUS_CODE
	}

	// Get body
	oReader := oHttpResponse.Body
	if oGozzle.client.maxSizeBody > 0 {
		oReader = struct {
			io.Reader
			io.Closer
		}{io.LimitReader(oReader, int64(oGozzle.client.maxSizeBody)), oReader}
	}
	oBody, oErr := ioutil.ReadAll(oReader)
	oResponse.contentLength = len(oBody)
	if oErr != nil {
		oHttpResponse.Body = struct {
			io.Reader
			io.Closer
		}{bytes.NewReader(oBody), oHttpResponse.Body}
	}
	oResponse.body = &oBody

	// Add response
	(*oResponseSet)[oRequest.name] = &oResponse
}
