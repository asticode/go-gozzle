// Copyright 2015, Quentin RENARD. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gozzle

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExec(t *testing.T) {
	// Initialize
	n := 5
	f := "test %d"
	h := "Test"

	// Create server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Loop
		for i := 0; i < n; i++ {
			// Get formatted message
			m := fmt.Sprintf(f, i)

			// Valid header
			if r.Header.Get(h) == m {
				// Set header
				w.Header().Set(h, m)

				// Set body
				w.Write([]byte(m))
			}
		}
	}))
	defer server.Close()

	// Create request set
	reqSet := NewRequestSet()

	// Loop
	for i := 0; i < n; i++ {
		// Get formatted message
		m := fmt.Sprintf(f, i)

		// Create request
		req := NewRequest(m, MethodGet, server.URL)
		req.AddHeader(h, m)
		reqSet.AddRequest(req)
	}

	// Create gozzle
	g := NewGozzle(0)

	// Execute requests
	respSet := g.Exec(reqSet)

	// Loop
	assert.Len(t, respSet.Names(), n)
	for i := 0; i < n; i++ {
		// Get formatted message
		m := fmt.Sprintf(f, i)

		// Get response
		resp := respSet.GetResponse(m)

		// Assert
		assert.Len(t, resp.Errors(), 0)
		assert.Equal(t, m, resp.Header().Get(h))
		c, e := ioutil.ReadAll(resp.BodyReader())
		assert.NoError(t, e)
		assert.Equal(t, m, string(c))
	}
}

func TestExecRequestBeforeHandler(t *testing.T) {
	// Create server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	defer server.Close()

	// Create request set
	reqSet := NewRequestSet()

	// Create request
	req := NewRequest("test", MethodGet, server.URL)
	req.SetBeforeHandler(func(r Request) bool {
		return false
	})
	reqSet.AddRequest(req)

	// Create gozzle
	g := NewGozzle(0)

	// Execute requests
	respSet := g.Exec(reqSet)

	// Assert
	assert.Len(t, respSet.Names(), 0)
}

func TestExecRequestError(t *testing.T) {
	// Initialize
	c := 500

	// Create server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(c)
	}))
	defer server.Close()

	// Create request set
	reqSet := NewRequestSet()

	// Create requests
	req1 := NewRequest("test1", MethodGet, server.URL)
	reqSet.AddRequest(req1)
	req2 := NewRequest("test2", MethodGet, "test")
	reqSet.AddRequest(req2)

	// Create gozzle
	g := NewGozzle(0)

	// Execute requests
	respSet := g.Exec(reqSet)

	// Assert
	assert.Len(t, respSet.Names(), 2)
	assert.Equal(t, c, respSet.GetResponse("test1").StatusCode())
	assert.Len(t, respSet.GetResponse("test1").Errors(), 1)
	assert.EqualError(t, respSet.GetResponse("test1").Errors()[0], ErrInvalidStatusCode.Error())
	assert.Len(t, respSet.GetResponse("test2").Errors(), 1)
	assert.EqualError(t, respSet.GetResponse("test2").Errors()[0], "Get test: unsupported protocol scheme \"\"")
}

func TestQuery(t *testing.T) {
	// Initialize
	r1 := request{
		query: map[string]string{
			"a":     "b",
			"ké@lù": "ùl@ék",
		},
	}
	r2 := request{}

	// Assert
	assert.Contains(t, "?a=b&k%C3%A9%40l%C3%B9=%C3%B9l%40%C3%A9k", query(&r1))
	assert.Empty(t, query(&r2))
}

func TestBody(t *testing.T) {
	// Initialize
	r := request{
		body: map[string]string{
			"test": "message",
		},
		headers: map[string]string{
			"Content-Type": "application/json",
		},
	}

	// Get body reader
	b, e := body(&r)
	assert.NoError(t, e)

	// Read body
	c, e := ioutil.ReadAll(b)
	assert.NoError(t, e)

	// Assert
	assert.Equal(t, "{\"test\":\"message\"}", string(c))
}

func TestHeaders(t *testing.T) {
	// Initialize
	k := "Key"
	v := "Value"
	r := NewRequest("test", MethodGet, "/test")
	r.AddHeader(k, v)
	hr := http.Request{Header: http.Header{}}

	// Assert
	assert.Empty(t, hr.Header.Get(k))

	// Add headers
	headers(r, &hr)

	// Assert
	assert.Equal(t, v, hr.Header.Get(k))
}
