package gozzle

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQuerySorted(t *testing.T) {
	// Initialize
	p := "/"
	r1 := request{
		path: p,
		query: map[string]string{
			"2": "b",
			"1": "a",
			"3": "c",
			"5": "e",
			"4": "d",
		},
	}
	r2 := request{}

	// Assert
	assert.Equal(t, "/?1=a&2=b&3=c&4=d&5=e", r1.FullPath())
	assert.Empty(t, r2.FullPath())
}

func TestQueryEncoding(t *testing.T) {
	// Initialize
	r1 := request{
		query: map[string]string{
			"a":     "b",
			"ké@lù": "ùl@ék",
		},
	}
	r2 := request{}

	// Assert
	assert.Equal(t, "?a=b&k%C3%A9%40l%C3%B9=%C3%B9l%40%C3%A9k", r1.FullPath())
	assert.Empty(t, r2.FullPath())
}
