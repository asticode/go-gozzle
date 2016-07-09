// Copyright 2015, Quentin RENARD. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gozzle

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
)

// Variables
var (
	ErrInvalidStatusCode   = errors.New("Invalid status code")
	ErrNilOriginalResponse = errors.New("Nil original response")
)

// Response represents a response received by gozzle after sending a request
type Response interface {
	Errors() []error
	Status() string
	StatusCode() int
	Header() http.Header
	BodyReader() io.ReadCloser
	Body() ([]byte, error)
	Close() error
}

// NewResponseError creates a new response with an error set by default
func NewResponseError(e error) Response {
	// Create response
	r := response{}

	// Add error
	r.errors = append(r.errors, e)

	// Return
	return &r
}

// NewResponse creates a new response based on an *http.Response
func NewResponse(or *http.Response, maxSizeBody int) Response {
	// Initialize
	r := response{
		originalResponse: or,
	}

	// Check status code
	if r.StatusCode() < 200 || r.StatusCode() >= 300 {
		r.errors = append(r.errors, ErrInvalidStatusCode)
	}

	// Update body reader
	if maxSizeBody > 0 {
		r.originalResponse.Body = ioutil.NopCloser(
			io.LimitReader(r.originalResponse.Body, int64(maxSizeBody)),
		)
	}

	// Return
	return &r
}

type response struct {
	errors           []error
	originalResponse *http.Response
}

// Error returns the response error
func (r *response) Errors() []error {
	return r.errors
}

// Status returns the response status text
func (r *response) Status() string {
	if r.originalResponse == nil {
		return ""
	}
	return r.originalResponse.Status
}

// StatusCode returns the status code
func (r *response) StatusCode() int {
	if r.originalResponse == nil {
		return 0
	}
	return r.originalResponse.StatusCode
}

// Header returns the http.Header object of the http.Response
func (r *response) Header() http.Header {
	if r.originalResponse == nil {
		return http.Header{}
	}
	return r.originalResponse.Header
}

// BodyReader returns the response BodyReader
func (r *response) BodyReader() io.ReadCloser {
	if r.originalResponse == nil {
		return ioutil.NopCloser(bytes.NewReader([]byte{}))
	}
	return r.originalResponse.Body
}

// Body returns the response body without compromising the BodyReader
func (r *response) Body() ([]byte, error) {
	var b []byte
	if r.originalResponse == nil {
		return b, nil
	}
	c, err := ioutil.ReadAll(r.originalResponse.Body)
	if err != nil {
		return b, err
	}
	r.originalResponse.Body = ioutil.NopCloser(bytes.NewReader(c))
	return c, nil
}

// Close closes the response
func (r *response) Close() error {
	if r.originalResponse == nil {
		return ErrNilOriginalResponse
	}
	return r.originalResponse.Body.Close()
}
