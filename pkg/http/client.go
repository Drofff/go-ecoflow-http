package http

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	neturl "net/url"
	"slices"
	"strings"
	"time"
)

// Client wraps http.Client handling the EcoFlow API requests configuration like auth, signing parameters
// and setting the target host URL - all based on the pre-provided ClientConfig.
type Client interface {
	// NewRequest wraps http.NewRequest() setting the pre-configured host URL into the created request.
	NewRequest(method, url string, body io.Reader) (*http.Request, error)
	// Do wraps http.Client#Do() inserting auth/signature headers into the request.
	Do(req *http.Request) (*http.Response, error)
}

type ClientConfig struct {
	// Host is the target EcoFlow OpenPlatform API URL f.e. "https://api-e.ecoflow.com".
	Host string
	// AccessKey issued to your application.
	AccessKey string
	// SecretKey issued to your application.
	SecretKey string
}

type client struct {
	conf       ClientConfig
	httpClient *http.Client
}

type signature struct {
	hash      string
	nonce     string
	timestamp string
}

const (
	paramAccessKey = "accessKey"
	paramNonce     = "nonce"
	paramTimestamp = "timestamp"

	headerAccessKey   = paramAccessKey
	headerNonce       = paramNonce
	headerTimestamp   = paramTimestamp
	headerSignature   = "sign"
	headerContentType = "content-type"

	contentTypeJSON = "application/json"
)

func NewClient(conf ClientConfig, httpClient *http.Client) Client {
	return &client{
		conf:       conf,
		httpClient: httpClient,
	}
}

func (c *client) NewRequest(method, url string, body io.Reader) (*http.Request, error) {
	url, err := neturl.JoinPath(c.conf.Host, url)
	if err != nil {
		return nil, fmt.Errorf("build request url: %w", err)
	}
	return http.NewRequest(method, url, body)
}

func toParamStr(k, v string) string {
	return fmt.Sprintf("%v=%v", k, v)
}

func asciiCompare(a, b string) int {
	if a == b {
		return 0
	}

	minSize := min(len(a), len(b))
	for i := 0; i < minSize; i++ {
		if a[i] < b[i] {
			return -1
		}
		if a[i] > b[i] {
			return 1
		}
	}

	if minSize == len(a) {
		return -1
	}
	return 1
}

func hasJSONBody(req *http.Request) bool {
	ct := req.Header.Get(headerContentType)
	return strings.ToLower(ct) == contentTypeJSON
}

func parseJSONParams(req *http.Request) []string {
	// TODO
	return nil
}

func parseQueryParams(req *http.Request) []string {
	var params []string

	q := req.URL.Query()
	for k, vs := range q {
		if len(vs) > 1 {
			slices.SortFunc(vs, asciiCompare)
		}

		v := strings.Join(vs, ",")
		params = append(params, toParamStr(k, v))
	}

	slices.SortFunc(params, asciiCompare)
	return params
}

func parseParams(req *http.Request) []string {
	if hasJSONBody(req) {
		return parseJSONParams(req)
	}
	return parseQueryParams(req)
}

func newNonce() string {
	return fmt.Sprint(rand.Intn(999999))
}

func newTimestamp() string {
	return fmt.Sprint(time.Now().UTC().UnixMilli())
}

func (c *client) calcSignature(params []string) (signature, error) {
	n := newNonce()
	t := newTimestamp()

	params = append(params,
		toParamStr(paramAccessKey, c.conf.AccessKey),
		toParamStr(paramNonce, n),
		toParamStr(paramTimestamp, t))
	payload := strings.Join(params, "&")

	h := hmac.New(sha256.New, []byte(c.conf.SecretKey))
	_, err := h.Write([]byte(payload))
	if err != nil {
		return signature{}, fmt.Errorf("hash payload: %w", err)
	}

	hash := hex.EncodeToString(h.Sum([]byte{}))

	return signature{
		hash:      hash,
		nonce:     n,
		timestamp: t,
	}, nil
}

func (c *client) Do(req *http.Request) (*http.Response, error) {
	params := parseParams(req)
	sign, err := c.calcSignature(params)
	if err != nil {
		return nil, fmt.Errorf("calculate signature: %w", err)
	}

	req.Header.Add(headerAccessKey, c.conf.AccessKey)
	req.Header.Add(headerNonce, sign.nonce)
	req.Header.Add(headerTimestamp, sign.timestamp)
	req.Header.Add(headerSignature, sign.hash)

	return c.httpClient.Do(req)
}
