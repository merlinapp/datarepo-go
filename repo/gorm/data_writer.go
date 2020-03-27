package gorm

import (
	"context"
	"errors"
	"github.com/jinzhu/gorm"
	"github.com/merlinapp/datarepo-go/drreflect"
)

type dataWriter struct {
	db          *gorm.DB
	typeHandler drreflect.TypeHandler
}

func (w *dataWriter) Create(ctx context.Context, value interface{}) error {
	err := w.ensurePointer(value)
	if err != nil {
		return err
	}

	err = w.db.Create(value).Error
	if err != nil {
		return err
	}

	return nil
}

func (w *dataWriter) Update(ctx context.Context, value interface{}) error {
	err := w.ensurePointer(value)
	if err != nil {
		return err
	}

	err = w.db.Save(value).Error
	if err != nil {
		return err
	}

	return nil
}

func (w *dataWriter) PartialUpdate(ctx context.Context, value interface{}) error {
	err := w.ensurePointer(value)
	if err != nil {
		return err
	}
	t := w.typeHandler.NewPtrToElement().Element()

	err = w.db.Model(t).Updates(value).Find(value).Error

	if err != nil {
		return err
	}

	return nil
}

func (w *dataWriter) ensurePointer(value interface{}) error {
	if !w.typeHandler.IsOfPtrType(value) {
		return errors.New("The provided value isn't of the expected type: " + w.typeHandler.Type().String())
	}
	return nil
}
