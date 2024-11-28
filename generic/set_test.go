package generic

import (
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringSet(t *testing.T) {
	set := NewSet[string]()
	a := "a"
	b := "b"
	c := "c"
	set.Add(a, b)

	assert.True(t, set.Contains(a))
	assert.True(t, set.Contains(b))
	assert.False(t, set.Contains(c))

	set.Remove(b)
	assert.False(t, set.Contains(b))

	another := NewSet[string]()
	another.Add(a, b, c)
	set.Merge(another)
	assert.True(t, set.Contains(a))
	assert.True(t, set.Contains(b))
	assert.True(t, set.Contains(c))

	testStr := ""
	f := func(s string) bool {
		testStr += s
		return true
	}

	set.Enumerate(f)
	slice := strings.Split(testStr, "")
	sort.Strings(slice)
	assert.True(t, strings.Join(slice, "") == a+b+c)
}
