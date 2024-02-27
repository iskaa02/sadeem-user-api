package main

import (
	"net/http"
	"strconv"

	"github.com/iskaa02/sadeem-user-api/auth"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

func registerUserRoute(g *echo.Group, db *sqlx.DB) {
	// require auth
	g.GET("/users/me", func(c echo.Context) error {
		id, _ := c.Get(auth.UserIDContextKey).(string)
		isAdmin, _ := c.Get(auth.IsAdminContextKey).(bool)
		u, err := getUser(db, id, isAdmin)
		if err != nil {
			return err
		}
		return c.JSON(200, u)
	}, auth.RequireAuthMiddleWare)
	g.PUT("/users/me", func(c echo.Context) error {
		id := c.Get(auth.UserIDContextKey).(string)
		u := &User{}
		if err := c.Bind(&u); err != nil {
			return err
		}
		return updateUser(db, id, u)
	}, auth.RequireAuthMiddleWare)
	g.POST("/users/change_password", func(c echo.Context) error {
		id := c.Get(auth.UserIDContextKey).(string)
		data := ChangePasswordParams{}
		if err := c.Bind(&data); err != nil {
			return err
		}
		return changePassword(db, id, data.OldPassword, data.NewPassword)
	}, auth.RequireAuthMiddleWare)
	g.POST("/users/change_image", func(c echo.Context) error {
		id, _ := c.Get(auth.UserIDContextKey).(string)
		return updateImage(db, c.Request(), id)
	})

	g.GET("/category", func(c echo.Context) error {
		isAdmin := c.Get(auth.IsAdminContextKey).(bool)
		page, _ := strconv.Atoi(c.Request().URL.Query().Get("page"))
		page -= 1
		if page < 0 {
			page = 0
		}
		searchQuery := c.Request().URL.Query().Get("q")
		result, err := listCategories(db, isAdmin, page, searchQuery)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, result)
	})
}
