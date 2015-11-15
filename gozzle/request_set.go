// Copyright 2015, Quentin RENARD. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gozzle

// RequestSet represents a set of sendable requests
type RequestSet interface {
	Names() []string
	AddRequest(r Request)
	GetRequest(name string) Request
	DelRequest(name string)
}

// NewRequestSet creates a new request set
func NewRequestSet() RequestSet {
	return &requestSet{}
}

type requestSet map[string]Request

// Requests returns the list of names
func (reqSet *requestSet) Names() []string {
	var n []string
	for k := range *reqSet {
		n = append(n, k)
	}
	return n
}

// AddRequest adds a new request to the request set
func (reqSet *requestSet) AddRequest(r Request) {
	(*reqSet)[r.Name()] = r
}

// GetRequest returns a request based on its name
func (reqSet *requestSet) GetRequest(name string) Request {
	return (*reqSet)[name]
}

// DelRequest removes a request from the request set
func (reqSet *requestSet) DelRequest(name string) {
	delete((*reqSet), name)
}
