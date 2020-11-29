package dao

import "../db"

type Dao struct {
	engine *db.DB
}

func New(engine *db.DB) *Dao {
	return &Dao{engine: engine}
}
