package biz

type Book struct {
	Name       string
	Author     string
	Publisher  string
	Price      float64
	InAppPrice float64
	URL        string
}

type BookRepo interface {
	SearchBook(queryString string) ([]Book, error)
	CreateBook(book *Book) (Book, error)
}

func NewBookUsecase(repo BookRepo) *BookUsecase {
	return &BookUsecase{repo: repo}
}

type BookUsecase struct {
	repo BookRepo
}

func (uc *BookUsecase) SearchBook(queryString string) ([]Book, error) {
	return uc.repo.SearchBook(queryString)
}
func (uc *BookUsecase) CreateBook(b *Book) (Book, error) {
	return uc.CreateBook(b)
}
