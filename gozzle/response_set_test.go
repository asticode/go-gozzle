package gozzle
import (
	"testing"
	"github.com/stretchr/testify/assert"
	"sort"
)

func TestResponseSet(t *testing.T) {
	// Initialize
	respSet := responseSet{
		"1": &response{},
		"2": &response{},
		"3": &response{},
		"4": &response{},
	}

	// Assert names
	e1 := []string{"1","2","3","4"}
	n1 := respSet.Names()
	sort.Strings(n1)
	assert.EqualValues(t, e1, n1)

	// Delete 2
	respSet.DelResponse("2")

	// Assert names
	e2 := []string{"1","3","4"}
	n2 := respSet.Names()
	sort.Strings(n2)
	assert.EqualValues(t, e2, n2)
}
