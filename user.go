package main

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/iskaa02/sadeem-user-api/auth"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID             string `json:"id" db:"id"`
	Name           string `json:"name" db:"name"`
	Email          string `json:"email" db:"email"`
	Image          string `json:"image" db:"image"`
	HashedPassword string `json:"_" db:"hashed_password"`
}

func registerUserRoute(r chi.Router, db *sqlx.DB) {
	// require auth
	r.Group(func(authRouter chi.Router) {
		authRouter.Use(auth.RequireAuthMiddleWare)
		authRouter.Get("/me", func(w http.ResponseWriter, r *http.Request) {
			id := r.Context().Value(auth.UserContextKey).(string)
			getUser(db, id)
		})
		authRouter.Put("/me", func(w http.ResponseWriter, r *http.Request) {
			id := r.Context().Value(auth.UserContextKey).(string)
			u := &User{}
			updateUser(db, id, u)
		})
		authRouter.Post("/change_password", func(w http.ResponseWriter, r *http.Request) {
			id := r.Context().Value(auth.UserContextKey).(string)
			oldPassword := ""
			newPassowrd := ""
			changePassword(db, id, oldPassword, newPassowrd)
		})
	})
	r.Get("/list", func(w http.ResponseWriter, r *http.Request) {
		page, _ := strconv.Atoi(r.URL.Query().Get("page"))
		page--
		if page < 0 {
			page = 0
		}
		listUsers(db, page)
	})
}

func login(db *sqlx.DB, usernameOrEmail, password string) (string, error) {
	user := &User{}
	err := db.Select(&user, "SELECT * FROM user WHERE email=$1", usernameOrEmail)
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
	err := db.Get(&id, "INSERT INTO user(username,email,password) VALUES($1,$2,$3) RETURNING id", username, email, password)
	if err != nil {
		return "", err
	}
	token := auth.GenerateToken(id)
	return token, nil
}

func getUser(db *sqlx.DB, id string) {
	u := &User{}
	db.Select(u, "SELECT id,name,email FROM user WHERE id=$1", id)
}

func updateUser(db *sqlx.DB, id string, u *User) error {
	_, err := db.Exec("UPDATE user SET name=$1,email=$2")
	if err != nil {
		return err
	}
	return nil
}

func updateImage(db *sqlx.DB, r *http.Request, id string) error {
	filename, err := uploadFile(r, id)
	if err != nil {
		return err
	}
	_, err = db.Exec("UPDATE user SET image=$2 WHERE id=$1", id, filename)
	if err != nil {
		return err
	}
	return nil
}

func changePassword(db *sqlx.DB, id, oldPassword, newPassowrd string) error {
	u := &User{}
	err := db.Select(u, "SELECT * FROM user WERE id=$1", id)
	if err != nil {
		return err
	}
	err = bcrypt.CompareHashAndPassword([]byte(u.HashedPassword), []byte(oldPassword))
	if err != nil {
		return err
	}
	newHashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassowrd), 10)
	_, err = db.Exec("UPDATE user SET hashed_password=$2 WHERE id=$1", id, newHashedPassword)
	if err != nil {
		return err
	}
	return nil
}

func listUsers(db *sqlx.DB, page int) error {
	users := []User{}
	err := db.Select(&users, "SELECT * FROM user LIMIT 10 OFFSET=$1", page*10)
	if err != nil {
		return err
	}
	return err
}
