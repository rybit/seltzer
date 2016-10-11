package api

import (
	"fmt"

	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"

	"github.com/rybit/seltzer/conf"
)

// API is the data holder for the API
type API struct {
	log    *logrus.Entry
	config *conf.Config
	echo   *echo.Echo
}

// Start will start the API on the specified port
func (api *API) Start() error {
	return api.echo.Run(standard.New(fmt.Sprintf(":%d", api.config.Port)))
}

// Stop will shutdown the engine internally
func (api *API) Stop() error {
	return api.echo.Stop()
}

// NewAPI will create an api instance that is ready to start
func NewAPI(log *logrus.Entry, config *conf.Config) *API {
	// create the api
	api := &API{
		config: config,
		log:    log.WithField("component", "api"),
	}

	// add the endpoints
	e := echo.New()
	e.Get("/info", api.Info)

	e.SetHTTPErrorHandler(api.handleError)
	e.SetLogger(wrapper{api.log})
	api.echo = e

	return api
}

func (api *API) Info(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, map[string]string{
		"version":     "testing",
		"description": "a boiler plate project",
		"name":        "seltzer",
	})
}
