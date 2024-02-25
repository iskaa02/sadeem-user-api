package auth

import (
	"errors"
	"net/http"
	"strings"

	"github.com/iskaa02/sadeem-user-api/api_error"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

var (
	IsAdminContextKey = "admin"
	UserIDContextKey  = "user"
)

// HTTP middleware setting a value on the request context
func LoadToken(db *sqlx.DB) func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return echo.HandlerFunc(func(c echo.Context) error {
			token := c.Request().Header.Get("Authorization")
			if token != "" {
				token = strings.Trim(token, "Bearer")
				id, isAdmin := isAdmin(token, db)
				c.Set(IsAdminContextKey, isAdmin)
				c.Set(UserIDContextKey, id)
			}
			return next(c)
		})
	}
}

func RequireAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return echo.HandlerFunc(func(c echo.Context) error {
		isAdmin, _ := c.Get(IsAdminContextKey).(bool)
		if !isAdmin {
			return api_error.NewForbiddenError("", errors.New(""))
		}
		return next(c)
	})
}

func RequireAuthMiddleWare(next echo.HandlerFunc) echo.HandlerFunc {
	return echo.HandlerFunc(func(c echo.Context) error {
		userID, ok := c.Get(UserIDContextKey).(string)
		if userID == "" || !ok {
			return api_error.NewUnauthorizedError("", errors.New("User not authnticated"))
		}
		return next(c)
	})
}

func RequireNoAuthMiddleWare(next echo.HandlerFunc) echo.HandlerFunc {
	return echo.HandlerFunc(func(c echo.Context) error {
		userID, _ := c.Get(UserIDContextKey).(string)
		if userID != "" {
			return c.NoContent(http.StatusOK)
		}
		return next(c)
	})
}
