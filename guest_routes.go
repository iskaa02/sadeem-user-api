package main

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/iskaa02/sadeem-user-api/api_error"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

func registerGuestRoute(g *echo.Group, db *sqlx.DB) {
	// require guest
	g.POST("/login", func(c echo.Context) error {
		data := map[string]string{}
		if err := c.Bind(&data); err != nil {
			return err
		}
		token, err := login(db, data["username"], data["email"], data["password"])
		if err != nil {
			if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) || errors.Is(err, sql.ErrNoRows) {
				return api_error.NewNotFoundError("invalid_credentials", err)
			}
			return api_error.NewBadRequestError("invalid_login_data", err)
		}
		return c.JSON(http.StatusOK, map[string]string{"token": token})
	})

	g.POST("/register", func(c echo.Context) error {
		data := map[string]string{}
		if err := c.Bind(&data); err != nil {
			return err
		}
		if data["username"] == "" || data["email"] == "" || data["password"] == "" {
			return api_error.NewBadRequestError("invalid_registration_data", errors.New("Empty required data"))
		}
		token, err := register(db, data["username"], data["email"], data["password"])
		if err != nil {
			pgErr, ok := err.(*pq.Error)
			if ok {
				if pgErr.Code == "23505" {
					return api_error.NewBadRequestError("username_or_email_already_exists", err)
				}
			}
			return api_error.NewBadRequestError("", err)
		}
		return c.JSON(http.StatusOK, map[string]string{"token": token})
	})
}
