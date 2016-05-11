// Copyright 2015, Quentin RENARD. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gozzle

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"io"
	"io/ioutil"
	"net/http"
	"sync"
	"github.com/rs/xlog"
	"time"
	"fmt"
)

// Constants
const (
	MethodGet     string = "GET"
	MethodPost    string = "POST"
	MethodPut     string = "PUT"
	MethodPatch   string = "PATCH"
	MethodDelete  string = "DELETE"
	MethodOptions string = "OPTIONS"
	MethodHead    string = "HEAD"
)

// Gozzle represents an object capable of executing a set of requests
type Gozzle interface {
	Exec(reqSet RequestSet) ResponseSet
	ExecWithLogger(reqSet RequestSet, f xlog.F) ResponseSet
	MaxSizeBody() int
	SetMaxSizeBody(maxSizeBody int) Gozzle
}

// Configuration represents a JSON-friendly gozzle configuration
type Configuration struct {
	MaxSizeBody int `json:"max_size_body"`
}

// NewGozzle creates a new Gozzle object
func NewGozzle() Gozzle {
	return &gozzle{
		client: &http.Client{},
	}
}

// NewGozzleFromConfiguration creates a new Gozzle object based on a configuration
func NewGozzleFromConfiguration(c Configuration) Gozzle {
	return NewGozzle().SetMaxSizeBody(c.MaxSizeBody)
}

type gozzle struct {
	maxSizeBody int
	client      *http.Client
}

func (g *gozzle) SetMaxSizeBody(maxSizeBody int) Gozzle {
	g.maxSizeBody = maxSizeBody
	return g
}

func (g *gozzle) MaxSizeBody() int {
	return g.maxSizeBody
}

// Exec executes a set of requests
func (g gozzle) Exec(reqSet RequestSet) ResponseSet {
	// Initialize
	respSet := NewResponseSet()
	reqNames := reqSet.Names()

	// Create wait group
	wg := sync.WaitGroup{}
	wg.Add(len(reqNames))

	// Loop through requests
	for _, name := range reqNames {
		go func(req Request) {
			// Execute request
			resp := g.execRequest(req)

			// Add response
			if resp != nil {
				respSet.AddResponse(req, resp)
			}

			// Update wait group
			wg.Done()
		}(reqSet.GetRequest(name))
	}

	// Wait
	wg.Wait()

	// Return
	return respSet
}

// ExecWithLogger executes a set of requests and logs requests duration
func (g gozzle) ExecWithLogger(reqSet RequestSet, f xlog.F) ResponseSet {
	// Log total duration
	defer func(t time.Time) {
		f["duration_gozzle_total"] = time.Since(t)
	}(time.Now())
	n := time.Now()

	// Initialize
	respSet := NewResponseSet()
	reqNames := reqSet.Names()

	// Create wait group
	wg := sync.WaitGroup{}
	wg.Add(len(reqNames))
	f["duration_gozzle_init_wg"] = time.Since(n)
	n = time.Now()

	// Loop through requests
	for _, name := range reqNames {
		go func(req Request) {
			// Execute request
			resp := g.execRequestWithLogger(req, f)
			now := time.Now()

			f[fmt.Sprintf("duration_gozzle_request_executed_%s", req.Name())] = time.Since(now)
			now = time.Now()

			// Add response
			if resp != nil {
				respSet.AddResponse(req, resp)
			}

			f[fmt.Sprintf("duration_gozzle_response_added_%s", req.Name())] = time.Since(now)
			now = time.Now()

			// Update wait group
			wg.Done()
			f[fmt.Sprintf("duration_gozzle_done_wait_%s", req.Name())] = time.Since(now)
		}(reqSet.GetRequest(name))
	}

	f["duration_gozzle_start_wait_wg"] = time.Since(n)
	n = time.Now()

	// Wait
	wg.Wait()

	f["duration_gozzle_done_wait_wg"] = time.Since(n)

	// Return
	return respSet
}

func (g gozzle) execRequest(req Request) Response {
	// Before handler
	if req.BeforeHandler() != nil {
		cont := req.BeforeHandler()(req)
		if !cont {
			return nil
		}
	}

	// Get body
	b, e := body(req)
	defer b.Close()
	if e != nil {
		return NewResponseError(e)
	}

	// Create http request
	httpReq, e := http.NewRequest(
		req.Method(),
		req.FullPath(),
		b,
	)
	if e != nil {
		return NewResponseError(e)
	}
	httpReq.Close = true

	// Add headers
	headers(req, httpReq)

	// TODO Add timeout and context

	// Send request
	httpResp, e := g.client.Do(httpReq)
	if e != nil {
		return NewResponseError(e)
	}

	// Create response
	resp := NewResponse(httpResp, g.maxSizeBody)

	// After handler
	if req.AfterHandler() != nil {
		req.AfterHandler()(req, resp)
	}

	// Return
	return resp
}

func (g gozzle) execRequestWithLogger(req Request, f xlog.F) Response {
	// Init
	now := time.Now()

	// Before handler
	if req.BeforeHandler() != nil {
		cont := req.BeforeHandler()(req)
		if !cont {
			return nil
		}
	}

	f[fmt.Sprintf("duration_gozzle_before_handler_%s", req.Name())] = time.Since(now)
	now = time.Now()

	// Get body
	b, e := body(req)
	defer b.Close()
	if e != nil {
		return NewResponseError(e)
	}

	f[fmt.Sprintf("duration_gozzle_body_%s", req.Name())] = time.Since(now)
	now = time.Now()

	// Create http request
	httpReq, e := http.NewRequest(
		req.Method(),
		req.FullPath(),
		b,
	)
	if e != nil {
		return NewResponseError(e)
	}
	httpReq.Close = true

	f[fmt.Sprintf("duration_gozzle_create_request_%s", req.Name())] = time.Since(now)
	now = time.Now()

	// Add headers
	headers(req, httpReq)

	f[fmt.Sprintf("duration_gozzle_headers_%s", req.Name())] = time.Since(now)
	now = time.Now()

	// TODO Add timeout and context

	// Send request
	httpResp, e := g.client.Do(httpReq)
	if e != nil {
		return NewResponseError(e)
	}

	f[fmt.Sprintf("duration_gozzle_do_request_%s", req.Name())] = time.Since(now)
	now = time.Now()

	// Create response
	resp := NewResponse(httpResp, g.maxSizeBody)

	f[fmt.Sprintf("duration_gozzle_create_response_%s", req.Name())] = time.Since(now)
	now = time.Now()

	// After handler
	if req.AfterHandler() != nil {
		req.AfterHandler()(req, resp)
	}

	f[fmt.Sprintf("duration_gozzle_after_handler_%s", req.Name())] = time.Since(now)

	// Return
	return resp
}

func body(r Request) (io.ReadCloser, error) {
	// Initialize
	var body []byte
	var e error

	// Encode body
	if r.BodyReader() != nil {
		bodyReader, ok := r.BodyReader().(io.ReadCloser)
		if !ok {
			bodyReader = ioutil.NopCloser(r.BodyReader())
		}
		return bodyReader, e
	} else if r.Body() != nil {
		// Get body reader
		if r.GetHeader("Content-Type") == "application/xml" {
			// XML marshall
			body, e = xml.Marshal(r.Body())
		} else {
			// JSON marshall
			body, e = json.Marshal(r.Body())
		}
	}

	// Return
	return ioutil.NopCloser(bytes.NewBuffer(body)), e
}

func headers(r Request, hr *http.Request) {
	// Loop through headers
	for k, v := range r.Headers() {
		hr.Header.Set(k, v)
	}
}
