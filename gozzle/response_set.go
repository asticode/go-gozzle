// Copyright 2015, Quentin RENARD. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gozzle

import "sync"

// ResponseSet represents a set of responses
type ResponseSet interface {
	Names() []string
	AddResponse(req Request, resp Response) ResponseSet
	GetResponse(name string) Response
	DelResponse(name string) ResponseSet
	Close() map[string]error
}

// NewResponseSet creates a new response set
func NewResponseSet() ResponseSet {
	return &responseSet{
		responses: make(map[string]Response),
	}
}

type responseSet struct {
	responses map[string]Response
	mutex     sync.Mutex
}

// Responses returns the list of names
func (respSet *responseSet) Names() []string {
	var n []string
	for k := range respSet.responses {
		n = append(n, k)
	}
	return n
}

// AddResponse adds a new response to the response set
func (respSet *responseSet) AddResponse(req Request, resp Response) ResponseSet {
	respSet.mutex.Lock()
	respSet.responses[req.Name()] = resp
	respSet.mutex.Unlock()
	return respSet
}

// GetResponse returns a request based on its name
func (respSet *responseSet) GetResponse(name string) Response {
	return respSet.responses[name]
}

// DelResponse removes a request from the request set
func (respSet *responseSet) DelResponse(name string) ResponseSet {
	delete(respSet.responses, name)
	return respSet
}

// Close closes the responses in the response set
func (respSet *responseSet) Close() map[string]error {
	errors := make(map[string]error)
	for k, v := range respSet.responses {
		errors[k] = v.Close()
	}
	return errors
}
