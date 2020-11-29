package service

import (
	"../dao"
	"../db"
)

type Service struct {
	dao *dao.Dao
}

func New() Service {
	return Service{dao: dao.New(db.New())}
}
