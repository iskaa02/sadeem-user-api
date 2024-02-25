package main

import (
	"net/http"
	"strconv"

	"github.com/iskaa02/sadeem-user-api/api_error"
	"github.com/iskaa02/sadeem-user-api/auth"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
)

func main() {
	db, err := sqlx.Connect("postgres", "postgresql://postgres:password@localhost:5432?sslmode=disable")
	if err != nil {
		panic(err)
	}
	e := echo.New()

	e.HTTPErrorHandler = api_error.GlobalErrorHandler

	e.Use(auth.LoadToken(db))
	registerAdminRoutes(e.Group("/api", auth.RequireAdmin), db)
	registerUserRoute(e.Group("/api/users", auth.RequireAuthMiddleWare), db)
	registerGuestRoute(e.Group("/api", auth.RequireNoAuthMiddleWare), db)

	// anyone can see
	e.GET("/api/category", func(c echo.Context) error {
		isAdmin := c.Get(auth.IsAdminContextKey).(bool)
		page, _ := strconv.Atoi(c.Request().URL.Query().Get("page"))
		page -= 1
		if page < 0 {
			page = 0
		}
		result := listCategories(db, isAdmin, page)
		return c.JSON(http.StatusOK, result)
	})
	e.Start("localhost:3000")
}
