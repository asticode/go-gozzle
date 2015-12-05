// Copyright 2015, Quentin RENARD. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gozzle

import (
	"io"

	"github.com/asticode/go-toolbox/array"
)

// Request represents a request sendable by gozzle
type Request interface {
	Name() string
	SetName(n string) Request
	Method() string
	SetMethod(m string) Request
	Path() string
	SetPath(p string) Request
	Headers() map[string]string
	SetHeaders(h map[string]string) Request
	AddHeader(k string, v string) Request
	GetHeader(k string) string
	DelHeader(k string) Request
	Query() map[string]string
	SetQuery(q map[string]string) Request
	AddQuery(k string, v string) Request
	GetQuery(k string) string
	DelQuery(k string) Request
	Body() interface{}
	SetBody(b interface{}) Request
	BodyReader() io.Reader
	SetBodyReader(reader io.Reader) Request
	BeforeHandler() func(r Request) bool
	SetBeforeHandler(f func(r Request) bool) Request
	AfterHandler() func(req Request, resp Response)
	SetAfterHandler(f func(req Request, resp Response)) Request
}

// NewRequest creates a new request
func NewRequest(name string, method string, path string) Request {
	return &request{
		name:    name,
		method:  method,
		path:    path,
		headers: make(map[string]string),
		query:   make(map[string]string),
	}
}

type request struct {
	name          string
	method        string
	path          string
	headers       map[string]string
	query         map[string]string
	body          interface{}
	bodyReader    io.Reader
	beforeHandler func(r Request) bool
	afterHandler  func(r Request, resp Response)
}

// Name returns the request name
func (r *request) Name() string {
	return r.name
}

// SetName sets the request name
func (r *request) SetName(n string) Request {
	r.name = n
	return r
}

// Method returns the request method
func (r *request) Method() string {
	return r.method
}

// SetMethod sets the request method
func (r *request) SetMethod(m string) Request {
	r.method = m
	return r
}

// Path returns the request path
func (r *request) Path() string {
	return r.path
}

// SetPath sets the request path
func (r *request) SetPath(p string) Request {
	r.path = p
	return r
}

// Headers returns the request headers
func (r *request) Headers() map[string]string {
	return r.headers
}

// SetHeaders sets the whole request headers
func (r *request) SetHeaders(h map[string]string) Request {
	r.headers = array.CloneMap(h)
	return r
}

// AddHeader adds a new header for a specific key
func (r *request) AddHeader(k string, v string) Request {
	r.headers[k] = v
	return r
}

// GetHeader returns the value of a specific header key
func (r *request) GetHeader(k string) string {
	return r.headers[k]
}

// DelHeader deletes a specific header key
func (r *request) DelHeader(k string) Request {
	delete(r.headers, k)
	return r
}

// Query returns the whole request query
func (r *request) Query() map[string]string {
	return r.query
}

// SetQuery sets the whole request query
func (r *request) SetQuery(q map[string]string) Request {
	r.query = array.CloneMap(q)
	return r
}

// AddQuery adds a new query for a specific key
func (r *request) AddQuery(k string, v string) Request {
	r.query[k] = v
	return r
}

// GetQuery returns the value of a specific query
func (r *request) GetQuery(k string) string {
	return r.query[k]
}

// DelQuery deletes a specific query key
func (r *request) DelQuery(k string) Request {
	delete(r.query, k)
	return r
}

// Body returns the whole request body
func (r *request) Body() interface{} {
	return r.body
}

// SetBody sets the whole request body
func (r *request) SetBody(b interface{}) Request {
	r.body = b
	return r
}

// BodyReader returns the body reader
func (r *request) BodyReader() io.Reader {
	return r.bodyReader
}

// SetBodyReader sets the body reader
func (r *request) SetBodyReader(reader io.Reader) Request {
	r.bodyReader = reader
	return r
}

// SetBeforeHandler sets the handler executed before sending the request
func (r *request) SetBeforeHandler(f func(r Request) bool) Request {
	r.beforeHandler = f
	return r
}

// BeforeHandler returns the handler executed before sending the request
func (r *request) BeforeHandler() func(r Request) bool {
	return r.beforeHandler
}

// SetAfterHandler sets the handler executed after sending the request
func (r *request) SetAfterHandler(f func(req Request, resp Response)) Request {
	r.afterHandler = f
	return r
}

// SetAfterHandler returns the handler executed after sending the request
func (r *request) AfterHandler() func(req Request, resp Response) {
	return r.afterHandler
}
