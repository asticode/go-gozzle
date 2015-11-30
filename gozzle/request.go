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
	SetName(n string)
	Method() string
	SetMethod(m string)
	Path() string
	SetPath(p string)
	Headers() map[string]string
	SetHeaders(h map[string]string)
	AddHeader(k string, v string)
	GetHeader(k string) string
	DelHeader(k string)
	Query() map[string]string
	SetQuery(q map[string]string)
	AddQuery(k string, v string)
	GetQuery(k string) string
	DelQuery(k string)
	Body() interface{}
	SetBody(b interface{})
	BodyReader() io.Reader
	SetBodyReader(reader io.Reader)
	BeforeHandler() func(r Request) bool
	SetBeforeHandler(f func(r Request) bool)
	AfterHandler() func(req Request, resp Response)
	SetAfterHandler(f func(req Request, resp Response))
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
func (r *request) SetName(n string) {
	r.name = n
}

// Method returns the request method
func (r *request) Method() string {
	return r.method
}

// SetMethod sets the request method
func (r *request) SetMethod(m string) {
	r.method = m
}

// Path returns the request path
func (r *request) Path() string {
	return r.path
}

// SetPath sets the request path
func (r *request) SetPath(p string) {
	r.path = p
}

// Headers returns the request headers
func (r *request) Headers() map[string]string {
	return r.headers
}

// SetHeaders sets the whole request headers
func (r *request) SetHeaders(h map[string]string) {
	r.headers = array.CloneMap(h)
}

// AddHeader adds a new header for a specific key
func (r *request) AddHeader(k string, v string) {
	r.headers[k] = v
}

// GetHeader returns the value of a specific header key
func (r *request) GetHeader(k string) string {
	return r.headers[k]
}

// DelHeader deletes a specific header key
func (r *request) DelHeader(k string) {
	delete(r.headers, k)
}

// Query returns the whole request query
func (r *request) Query() map[string]string {
	return r.query
}

// SetQuery sets the whole request query
func (r *request) SetQuery(q map[string]string) {
	r.query = array.CloneMap(q)
}

// AddQuery adds a new query for a specific key
func (r *request) AddQuery(k string, v string) {
	r.query[k] = v
}

// GetQuery returns the value of a specific query
func (r *request) GetQuery(k string) string {
	return r.query[k]
}

// DelQuery deletes a specific query key
func (r *request) DelQuery(k string) {
	delete(r.query, k)
}

// Body returns the whole request body
func (r *request) Body() interface{} {
	return r.body
}

// SetBody sets the whole request body
func (r *request) SetBody(b interface{}) {
	r.body = b
}

// BodyReader returns the body reader
func (r *request) BodyReader() io.Reader {
	return r.bodyReader
}

// SetBodyReader sets the body reader
func (r *request) SetBodyReader(reader io.Reader) {
	r.bodyReader = reader
}

// SetBeforeHandler sets the handler executed before sending the request
func (r *request) SetBeforeHandler(f func(r Request) bool) {
	r.beforeHandler = f
}

// BeforeHandler returns the handler executed before sending the request
func (r *request) BeforeHandler() func(r Request) bool {
	return r.beforeHandler
}

// SetAfterHandler sets the handler executed after sending the request
func (r *request) SetAfterHandler(f func(req Request, resp Response)) {
	r.afterHandler = f
}

// SetAfterHandler returns the handler executed after sending the request
func (r *request) AfterHandler() func(req Request, resp Response) {
	return r.afterHandler
}
