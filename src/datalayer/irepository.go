package datalayer

import (
	"context"
)

type IRepository[T any] interface {
    Create(ctx context.Context, entity *T) error
    GetByID(ctx context.Context, id int) (*T, error)
    GetAll(ctx context.Context) ([]T, error)
    GetFiltered(ctx context.Context, conditions map[string]any, limit int) ([]T, error)
    Update(ctx context.Context, entity *T) error
    Remove(ctx context.Context, id int) error
}