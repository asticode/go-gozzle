// Copyright 2015, Quentin RENARD. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gozzle

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"log/syslog"
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
}

// Configuration represents a JSON-friendly gozzle configuration
type Configuration struct {
	MaxSizeBody int `json:"max_size_body"`
}

// NewGozzle creates a new Gozzle object
func NewGozzle(maxSizeBody int) Gozzle {
	return &gozzle{
		maxSizeBody: maxSizeBody,
		client:      &http.Client{},
	}
}

// NewGozzleFromConfiguration creates a new Gozzle object based on a configuration
func NewGozzleFromConfiguration(c Configuration) Gozzle {
	return NewGozzle(
		c.MaxSizeBody,
	)
}

type gozzle struct {
	maxSizeBody int
	client      *http.Client
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
		req := reqSet.GetRequest(name)
		go func() {
			// Execute request
			resp := g.execRequest(req)

			// Add response
			if resp != nil {
				respSet.AddResponse(req, resp)
			}

			// Update wait group
			wg.Done()
		}()
	}

	// Wait
	wg.Wait()

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
	if e != nil {
		return NewResponseError(e)
	}

	// Create http request
	httpReq, e := http.NewRequest(
		req.Method(),
		req.Path()+query(req),
		b,
	)
	if e != nil {
		return NewResponseError(e)
	}

	// Add headers
	headers(req, httpReq)

	// TODO Add timeout and context

	// TODO Remove
	w, e := syslog.New(syslog.LOG_ERR, "bite")
	defer w.Close()
	w.Err(fmt.Sprintf("Sending HTTP Verb <%s>", httpReq.Method))

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

func query(r Request) string {
	var query string
	if len(r.Query()) > 0 {
		// Add "?"
		query += "?"

		// Loop through query parameters
		for k, v := range r.Query() {
			query += url.QueryEscape(k) + "=" + url.QueryEscape(v) + "&"
		}

		// Trim "&"
		query = strings.Trim(query, "&")
	}
	return query
}

func body(r Request) (io.Reader, error) {
	// Initialize
	var body []byte
	var e error

	// Encode body
	if r.GetHeader("Content-Type") == "application/xml" {
		// XML marshall
		body, e = xml.Marshal(r.Body())
	} else {
		// JSON marshall
		body, e = json.Marshal(r.Body())
	}

	// Return
	return bytes.NewBuffer(body), e
}

func headers(r Request, hr *http.Request) {
	// Loop through headers
	for k, v := range r.Headers() {
		hr.Header.Set(k, v)
	}
}
