package auth

import (
	"context"
	"net/http"
	"strings"

	"github.com/jmoiron/sqlx"
)

var (
	IsAdminContextKey = "admin"
	UserIDContextKey  = "user"
)

// HTTP middleware setting a value on the request context
func ChekIsAdminMiddleWare(db *sqlx.DB, required bool) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("Authorization")
			token = strings.Trim(token, "Bearer")
			id, isAdmin := isAdmin(token, db)
			if !isAdmin && required {
				w.WriteHeader(http.StatusForbidden)
				// w.Write(string)
				return
			}
			ctx := context.WithValue(r.Context(), IsAdminContextKey, isAdmin)
			ctx = context.WithValue(ctx, UserIDContextKey, id)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func RequireAuthMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		token = strings.Trim(token, "Bearer")
		id, isAuth := isAuthenticated(token)
		if !isAuth {
			w.WriteHeader(http.StatusForbidden)
			// w.Write(string)
			return
		}
		ctx := context.WithValue(r.Context(), UserIDContextKey, id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func RequireNoAuthMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		token = strings.Trim(token, "Bearer")
		id, isAuth := isAuthenticated(token)
		if isAuth {
			w.WriteHeader(http.StatusOK)
			// w.Write(string)
			return
		}
		ctx := context.WithValue(r.Context(), UserIDContextKey, id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
