package model

import (
	dbm "../db"
	"database/sql"
	"fmt"
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
		return user, errors.Wrapf(dbm.ErrRecordNotFound, fmt.Sprintf("sql: select * from users where id = %d error: %v", u.ID, err))
	}
	if err != nil {
		return user, err
	}
	return user, nil
}
