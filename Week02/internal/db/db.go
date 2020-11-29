package db

import (
	"database/sql"
	"github.com/pkg/errors"
)

var ErrRecordNotFound = errors.New("record not found")

type DB struct {
}

func (db DB) Find(interface{}, int) error {
	return sql.ErrNoRows
}

func New() *DB {
	return &DB{}
}
