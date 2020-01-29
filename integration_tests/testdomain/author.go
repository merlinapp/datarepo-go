package testdomain

import (
	"context"
	"github.com/merlinapp/datarepo-go"
	"github.com/merlinapp/datarepo-go/integration_tests/bootstrap"
	"github.com/merlinapp/datarepo-go/integration_tests/model"
	"github.com/satori/uuid"
)

type Author struct {
	System   *bootstrap.SystemInstance
	AuthorId string
	DBAuthor *model.Author
}

const authorCachePrefix = "a:"

func CreateAuthor(system *bootstrap.SystemInstance) *Author {
	author := &model.Author{
		ID:   uuid.NewV4().String(),
		Name: "Name",
	}
	system.DB.Create(&author)
	return &Author{
		System:   system,
		AuthorId: author.ID,
		DBAuthor: author,
	}
}

func (a *Author) ClearData(ctx context.Context) {
	a.System.DB.Delete(&model.Book{}, "author_id = ?", a.AuthorId)
	a.ClearCacheData(ctx)
}

func (a *Author) ClearCacheData(ctx context.Context) {
	a.System.BookCacheStore.Delete(ctx, authorCachePrefix+a.AuthorId)
}

func (a *Author) CreateBook(ctx context.Context, status string) (*Book, error) {
	return CreateBook(a.System, a.AuthorId, status)
}

func (a *Author) VerifyBookIsCached(ctx context.Context, bookId, status string) bool {
	var books []*model.Book
	found, err := a.System.BookCacheStore.Get(ctx, authorCachePrefix+a.AuthorId, &books)
	if err != nil || !found {
		return false
	}
	for _, b := range books {
		if b.ID == bookId && b.Status == status {
			return true
		}
	}
	return false
}

func (a *Author) GetBooks(ctx context.Context) ([]*model.Book, error) {
	result, err := a.System.BookRepo.FindByKey(ctx, "AuthorID", a.AuthorId)
	if err != nil {
		return nil, err
	}
	var books []*model.Book
	result.InjectResult(&books)
	return books, nil
}

func GetBooksForAuthors(system *bootstrap.SystemInstance, ids []string) ([][]*model.Book, error) {
	results, err := system.BookRepo.FindByKeys(system.Ctx, "AuthorID", ids)
	if err != nil {
		return nil, err
	}
	var booksPerAuthor [][]*model.Book
	datarepo.InjectResults(results, &booksPerAuthor)
	return booksPerAuthor, nil
}
