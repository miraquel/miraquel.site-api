package repository

import (
	"context"

	"miraquel.site/api/model"
)

type RepositoryPost interface {
	Migrate(ctx context.Context) error
	Create(ctx context.Context, entity model.Post) (*model.Post, error)
	All(ctx context.Context) (*[]model.Post, error)
	GetById(ctx context.Context, id int64) (*model.Post, error)
	GetByParams(ctx context.Context, entity model.Post) (*model.Post, error)
	Update(ctx context.Context, id int64, entity model.Post) (*int64, error)
	Delete(ctx context.Context, id int64) error
}
