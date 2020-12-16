package service

import (
	"context"
	"fmt"
	pb "github.com/shaojunda/Go-000/Week4/explorer/api/book/v1"
	"github.com/shaojunda/Go-000/Week4/explorer/internal/biz"
	"strconv"
)

var _ pb.BookServiceServer = (*BookService)(nil)

type BookService struct {
	buc *biz.BookUsecase
}

func NewBookService(buc *biz.BookUsecase) *BookService {
	return &BookService{buc: buc}
}

func (svr *BookService) CreateBook(ctx context.Context, r *pb.CreateBookRequest) (*pb.CreateBookResponse, error) {
	book := new(biz.Book)
	book.Name = r.Name
	book.Author = r.Author
	book.Publisher = r.Publisher
	book.URL = r.Url
	f, err := strconv.ParseFloat(r.Price, 64)
	if err != nil {
		return nil, err
	}
	book.Price = f
	book.InAppPrice = f * 0.75
	b, err := svr.buc.CreateBook(book)
	if err != nil {
		return nil, err
	}
	return &pb.CreateBookResponse{
		Book: &pb.Book{
			Name:      b.Name,
			Author:    b.Author,
			Publisher: b.Publisher,
			Price:     fmt.Sprintf("%f", b.Price),
			Url:       b.URL,
		},
	}, err
}

func (svr *BookService) SearchBook(ctx context.Context, r *pb.SearchBookRequest) (*pb.SearchBookResponse, error) {
	books, err := svr.buc.SearchBook(r.QueryString)
	if err != nil {
		return nil, err
	}
	rbs := make([]*pb.Book, 0, len(books))
	for _, book := range books {
		rbs = append(rbs, &pb.Book{
			Name:      book.Name,
			Author:    book.Author,
			Publisher: book.Publisher,
			Price:     fmt.Sprintf("%f", book.Price),
			Url:       book.URL,
		})
	}

	return &pb.SearchBookResponse{
		Books: rbs,
	}, nil
}
