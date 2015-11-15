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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewResponseError(t *testing.T) {
	// Initialize
	resp := NewResponseError(errors.New("test"))

	// Assert
	assert.Equal(t, 1, len(resp.Errors()))
}

func TestNewResponseInvalidStatusCode(t *testing.T) {
	// Initialize
	httpRespSuccess := NewResponse(&http.Response{StatusCode: 200}, 0)
	httpRespError := NewResponse(&http.Response{StatusCode: 400}, 0)

	// Assert
	assert.Len(t, httpRespSuccess.Errors(), 0)
	assert.Len(t, httpRespError.Errors(), 1)
	assert.EqualError(t, httpRespError.Errors()[0], ErrInvalidStatusCode.Error())
}

type mockedCloser struct{}

func (m mockedCloser) Close() error { return nil }

func mockedIoReaderCloser(b []byte) io.ReadCloser {
	return struct {
		io.Reader
		io.Closer
	}{bytes.NewReader(b), mockedCloser{}}
}

func TestNewResponseBodyReader(t *testing.T) {
	// Initialize
	b := []byte("testmessage")
	max := 3

	// New responses
	resp1 := NewResponse(&http.Response{Body: mockedIoReaderCloser(b)}, 0)
	resp2 := NewResponse(&http.Response{Body: mockedIoReaderCloser(b)}, max)

	// Read bodies
	b1, e := ioutil.ReadAll(resp1.BodyReader())
	assert.NoError(t, e)
	b2, e := ioutil.ReadAll(resp2.BodyReader())
	assert.NoError(t, e)

	// Assert
	assert.Len(t, b1, len(b))
	assert.Len(t, b2, max)
}
