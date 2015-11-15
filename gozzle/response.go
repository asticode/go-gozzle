// Copyright 2015, Quentin RENARD. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gozzle

import (
	"errors"
	"io"
	"net/http"
)

// Variables
var (
	ErrInvalidStatusCode = errors.New("Invalid status code")
)

// Response represents a response received by gozzle after sending a request
type Response interface {
	Errors() []error
	Status() string
	StatusCode() int
	Header() http.Header
	BodyReader() io.ReadCloser
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
		maxSizeBody:      maxSizeBody,
	}

	// Check status code
	if r.StatusCode() < 200 || r.StatusCode() >= 300 {
		r.errors = append(r.errors, ErrInvalidStatusCode)
	}

	// Return
	return &r
}

type response struct {
	errors           []error
	maxSizeBody      int
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
	return r.originalResponse.Header
}

// Body returns the response body
func (r *response) BodyReader() io.ReadCloser {
	if r.maxSizeBody > 0 {
		return struct {
			io.Reader
			io.Closer
		}{io.LimitReader(r.originalResponse.Body, int64(r.maxSizeBody)), r.originalResponse.Body}
	}
	return r.originalResponse.Body
}
