package model

import (
	dbm "../db"
	"database/sql"
	"github.com/pkg/errors"
)

type User struct {
	ID   int
	Name string
}

func (u User) Get(db *dbm.DB) (User, error) {
	var user User
	err := db.Find(&user, u.ID)
	if errors.Is(err, sql.ErrNoRows) {
		return user, dbm.ErrRecordNotFound
	}
	if err != nil {
		return user, err
	}
	return user, nil
}
