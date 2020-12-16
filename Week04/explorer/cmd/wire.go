//+build wireinject

package main

import (
	"github.com/shaojunda/Go-000/Week4/explorer/internal/service"
	"github.com/shaojunda/Go-000/Week4/explorer/internal/biz"
	"github.com/shaojunda/Go-000/Week4/explorer/internal/data"
	"github.com/google/wire"
)

func InitializeBookService() *service.BookService {
	wire.Build(service.NewBookService, biz.NewBookUsecase, data.NewBookRepo, data.NewDB)
	return &service.BookService{}
}
