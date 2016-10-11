package api

import (
	"net/http"

	"github.com/labstack/echo"
)

func (api *API) handleError(err error, ctx echo.Context) {
	if ctx.Response().Committed() {
		return
	}

	httpErr, ok := err.(*echo.HTTPError)
	if ok {
		ctx.JSON(httpErr.Code, httpErr)
	} else {
		api.log.WithError(err).Warn("Unexpected non http error")
		ctx.JSON(http.StatusInternalServerError, &echo.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Internal server error",
		})
	}
}
