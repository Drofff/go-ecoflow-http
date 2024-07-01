//go:build integration

package http

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	clientConfig ClientConfig
	testDeviceSN string
)

func TestMain(m *testing.M) {
	ak, exists := os.LookupEnv("TEST_ACCESS_KEY")
	if !exists {
		log.Fatalln("must set TEST_ACCESS_KEY env var")
	}

	sk, exists := os.LookupEnv("TEST_SECRET_KEY")
	if !exists {
		log.Fatalln("must set TEST_SECRET_KEY env var")
	}

	clientConfig = ClientConfig{
		Host:      "https://api-e.ecoflow.com",
		AccessKey: ak,
		SecretKey: sk,
	}

	sn, exists := os.LookupEnv("TEST_DEVICE_SN")
	if exists {
		testDeviceSN = sn
	}

	os.Exit(m.Run())
}

func parseResponse(t *testing.T, resp *http.Response) map[string]any {
	require.Equal(t, http.StatusOK, resp.StatusCode)

	data, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var dataMap map[string]any
	err = json.Unmarshal(data, &dataMap)
	require.NoError(t, err)

	return dataMap
}

func assertResponseCode(t *testing.T, data map[string]any) {
	respCode, cast := data["code"].(string)
	require.True(t, cast)
	require.Equal(t, "0", respCode)
}

func TestListDevices(t *testing.T) {
	c := NewClient(clientConfig, http.DefaultClient)

	req, err := c.NewRequest("GET", "/iot-open/sign/device/list", nil)
	require.NoError(t, err)

	resp, err := c.Do(req)
	require.NoError(t, err)

	data := parseResponse(t, resp)
	assertResponseCode(t, data)
}

func TestGetAllQuotaValues(t *testing.T) {
	require.NotEmpty(t, testDeviceSN)

	c := NewClient(clientConfig, http.DefaultClient)

	req, err := c.NewRequest("GET", "/iot-open/sign/device/quota/all", nil)
	require.NoError(t, err)

	queryParams := req.URL.Query()
	queryParams.Set("sn", testDeviceSN)
	req.URL.RawQuery = queryParams.Encode()

	resp, err := c.Do(req)
	require.NoError(t, err)

	data := parseResponse(t, resp)
	assertResponseCode(t, data)
}

func TestGetXBoostSwitch(t *testing.T) {
	require.NotEmpty(t, testDeviceSN)

	c := NewClient(clientConfig, http.DefaultClient)

	reqJSON := fmt.Sprintf("{\"sn\":\"%v\",\"params\":{\"quotas\":[\"inv.cfgAcEnabled\",\"inv.cfgAcXboost\"]}}", testDeviceSN)
	req, err := c.NewRequest("POST", "/iot-open/sign/device/quota", strings.NewReader(reqJSON))
	require.NoError(t, err)

	req.Header.Add(headerContentType, contentTypeJSON)

	resp, err := c.Do(req)
	require.NoError(t, err)

	data := parseResponse(t, resp)
	assertResponseCode(t, data)
}

func TestTurnOnAC(t *testing.T) {
	require.NotEmpty(t, testDeviceSN)

	c := NewClient(clientConfig, http.DefaultClient)

	reqJSON := fmt.Sprintf(`
{
	"sn": "%v",
	"params": {
		"cmdSet": 32,
		"id": 66,
		"enabled": 1, 
		"xboost": 1
	}
}
	`, testDeviceSN)
	req, err := c.NewRequest("PUT", "/iot-open/sign/device/quota", strings.NewReader(reqJSON))
	require.NoError(t, err)

	req.Header.Add(headerContentType, contentTypeJSON)

	resp, err := c.Do(req)

	data := parseResponse(t, resp)
	assertResponseCode(t, data)
}
