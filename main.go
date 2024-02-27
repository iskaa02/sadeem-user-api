package main

import (
	"os"

	"github.com/huandu/go-sqlbuilder"
	"github.com/iskaa02/sadeem-user-api/api_error"
	"github.com/iskaa02/sadeem-user-api/auth"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
)

func main() {
	db, err := sqlx.Connect("postgres", os.Getenv("DB_SOURCE_NAME"))
	if err != nil {
		panic(err)
	}
	e := echo.New()
	sqlbuilder.DefaultFlavor = sqlbuilder.PostgreSQL
	e.HTTPErrorHandler = api_error.GlobalErrorHandler

	e.Use(auth.LoadToken(db))
	registerAdminRoutes(e.Group("/api", auth.RequireAdmin), db)
	registerUserRoute(e.Group("/api", auth.RequireAuthMiddleWare), db)
	registerGuestRoute(e.Group("/api", auth.RequireNoAuthMiddleWare), db)

	e.Static("/images", "./images")
	e.Start("localhost:3000")
}
