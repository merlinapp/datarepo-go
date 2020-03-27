package testdomain

import (
	"context"
	"github.com/merlinapp/datarepo-go/integration_tests/model"
	"github.com/satori/uuid"
)

type Book struct {
	System *SystemInstance
	BookId string
	DBBook *model.Book
}

const bookCachePrefix = "b:"

func CreateBook(system *SystemInstance, authorId, status string) (*Book, error) {
	book := &model.Book{
		ID:       uuid.NewV4().String(),
		AuthorID: authorId,
		Status:   status,
	}
	err := system.BookRepo.Create(system.Ctx, book)
	if err != nil {
		return nil, err
	}
	return &Book{
		System: system,
		BookId: book.ID,
		DBBook: book,
	}, nil
}

func (b *Book) ClearData(ctx context.Context) {
	b.System.DB.Delete(&model.Book{}, "id = ?", b.BookId)
	b.ClearCacheData(ctx)
}

func (b *Book) ClearCacheData(ctx context.Context) {
	b.System.BookCacheStore.Delete(ctx, bookCachePrefix+b.BookId)
}

func (b *Book) VerifyBookExists(ctx context.Context) bool {
	sm := &model.Book{}
	err := b.System.DB.Where("id = ?", b.BookId).
		Where("status = ?", b.DBBook.Status).
		First(sm).
		Error
	return err == nil
}

func (b *Book) VerifyBookIsCached(ctx context.Context) bool {
	var book *model.Book
	found, err := b.System.BookCacheStore.Get(ctx, bookCachePrefix+b.BookId, &book)
	if err != nil || !found || book.Status != b.DBBook.Status {
		return false
	}

	return true
}

func (b *Book) UpdateStatus(ctx context.Context, newStatus string) error {
	updatedBook := &model.Book{
		ID:       b.BookId,
		AuthorID: b.DBBook.AuthorID,
		Status:   newStatus,
	}
	err := b.System.BookRepo.Update(ctx, updatedBook)
	if err != nil {
		return err
	}
	b.DBBook = updatedBook
	return nil
}

func (b *Book) PartialUpdateStatus(ctx context.Context, newStatus string) error {
	updatedBook := &model.Book{
		ID:     b.BookId,
		Status: newStatus,
	}
	err := b.System.BookRepo.PartialUpdate(ctx, updatedBook)
	if err != nil {
		return err
	}

	b.DBBook = updatedBook

	return nil
}
