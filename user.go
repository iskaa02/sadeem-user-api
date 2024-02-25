package main

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/guregu/null"
	"github.com/iskaa02/sadeem-user-api/api_error"
	"github.com/iskaa02/sadeem-user-api/auth"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
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
		u, err := getUser(db, id)
		fmt.Println(err)
		if err != nil {
			return err
		}
		return c.JSON(200, u)
	}, auth.RequireAuthMiddleWare)
	g.PUT("/me", func(c echo.Context) error {
		id := c.Get(auth.UserIDContextKey).(string)
		u := &User{}
		c.Bind(&u)
		return updateUser(db, id, u)
	}, auth.RequireAuthMiddleWare)
	g.POST("/me/change_password", func(c echo.Context) error {
		id := c.Get(auth.UserIDContextKey).(string)
		data := ChangePasswordParams{}
		c.Bind(&data)
		return changePassword(db, id, data.OldPassword, data.NewPassword)
	}, auth.RequireAuthMiddleWare)
	g.POST("/me/change_image", func(c echo.Context) error {
		id, _ := c.Get(auth.UserIDContextKey).(string)
		return updateImage(db, c.Request(), id)
	})

	// anyone can see
}

func getUser(db *sqlx.DB, id string) (User, error) {
	u := User{}
	err := db.Get(&u, "SELECT id,username,email,image_path FROM users WHERE id=$1", id)
	if err != nil {
		if err == sql.ErrNoRows {
			return u, api_error.NewNotFoundError("user_not_found", err)
		}
	}
	return u, err
}

func updateUser(db *sqlx.DB, id string, u *User) error {
	_, err := db.Exec("UPDATE users SET username=$1,email=$2 WHERE id=$3", u.Username, u.Email, id)
	if err != nil {
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
	fmt.Println(err)
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

func listUsers(db *sqlx.DB, page int) ([]User, error) {
	users := []User{}
	err := db.Select(&users, "SELECT * FROM users ORDER BY username LIMIT 10 OFFSET $1 ", page*10)
	if err != nil {
		return users, err
	}
	return users, err
}

func categorizeUser(db *sqlx.DB, userID, categoryID string) error {
	_, err := db.Exec("INSERT INTO user_category(user_id,category_id) VALUES($1,$2)", userID, categoryID)
	return err
}

func uncategorizeUser(db *sqlx.DB, userID, categoryID string) error {
	_, err := db.Exec("DELETE FROM user_category WHERE user_id=$1 AND category_id=$2", userID, categoryID)
	return err
}
