package api

import (
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	"github.com/rybit/config_example/conf"
)

// API is the root of the api. Methods and actions are tied to the api as makes sense
type API struct {
	handler *echo.Echo
	log     *logrus.Entry
	config  *conf.Config
}

// JWTClaims define what we are expecting when we deserialize the token
type JWTClaims struct {
	jwt.StandardClaims
	Name  string `json:"name"`
	Admin bool   `json:"is_admin"`
}

// NewAPI will create an API based on the configuration
func NewAPI(config *conf.Config) *API {
	api := &API{
		log:    logrus.StandardLogger().WithField("component", "api"),
		config: config,
	}

	requireToken := middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte(config.JWTSecret),
		Claims:     &JWTClaims{},
	})

	requireAdmin := func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// TODO ~ verify that 'is_admin' is enabled
			return next(c)
		}
	}

	e := echo.New()

	// unauthenticated
	e.Get("/", hello)

	// can add auth on a single endpoint
	e.Get("/admin", adminOnly, requireAdmin)

	// can add auth to a group of endpoints
	authenticatedEndpoints := e.Group("/", requireToken)
	authenticatedEndpoints.GET("private", authenticated)

	api.handler = e
	return api
}

func hello(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"application": "seltzer", "version": "1.0"})
}

func authenticated(c echo.Context) error {
	return c.String(200, "wtf")
}

func adminOnly(c echo.Context) error {
	return c.String(200, "def an admin")
}
