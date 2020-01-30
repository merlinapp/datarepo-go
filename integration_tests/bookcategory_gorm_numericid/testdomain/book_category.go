package testdomain

import (
	"context"
	"github.com/merlinapp/datarepo-go/integration_tests/model"
	"strconv"
)

type BookCategory struct {
	System         *SystemInstance
	BookCategoryId int
	DBBookCategory *model.BookCategory
}

const bookCategoryCachePrefix = "bc:"

var categoryId = 0

func CreateBookCategory(system *SystemInstance, categoryName string) (*BookCategory, error) {
	bookCategory := &model.BookCategory{
		ID:   categoryId + 1,
		Name: categoryName,
	}
	categoryId = categoryId + 1
	err := system.BookCategoryRepo.Create(system.Ctx, bookCategory)
	if err != nil {
		return nil, err
	}
	return &BookCategory{
		System:         system,
		BookCategoryId: bookCategory.ID,
		DBBookCategory: bookCategory,
	}, nil
}

func (b *BookCategory) ClearData(ctx context.Context) {
	b.System.DB.Delete(&model.BookType{}, "id = ?", b.BookCategoryId)
	b.ClearCacheData(ctx)
}

func (b *BookCategory) ClearCacheData(ctx context.Context) {
	b.System.BookCategoryCacheStore.Delete(ctx, bookCategoryCachePrefix+strconv.Itoa(b.BookCategoryId))
}

func (b *BookCategory) VerifyBookCategoryExists() bool {
	bt := &model.BookCategory{}
	err := b.System.DB.Where("id = ?", b.BookCategoryId).
		Where("name = ?", b.DBBookCategory.Name).
		First(bt).
		Error
	return err == nil
}

func (b *BookCategory) VerifyBookCategoryIsCached(ctx context.Context) bool {
	var bookCategory *model.BookCategory
	found, err := b.System.BookCategoryCacheStore.Get(ctx, bookCategoryCachePrefix+strconv.Itoa(b.BookCategoryId), &bookCategory)
	if err != nil || !found || bookCategory.Name != b.DBBookCategory.Name {
		return false
	}

	return true
}
