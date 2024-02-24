package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/iskaa02/sadeem-user-api/auth"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID             string `json:"id" db:"id"`
	Username       string `json:"username" db:"username"`
	Email          string `json:"email" db:"email"`
	Image          string `json:"image" db:"image"`
	HashedPassword string `json:"_" db:"hashed_password"`
}

func registerUserRoute(r chi.Router, db *sqlx.DB) {
	// require auth
	r.Group(func(authRouter chi.Router) {
		authRouter.Use(auth.RequireAuthMiddleWare)
		authRouter.Get("users/me", func(w http.ResponseWriter, r *http.Request) {
			id := r.Context().Value(auth.UserIDContextKey).(string)
			getUser(db, id)
		})
		authRouter.Put("users/me", func(w http.ResponseWriter, r *http.Request) {
			id := r.Context().Value(auth.UserIDContextKey).(string)
			u := &User{}
			updateUser(db, id, u)
		})
		authRouter.Post("users/me/change_password", func(w http.ResponseWriter, r *http.Request) {
			id := r.Context().Value(auth.UserIDContextKey).(string)
			oldPassword := ""
			newPassowrd := ""
			changePassword(db, id, oldPassword, newPassowrd)
		})
	})
	// require guest
	r.Group(func(noAuth chi.Router) {
		noAuth.Post("/login", func(w http.ResponseWriter, r *http.Request) {
			data := map[string]string{}
			json.NewDecoder(r.Body).Decode(&data)
			token, err := login(db, data["email"], data["username"], data["password"])
			if err != nil {
				w.WriteHeader(400)
				return
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{token: token})
		})

		noAuth.Post("/register", func(w http.ResponseWriter, r *http.Request) {
			data := map[string]string{}
			json.NewDecoder(r.Body).Decode(&data)
			token, err := register(db, data["username"], data["email"], data["password"])
			if err != nil {
				w.WriteHeader(400)
				return
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{token: token})
		})
	})
	// anyone can see
	r.Get("/users/list", func(w http.ResponseWriter, r *http.Request) {
		page, _ := strconv.Atoi(r.URL.Query().Get("page"))
		page--
		if page < 0 {
			page = 0
		}
		listUsers(db, page)
	})
}

func login(db *sqlx.DB, username, email, password string) (string, error) {
	user := &User{}
	var err error = nil
	if username == "" && email == "" {
		err = errors.New("Username and email cannot be empty")
	}
	if email == "" {
		err = db.Select(&user, "SELECT * FROM user WHERE email=$1", email)
	} else {
		err = db.Select(&user, "SELECT * FROM user WHERE username=$1", username)
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
	err := db.Get(&id, "INSERT INTO user(username,email,password) VALUES($1,$2,$3) RETURNING id", username, email, password)
	if err != nil {
		return "", err
	}
	token := auth.GenerateToken(id)
	return token, nil
}

func getUser(db *sqlx.DB, id string) User {
	u := User{}
	db.Select(&u, "SELECT id,username,email FROM user WHERE id=$1", id)
	return u
}

func updateUser(db *sqlx.DB, id string, u *User) error {
	_, err := db.Exec("UPDATE user SET username=$1,email=$2")
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
	err := db.Select(&users, "SELECT * FROM user LIMIT 10 OFFSET=$1 SORT BY name", page*10)
	if err != nil {
		return err
	}
	return err
}

func categorizeUser(db *sqlx.DB, userID, categoryID string) {
	db.Exec("INSERT INTO user_category(userID,categoryID) VALUES($1,$2)", userID, categoryID)
}

func uncategorizeUser(db *sqlx.DB, userID, categoryID string) {
	db.Exec("DELETE FROM user_category WHERE userID=$1 AND categoryID=$2", userID, categoryID)
}
