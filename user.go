package main

import (
	"database/sql"
	"errors"
	"net/http"
	"net/mail"

	"github.com/guregu/null"
	"github.com/huandu/go-sqlbuilder"
	"github.com/iskaa02/sadeem-user-api/api_error"
	"github.com/iskaa02/sadeem-user-api/auth"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID             string      `json:"id" db:"id"`
	Username       string      `json:"username" db:"username"`
	Email          string      `json:"email" db:"email"`
	ImagePath      null.String `json:"image_path" db:"image_path"`
	HashedPassword string      `json:"-" db:"hashed_password"`
}
type ChangePasswordParams struct {
	NewPassword string `json:"new_password"`
	OldPassword string `json:"old_password"`
}

func registerUserRoute(g *echo.Group, db *sqlx.DB) {
	// require auth
	g.GET("/me", func(c echo.Context) error {
		id, _ := c.Get(auth.UserIDContextKey).(string)
		isAdmin, _ := c.Get(auth.IsAdminContextKey).(bool)
		u, err := getUser(db, id, isAdmin)
		if err != nil {
			return err
		}
		return c.JSON(200, u)
	}, auth.RequireAuthMiddleWare)
	g.PUT("/me", func(c echo.Context) error {
		id := c.Get(auth.UserIDContextKey).(string)
		u := &User{}
		if err := c.Bind(&u); err != nil {
			return err
		}
		return updateUser(db, id, u)
	}, auth.RequireAuthMiddleWare)
	g.POST("/me/change_password", func(c echo.Context) error {
		id := c.Get(auth.UserIDContextKey).(string)
		data := ChangePasswordParams{}
		if err := c.Bind(&data); err != nil {
			return err
		}
		return changePassword(db, id, data.OldPassword, data.NewPassword)
	}, auth.RequireAuthMiddleWare)
	g.POST("/me/change_image", func(c echo.Context) error {
		id, _ := c.Get(auth.UserIDContextKey).(string)
		return updateImage(db, c.Request(), id)
	})

	// anyone can see
}

type GetUserRes struct {
	User
	Categories []Category `json:"category"`
}

func getUser(db *sqlx.DB, id string, isAdmin bool) (GetUserRes, error) {
	u := GetUserRes{}
	err := db.Get(&u, "SELECT id,username,email,image_path FROM users WHERE id=$1", id)
	if err != nil {
		if err == sql.ErrNoRows {
			return u, api_error.NewNotFoundError("user_not_found", err)
		}
	}
	category, err := getUserCategory(db, id, isAdmin)
	if err != nil {
		return u, err
	}
	u.Categories = category
	return u, err
}

func updateUser(db *sqlx.DB, id string, u *User) error {
	sb := sqlbuilder.Update("users")
	sb.Where(sb.EQ("id", id))
	if u.Email != "" {
		_, err := mail.ParseAddress(u.Email)
		if err != nil {
			return api_error.NewBadRequestError("invalid_email", err)
		}
		sb.SetMore(sb.Assign("email", u.Email))
	}
	if u.Username != "" {
		sb.SetMore(sb.Assign("username", u.Username))
	}
	sql, args := sb.Build()
	_, err := db.Exec(sql, args...)
	if err != nil {
		pgErr, ok := err.(*pq.Error)
		if ok {
			if pgErr.Code == "23505" {
				return api_error.NewBadRequestError("username_or_email_already_exists", err)
			}
		}
		return err
	}
	return err
}

func updateImage(db *sqlx.DB, r *http.Request, id string) error {
	filename, err := uploadFile(r, id)
	if err != nil {
		return err
	}
	_, err = db.Exec("UPDATE users SET image_path=$1 WHERE id=$2", filename, id)
	if err != nil {
		return err
	}
	return nil
}

func changePassword(db *sqlx.DB, id, oldPassword, newPassowrd string) error {
	u := &User{}
	err := db.Get(u, "SELECT * FROM users WHERE id=$1", id)
	if err != nil {
		return err
	}
	err = bcrypt.CompareHashAndPassword([]byte(u.HashedPassword), []byte(oldPassword))
	if err != nil {
		return api_error.NewBadRequestError("old_password_do_not_match", errors.New("old password is not correct"))
	}
	newHashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassowrd), 10)
	_, err = db.Exec("UPDATE users SET hashed_password=$2 WHERE id=$1", id, newHashedPassword)
	if err != nil {
		return err
	}
	return nil
}

func listUsers(db *sqlx.DB, page int, searchQuery string) ([]User, error) {
	users := []User{}
	var err error
	sb := sqlbuilder.Select("*").
		From("users").
		Limit(10).
		Offset(page * 10)

	if searchQuery != "" {
		sb.Where(sb.Like("username", "%"+searchQuery+"%"))
	}
	sql, args := sb.Build()
	err = db.Select(&users, sql, args...)
	return users, err
}
