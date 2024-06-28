//go:build integration

package http

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

var clientConfig ClientConfig

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

	os.Exit(m.Run())
}

func TestListDevices(t *testing.T) {
	c := NewClient(clientConfig, http.DefaultClient)

	req, err := c.NewRequest("GET", "/iot-open/sign/device/list", nil)
	require.NoError(t, err)

	resp, err := c.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	data, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var dataMap map[string]interface{}
	err = json.Unmarshal(data, &dataMap)
	require.NoError(t, err)

	respCode, cast := dataMap["code"].(string)
	require.True(t, cast)
	require.Equal(t, "0", respCode)
}

func TestGetAllQuotaValues(t *testing.T) {
	// TODO
}

func TestGetXBoostSwitch(t *testing.T) {
	// TODO
}
