package main

import (
	"github.com/jmoiron/sqlx"
)

type Category struct {
	ID        string `json:"id" db:"id"`
	Name      string `json:"name" db:"name"`
	Activated bool   `json:"activated" db:"activated"`
}

func listCategories(db *sqlx.DB, isAdmin bool, page int) []Category {
	result := []Category{}
	if isAdmin {
		err := db.Select(&result, "SELECT * FROM category SORT BY name LIMIT 10 OFFSET $1 SORT BY name", page*10)
		if err != nil {
			return nil
		}
		return result
	}
	err := db.Select(&result, "SELECT * FROM category WHERE activated=true SORT BY name LIMIT 10 OFFSET $1 SORT BY name", page*10)
	if err != nil {
		return nil
	}
	return result
}

func addCategory(name string, activated bool, db *sqlx.DB) error {
	_, err := db.Exec("INSERT INTO category(name,activated) VALUES($1,$2)", name, activated)
	if err != nil {
		return err
	}
	return nil
}

func deleteCategory(id string, db *sqlx.DB) error {
	_, err := db.Exec("DELETE FROM category WHERE id=$1", id)
	if err != nil {
		return err
	}
	return nil
}

func updateCategory(id string, name string, activated bool, db *sqlx.DB) error {
	_, err := db.Exec("UPDATE category SET name=$2,activated=$3 WHERE id=$1", id, name, activated)
	if err != nil {
		return err
	}
	return nil
}
