package http

import (
	"net/http"
	"strings"
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

func TestParseJSONParams(t *testing.T) {
	reqBody := `{
		"name": "demo1",
		"ids": [1, 2, 3],
		"deviceInfo": {
			"id": 1
		},
		"deviceList": [
			{
				"id": 1
			},
			{
				"id": 2
			}
		]
	}`
	req, err := http.NewRequest("POST", "/test", strings.NewReader(reqBody))
	require.NoError(t, err)

	req.Header.Add("Content-Type", "application/json")

	params, err := parseJSONParams(req)
	require.NoError(t, err)

	paramsStr := strings.Join(params, "&")
	require.Equal(t, "deviceInfo.id=1&deviceList[0].id=1&deviceList[1].id=2&ids[0]=1&ids[1]=2&ids[2]=3&name=demo1", paramsStr)
}

func TestParseJSONParams_EmptyRequest(t *testing.T) {
	req, err := http.NewRequest("GET", "/test", nil)
	require.NoError(t, err)

	req.Header.Add("Content-Type", "application/json")

	_, err = parseJSONParams(req)
	require.Error(t, err)
	require.Equal(t, "invalid request: must contain a non-nil body since content-type is JSON", err.Error())
}

func TestDo_ErrParseJSON(t *testing.T) {
	c := NewClient(ClientConfig{Host: "http://localhost:8080"}, nil)

	req, err := http.NewRequest("GET", "/test", strings.NewReader("{>"))
	require.NoError(t, err)

	req.Header.Add("Content-Type", "application/json")

	_, err = c.Do(req)
	require.Error(t, err)
	require.Equal(t, "parse request parameters: parse JSON request body: invalid character '>' looking for beginning of object key string", err.Error())
}
