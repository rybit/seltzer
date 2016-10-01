package api

import (
	"net/http"
	"testing"

	"github.com/labstack/echo/test"
	"github.com/rybit/config_example/conf"
	"github.com/stretchr/testify/assert"
)

func TestCanSayHello(t *testing.T) {
	config := &conf.Config{
		JWTSecret: "secret",
	}
	a := NewAPI(config)

	req := test.NewRequest("GET", "/", nil)
	rsp := test.NewResponseRecorder()
	a.handler.ServeHTTP(req, rsp)

	assert.Equal(t, http.StatusOK, rsp.Status())
}

func TestNoAuthProvided(t *testing.T) {
	config := &conf.Config{
		JWTSecret: "secret",
	}
	a := NewAPI(config)

	req := test.NewRequest("GET", "/private", nil)
	rsp := test.NewResponseRecorder()
	a.handler.ServeHTTP(req, rsp)

	assert.Equal(t, http.StatusBadRequest, rsp.Status())
}

func TestBadAuthProvided(t *testing.T) {
	config := &conf.Config{
		JWTSecret: "secret",
	}
	a := NewAPI(config)

	req := test.NewRequest("GET", "/private", nil)
	req.Header().Add("Authorization", "Bearer nonsense")
	rsp := test.NewResponseRecorder()
	a.handler.ServeHTTP(req, rsp)

	assert.Equal(t, http.StatusUnauthorized, rsp.Status())
}

func TestAuthOk(t *testing.T) {
	config := &conf.Config{
		JWTSecret: "secret",
	}
	a := NewAPI(config)

	req := test.NewRequest("GET", "/private", nil)
	req.Header().Add("Authorization", "Bearer nonsense")
	rsp := test.NewResponseRecorder()
	a.handler.ServeHTTP(req, rsp)

	assert.Equal(t, http.StatusUnauthorized, rsp.Status())
}
