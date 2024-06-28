package http

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewRequest_BuildURLFailure(t *testing.T) {
	c := NewClient(ClientConfig{
		Host: ".f:://g",
	}, nil)
	_, err := c.NewRequest("GET", "/test", nil)
	require.Error(t, err)
	require.Equal(t, "build request url: parse \".f:://g\": first path segment in URL cannot contain colon", err.Error())
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
	reqURI := "/test?test-list=b&test-list=a&example-1=12"

	req, err := http.NewRequest("GET", reqURI, nil)
	require.NoError(t, err)

	params := parseQueryParams(req)
	require.Equal(t, 2, len(params))
	require.Equal(t, "example-1=12", params[0])
	require.Equal(t, "test-list=a,b", params[1])
}
