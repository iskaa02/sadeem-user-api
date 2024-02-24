package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/iskaa02/sadeem-user-api/auth"
	"github.com/jmoiron/sqlx"
)

func main() {
	db, err := sqlx.Connect("postgres", "user=foo dbname=bar sslmode=disable")
	if err != nil {
		panic(err)
	}
	r := chi.NewMux()
	r.Route("/api/user", func(r chi.Router) {
		registerUserRoute(r, db)
	})
	r.Route("/api/", func(r chi.Router) {
		registerAdminRoutes(r, db)
	})

	r.Get("/api/category", func(w http.ResponseWriter, r *http.Request) {
		isAdmin := r.Context().Value(auth.IsAdminContextKey).(bool)
		page, _ := strconv.Atoi(r.URL.Query().Get("page"))
		page -= 1
		if page < 0 {
			page = 0
		}
		result := listCategories(db, isAdmin, page)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(result)
	})
}
