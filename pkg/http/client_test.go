package http

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRequest_BuildURLFailure(t *testing.T) {
	// TODO
}

func TestAsciiCompare(t *testing.T) {
	// identical
	cmp := asciiCompare("abc", "abc")
	assert.Equal(t, 0, cmp)

	// bigger ascii
	cmp = asciiCompare("abc1", "abc!")
	assert.Equal(t, 1, cmp)

	// bigger ascii and shorter string
	cmp = asciiCompare(":", ")*123")
	assert.Equal(t, 1, cmp)

	// smaller ascii
	cmp = asciiCompare("abc)12", "abc1!!")
	assert.Equal(t, -1, cmp)

	// smaller ascii and longer string
	cmp = asciiCompare("/123", "123")
	assert.Equal(t, -1, cmp)

	// equal ascii and longer string
	cmp = asciiCompare("1233", "123")
	assert.Equal(t, 1, cmp)

	// equal ascii and shorter string
	cmp = asciiCompare("ab", "abc")
	assert.Equal(t, -1, cmp)
}

func TestParseQueryParams_ListValue(t *testing.T) {
	// TODO
}

func TestDo_InvalidSecretKey(t *testing.T) {
	// TODO
}
