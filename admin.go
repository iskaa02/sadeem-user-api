package main

import (
	"net/http"
	"strconv"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

func registerAdminRoutes(g *echo.Group, db *sqlx.DB) {
	g.PUT("/users/:id", func(c echo.Context) error {
		id := ""
		echo.PathParamsBinder(c).String("id", &id)
		u := User{}
		c.Bind(&u)
		return updateUser(db, id, &u)
	})
	g.POST("/category", func(c echo.Context) error {
		data := Category{}
		err := c.Bind(&data)
		if err != nil {
			return err
		}
		return addCategory(data.Name, data.Activated, db)
	})
	g.PUT("/category/:id", func(c echo.Context) error {
		data := Category{}
		id := ""
		echo.PathParamsBinder(c).String("id", &id)
		err := c.Bind(&data)
		if err != nil {
			return err
		}
		return updateCategory(id, data.Name, data.Activated, db)
	})
	g.DELETE("/category/:id", func(c echo.Context) error {
		// id := chi.URLParam(r, "id")
		id := ""
		echo.PathParamsBinder(c).String("id", &id)
		return deleteCategory(id, db)
	})
	g.POST("/user/:id/categorize", func(c echo.Context) error {
		id := ""
		echo.PathParamsBinder(c).String("id", &id)
		data := map[string]string{}
		if err := c.Bind(&data); err != nil {
			return err
		}
		return categorizeUser(db, id, data["category_id"])
	})

	g.POST("/user/:id/uncategorize", func(c echo.Context) error {
		id := ""
		echo.PathParamsBinder(c).String("id", &id)
		data := map[string]string{}
		c.Bind(&data)
		return uncategorizeUser(db, id, data["category_id"])
	})

	g.GET("/users/list", func(c echo.Context) error {
		page, _ := strconv.Atoi(c.Request().URL.Query().Get("page"))
		page--
		if page < 0 {
			page = 0
		}
		users, err := listUsers(db, page)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, users)
	})
}
