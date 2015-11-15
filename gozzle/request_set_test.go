// Copyright 2015, Quentin RENARD. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gozzle

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequestSet(t *testing.T) {
	// Initialize
	reqSet := requestSet{
		"1": &request{},
		"2": &request{},
		"3": &request{},
		"4": &request{},
	}

	// Assert names
	e1 := []string{"1", "2", "3", "4"}
	n1 := reqSet.Names()
	sort.Strings(n1)
	assert.EqualValues(t, e1, n1)

	// Delete 2
	reqSet.DelRequest("2")

	// Assert names
	e2 := []string{"1", "3", "4"}
	n2 := reqSet.Names()
	sort.Strings(n2)
	assert.EqualValues(t, e2, n2)
}
