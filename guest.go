package main

import (
	"errors"
	"net/http"
	"net/mail"

	"github.com/iskaa02/sadeem-user-api/api_error"
	"github.com/iskaa02/sadeem-user-api/auth"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
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
			return api_error.NewBadRequestError("", err)
		}
		return c.JSON(http.StatusOK, map[string]string{"token": token})
	})
}

func login(db *sqlx.DB, username, email, password string) (string, error) {
	user := &User{}
	var err error = nil
	if username == "" && email == "" {
		err = api_error.NewBadRequestError("missing_both_email_and_username", errors.New("Username and email can't be both empty"))
	}
	if email != "" {
		err = db.Get(user, "SELECT * FROM users WHERE email=$1", email)
	} else {
		err = db.Get(user, "SELECT * FROM users WHERE username=$1", username)
	}
	if err != nil {
		return "", err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password)); err != nil {
		return "", err
	}
	token := auth.GenerateToken(user.ID)
	return token, nil
}

func register(db *sqlx.DB, username, email, password string) (string, error) {
	id := ""
	_, err := mail.ParseAddress(email)
	if err != nil {
		return "", api_error.NewBadRequestError("invalid_email", err)
	}
	hashed_password, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", api_error.NewBadRequestError("", err)
	}
	err = db.Get(&id, "INSERT INTO users(username,email,hashed_password) VALUES($1,$2,$3) RETURNING id", username, email, hashed_password)
	if err != nil {
		return "", err
	}
	token := auth.GenerateToken(id)
	return token, nil
}
