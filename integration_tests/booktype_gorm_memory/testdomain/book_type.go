package testdomain

import (
	"context"
	"github.com/merlinapp/datarepo-go/integration_tests/model"
	"github.com/satori/uuid"
)

type BookType struct {
	System     *SystemInstance
	BookTypeId string
	DBBookType *model.BookType
}

const bookTypeCachePrefix = "bt:"

func CreateBookType(system *SystemInstance, typeName string) (*BookType, error) {
	bookType := &model.BookType{
		ID:   uuid.NewV4().String(),
		Name: typeName,
	}
	err := system.BookTypeRepo.Create(system.Ctx, bookType)
	if err != nil {
		return nil, err
	}
	return &BookType{
		System:     system,
		BookTypeId: bookType.ID,
		DBBookType: bookType,
	}, nil
}

func (b *BookType) ClearData(ctx context.Context) {
	b.System.DB.Delete(&model.BookType{}, "id = ?", b.BookTypeId)
	b.ClearCacheData(ctx)
}

func (b *BookType) ClearCacheData(ctx context.Context) {
	b.System.BookTypeCacheStore.Delete(ctx, bookTypeCachePrefix+b.BookTypeId)
}

func (b *BookType) VerifyBookTypeExists(ctx context.Context) bool {
	bt := &model.BookType{}
	err := b.System.DB.Where("id = ?", b.BookTypeId).
		Where("name = ?", b.DBBookType.Name).
		First(bt).
		Error
	return err == nil
}

func (b *BookType) VerifyBookTypeIsCached(ctx context.Context) bool {
	var bookType *model.BookType
	found, err := b.System.BookTypeCacheStore.Get(ctx, bookTypeCachePrefix+b.BookTypeId, &bookType)
	if err != nil || !found || bookType.Name != b.DBBookType.Name {
		return false
	}

	return true
}
