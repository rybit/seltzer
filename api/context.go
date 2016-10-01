package api

import (
	"github.com/labstack/echo"

	"github.com/rybit/config_example/conf"
)

const (
	tokenKey  = "token"
	configKey = "config"
)

func getConfig(c echo.Context) *conf.Config {
	obj := c.Get(configKey)
	return obj.(*conf.Config)
}

func setConfig(c echo.Context, config *conf.Config) {
	c.Set(configKey, config)
}
