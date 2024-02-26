package main

import (
	"fmt"

	"github.com/guregu/null"
	"github.com/huandu/go-sqlbuilder"
	"github.com/iskaa02/sadeem-user-api/api_error"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type Category struct {
	ID        string    `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Activated null.Bool `json:"activated" db:"activated"`
}

func listCategories(db *sqlx.DB, isAdmin bool, page int, searchQuery string) ([]Category, error) {
	result := []Category{}
	sb := sqlbuilder.Select("*").
		From("category").
		Limit(10).
		Offset(page * 10)
	if !isAdmin {
		sb.Where(sb.NotEqual("activated", false))
	}
	if searchQuery != "" {
		sb.Where(sb.Like("name", "%"+searchQuery+"%"))
	}
	sql, args := sb.Build()
	err := db.Select(&result, sql, args...)
	return result, err
}

func addCategory(name string, activated null.Bool, db *sqlx.DB) error {
	_, err := db.Exec("INSERT INTO category(name,activated) VALUES($1,$2)", name, activated.Bool)
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok {
			if pgErr.Code == "23505" {
				return api_error.NewBadRequestError("category_already_exists", err)
			}
		}
	}
	return err
}

func deleteCategory(id string, db *sqlx.DB) error {
	_, err := db.Exec("DELETE FROM category WHERE id=$1", id)
	return err
}

func updateCategory(id string, name string, activated null.Bool, db *sqlx.DB) error {
	sb := sqlbuilder.Update("category")
	sb.Where(sb.EQ("id", id))
	fmt.Println(activated.Valid)
	if activated.Valid {
		sb.SetMore(sb.Assign("activated", activated.Bool))
	}
	if name != "" {
		sb.SetMore(sb.Assign("name", name))
	}

	sql, args := sb.Build()
	_, err := db.Exec(sql, args...)
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok {
			if pgErr.Code == "23505" {
				return api_error.NewBadRequestError("category_already_exists", err)
			}
		}
	}
	return err
}

func categorizeUser(db *sqlx.DB, userID, categoryID string) error {
	_, err := db.Exec("INSERT INTO user_category(user_id,category_id) VALUES($1,$2)", userID, categoryID)
	return err
}

func uncategorizeUser(db *sqlx.DB, userID, categoryID string) error {
	_, err := db.Exec("DELETE FROM user_category WHERE user_id=$1 AND category_id=$2", userID, categoryID)
	return err
}

func getUserCategory(db *sqlx.DB, userID string, isAdmin bool) ([]Category, error) {
	result := []Category{}
	sb := sqlbuilder.Select("c.*").From("category c")
	sb.Join("user_category uc", "uc.category_id=c.id")
	sb.Where(sb.Equal("uc.user_id", userID))
	if !isAdmin {
		sb.Where(sb.NotEqual("activated", false))
	}
	sql, args := sb.Build()
	err := db.Select(&result, sql, args...)
	return result, err
}
