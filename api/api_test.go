package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
	"github.com/labstack/echo/test"
	"github.com/rybit/seltzer/conf"
	"github.com/stretchr/testify/assert"
)

var api *API
var testConfig *conf.Config

func TestMain(m *testing.M) {
	testConfig = &conf.Config{
		Port:      0,
		JWTSecret: "secret",
	}

	testLogger := logrus.StandardLogger().WithField("testing", true)
	api = NewAPI(testLogger, testConfig)

	os.Exit(m.Run())
}

func TestInfoEndpoint(t *testing.T) {
	code, body := request(t, "GET", "/info", nil)
	if assert.Equal(t, http.StatusOK, code) {
		raw := make(map[string]string)
		extractPayload(t, body, &raw)
		assert.NotEmpty(t, raw["version"])
		assert.NotEmpty(t, raw["description"])
		assert.NotEmpty(t, raw["name"])
	}
}

func TestMissingEndpoint(t *testing.T) {
	code, body := request(t, "GET", "/missing", nil)
	assert.Equal(t, http.StatusNotFound, code)
	err := extractError(t, body)
	assert.Equal(t, http.StatusNotFound, err.Code)
	assert.NotEmpty(t, err.Message)
}

func TestGenerateWithBadPayload(t *testing.T) {
	payload := &TokenRequest{
		Email: "batman@dc.com",
		Pass:  "",
	}
	code, body := request(t, "POST", "/login", payload)
	if assert.Equal(t, http.StatusBadRequest, code) {
		err := extractError(t, body)
		assert.Equal(t, http.StatusBadRequest, err.Code)
		assert.NotEmpty(t, err.Message)
	}

	payload.Email = ""
	payload.Pass = "some-magic-string"
	code, body = request(t, "POST", "/login", payload)
	if assert.Equal(t, http.StatusBadRequest, code) {
		err := extractError(t, body)
		assert.Equal(t, http.StatusBadRequest, err.Code)
		assert.NotEmpty(t, err.Message)
	}
}

func TestGenerateNewToken(t *testing.T) {
	code, body := request(t, "POST", "/login", &TokenRequest{
		Email: "batman@dc.com",
		Pass:  "some-magic-string",
	})
	if assert.Equal(t, http.StatusCreated, code) {
		rsp := new(TokenResponse)
		extractPayload(t, body, rsp)

		assert.NotEmpty(t, rsp.Key)
	}
}



// ------------------------------------------------------------------------------------------------
// Helpers
// ------------------------------------------------------------------------------------------------

func extractError(t *testing.T, body string) *echo.HTTPError {
	raw := new(echo.HTTPError)
	extractPayload(t, body, raw)
	return raw
}

func extractPayload(t *testing.T, body string, out interface{}) {
	err := json.Unmarshal([]byte(body), out)
	assert.NoError(t, err)
}

func request(t *testing.T, method, path string, body interface{}) (int, string) {
	req := test.NewRequest(method, path, nil)

	if body != nil {
		bs, err := json.Marshal(body)
		if err != nil {
			assert.FailNow(t, "failed to serialize request body: "+err.Error())
		}

		req = test.NewRequest(method, path, bytes.NewBuffer(bs))
		req.Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	}

	rsp := test.NewResponseRecorder()
	api.echo.ServeHTTP(req, rsp)
	return rsp.Status(), rsp.Body.String()
}
