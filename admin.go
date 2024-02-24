package main

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/iskaa02/sadeem-user-api/auth"
	"github.com/jmoiron/sqlx"
)

func registerAdminRoutes(r chi.Router, db *sqlx.DB) {
	r.Use(auth.ChekIsAdminMiddleWare(db, true))
	r.Put("/users/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		u := User{}
		_ = json.NewDecoder(r.Body).Decode(&u)
		updateUser(db, id, &u)
	})
	r.Post("/category", func(w http.ResponseWriter, r *http.Request) {
		c := Category{}
		err := json.NewDecoder(r.Body).Decode(&c)
		if err != nil {
		}
		addCategory(c.Name, c.Activated, db)
	})
	r.Put("/category/{id}", func(w http.ResponseWriter, r *http.Request) {
		c := Category{}
		id := chi.URLParam(r, "id")
		err := json.NewDecoder(r.Body).Decode(&c)
		if err != nil {
		}
		updateCategory(id, c.Name, c.Activated, db)
	})
	r.Delete("/category/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		err := deleteCategory(id, db)
		if err != nil {
			return
		}
	})
	r.Post("/user/{id}/categorize", func(w http.ResponseWriter, r *http.Request) {
		userID := chi.URLParam(r, "id")
		data := map[string]string{}
		json.NewDecoder(r.Body).Decode(&data)
		categorizeUser(db, userID, data["category_id"])
	})

	r.Post("/user/{id}/uncategorize", func(w http.ResponseWriter, r *http.Request) {
		userID := chi.URLParam(r, "id")
		data := map[string]string{}
		json.NewDecoder(r.Body).Decode(&data)
		uncategorizeUser(db, userID, data["category_id"])
	})
}
