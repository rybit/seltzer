package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/labstack/echo/test"
	"github.com/rybit/seltzer/conf"
	"github.com/stretchr/testify/assert"
)

var api *API
var testConfig *conf.Config

func TestMain(m *testing.M) {
	testConfig = &conf.Config{
		Port: 0,
	}

	testLogger := logrus.StandardLogger().WithField("testing", true)
	api = NewAPI(testLogger, testConfig)

	os.Exit(m.Run())
}

func TestInfoEndpoint(t *testing.T) {
	code, body := request(t, "GET", "/info", nil)
	if assert.Equal(t, http.StatusOK, code) {
		raw := extractRawPayload(t, body)
		assert.NotEmpty(t, raw["version"])
		assert.NotEmpty(t, raw["description"])
		assert.NotEmpty(t, raw["name"])
	}
}

func TestMissingEndpoint(t *testing.T) {
	code, _ := request(t, "GET", "/missing", nil)
	assert.Equal(t, http.StatusNotFound, code)
}

// ------------------------------------------------------------------------------------------------
// Helpers
// ------------------------------------------------------------------------------------------------

func extractRawPayload(t *testing.T, body string) map[string]interface{} {
	raw := map[string]interface{}{}

	err := json.Unmarshal([]byte(body), &raw)
	assert.NoError(t, err)
	return raw
}

func request(t *testing.T, method, path string, body interface{}) (int, string) {
	req := test.NewRequest(method, path, nil)

	if body != nil {
		bs, err := json.Marshal(body)
		if err != nil {
			assert.FailNow(t, "failed to serialize request body: "+err.Error())
		}

		req = test.NewRequest(method, path, bytes.NewBuffer(bs))
	}

	rsp := test.NewResponseRecorder()
	api.echo.ServeHTTP(req, rsp)
	return rsp.Status(), rsp.Body.String()
}
