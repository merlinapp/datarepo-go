package datarepo

import (
	"context"
)

type DataWriter interface {
	// Inserts the provided value into the repository
	Create(ctx context.Context, value interface{}) error
	// Updates the provided value in the repository
	Update(ctx context.Context, value interface{}) error
}
