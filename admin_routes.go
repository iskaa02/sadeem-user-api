package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/iskaa02/sadeem-user-api/api_error"
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
	g.POST("/users/:user_id/categorize", func(c echo.Context) error {
		user_id := ""
		echo.PathParamsBinder(c).String("user_id", &user_id)
		data := map[string]string{}
		if err := c.Bind(&data); err != nil {
			return err
		}
		if user_id == "" {
			return api_error.NewBadRequestError("user_id_cannot_be_empty", errors.New("categorize user, empty user_id"))
		}
		if data["category_id"] == "" {
			return api_error.NewBadRequestError("category_id_cannot_be_empty", errors.New("categorize user, empty category_id"))
		}
		err := categorizeUser(db, user_id, data["category_id"])
		if err != nil {
			pgErr, ok := err.(*pq.Error)
			if ok {
				// if user already cateogrized do nothing
				if pgErr.Code == "23505" {
					return c.NoContent(200)
				}
			}
		}
		return err
	})

	g.POST("/users/:user_id/uncategorize", func(c echo.Context) error {
		user_id := ""
		echo.PathParamsBinder(c).String("user_id", &user_id)
		data := map[string]string{}
		if user_id == "" {
			return api_error.NewBadRequestError("user_id_cannot_be_empty", errors.New("uncategorize user, empty user_id"))
		}
		if data["category_id"] == "" {
			return api_error.NewBadRequestError("category_id_cannot_be_empty", errors.New("uncategorize user, empty category_id"))
		}
		if err := c.Bind(&data); err != nil {
			return err
		}
		return uncategorizeUser(db, user_id, data["category_id"])
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
