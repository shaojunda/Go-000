package data

import "github.com/shaojunda/Go-000/Week4/explorer/internal/biz"

var _ biz.BookRepo = (*bookRepo)(nil)

type bookRepo struct {
	db *DB
}

func NewBookRepo(db *DB) biz.BookRepo{
	return &bookRepo{db}
}

func (b *bookRepo) SearchBook(queryString string) ([]biz.Book, error) {
	return b.db.SearchBook(queryString)
}

func (b *bookRepo) CreateBook(book *biz.Book) (biz.Book, error) {
	return b.db.CreateBook(book)
}

type DB struct{}

func NewDB() *DB {
	return &DB{}
}

func (db *DB) SearchBook(queryString string) ([]biz.Book, error) {
	return []biz.Book{
		{
			Name:       "The Soul of Care",
			Author:     "凯博文",
			Publisher:  "中信出版社",
			Price:      58,
			InAppPrice: 36,
			URL:        "https://book.douban.com/subject/35217694/",
		},
	}, nil
}

func (db *DB) CreateBook(book *biz.Book) (biz.Book, error) {
	return biz.Book{
		Name:       "你当像鸟飞往你的山",
		Author:     "塔拉·韦斯特弗",
		Publisher:  "南海出版公司",
		Price:      59,
		InAppPrice: 37,
		URL:        "https://book.douban.com/subject/33440205/",
	}, nil
}
