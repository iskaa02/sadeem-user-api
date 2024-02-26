package main

import (
	"net/http"
	"strconv"

	"github.com/iskaa02/sadeem-user-api/auth"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
)

func registerAdminRoutes(g *echo.Group, db *sqlx.DB) {
	g.PUT("/users/:id", func(c echo.Context) error {
		id := ""
		echo.PathParamsBinder(c).String("id", &id)
		u := User{}
		c.Bind(&u)
		return updateUser(db, id, &u)
	})
	g.GET("/users/:id", func(c echo.Context) error {
		id := ""
		echo.PathParamsBinder(c).String("id", &id)
		isAdmin, _ := c.Get(auth.IsAdminContextKey).(bool)
		u, err := getUser(db, id, isAdmin)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, u)
	})
	g.POST("/category", func(c echo.Context) error {
		data := Category{}
		if err := c.Bind(&data); err != nil {
			return err
		}
		return addCategory(data.Name, data.Activated, db)
	})
	g.PUT("/category/:id", func(c echo.Context) error {
		data := Category{}
		id := ""
		echo.PathParamsBinder(c).String("id", &id)
		if err := c.Bind(&data); err != nil {
			return err
		}
		return updateCategory(id, data.Name, data.Activated, db)
	})
	g.DELETE("/category/:id", func(c echo.Context) error {
		id := ""
		echo.PathParamsBinder(c).String("id", &id)
		return deleteCategory(id, db)
	})
	g.POST("/users/:id/categorize", func(c echo.Context) error {
		id := ""
		echo.PathParamsBinder(c).String("id", &id)
		data := map[string]string{}
		if err := c.Bind(&data); err != nil {
			return err
		}
		err := categorizeUser(db, id, data["category_id"])
		if err != nil {
			pgErr, ok := err.(*pq.Error)
			if ok {
				// if user already cateogrized do nothing
				if pgErr.Code == "23505" {
					return nil
				}
			}
		}
		return err
	})

	g.POST("/users/:id/uncategorize", func(c echo.Context) error {
		id := ""
		echo.PathParamsBinder(c).String("id", &id)
		data := map[string]string{}
		if err := c.Bind(&data); err != nil {
			return err
		}
		return uncategorizeUser(db, id, data["category_id"])
	})

	g.GET("/users/list", func(c echo.Context) error {
		page, _ := strconv.Atoi(c.Request().URL.Query().Get("page"))
		page--
		if page < 0 {
			page = 0
		}
		searchQuery := c.Request().URL.Query().Get("q")
		users, err := listUsers(db, page, searchQuery)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, users)
	})
}
