package api_error

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func GlobalErrorHandler(err error, c echo.Context) {
	apiError, errorCasted := err.(ApiError)
	// fallback to default error handler
	if !errorCasted {
		fmt.Println(err)
		c.Echo().DefaultHTTPErrorHandler(err, c)
		return
	}
	if !c.Response().Committed {
		if c.Request().Method == http.MethodHead {
			c.NoContent(apiError.Code)
			return
		}
		i18n, ok := translated_errors[apiError.MessageCode]
		if !ok {
			i18n = translated_errors["something_went_wrong"]
		}
		msg := struct {
			Msg string `json:"message"`
		}{Msg: i18n.Translate(c.Request().Header.Get("Accept-Language"))}
		err = c.JSON(apiError.Code, msg)

		if err != nil {
			fmt.Println(err)
		}
	}
}
