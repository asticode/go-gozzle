// Copyright 2015, Quentin RENARD. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gozzle

// ResponseSet represents a set of responses
type ResponseSet interface {
	Names() []string
	AddResponse(req Request, resp Response)
	GetResponse(name string) Response
	DelResponse(name string)
}

// NewResponseSet creates a new response set
func NewResponseSet() ResponseSet {
	return &responseSet{}
}

type responseSet map[string]Response

// Responses returns the list of names
func (respSet *responseSet) Names() []string {
	var n []string
	for k := range *respSet {
		n = append(n, k)
	}
	return n
}

// AddResponse adds a new response to the response set
func (respSet *responseSet) AddResponse(req Request, resp Response) {
	(*respSet)[req.Name()] = resp
}

// GetResponse returns a request based on its name
func (respSet *responseSet) GetResponse(name string) Response {
	return (*respSet)[name]
}

// DelResponse removes a request from the request set
func (respSet *responseSet) DelResponse(name string) {
	delete((*respSet), name)
}
