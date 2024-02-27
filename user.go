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
