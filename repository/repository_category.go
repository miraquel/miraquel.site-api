package repository

import (
	"context"

	"miraquel.site/api/model"
)

type RepositoryCategory interface {
	Migrate(ctx context.Context) error
	Create(ctx context.Context, entity model.Category) (*model.Category, error)
	All(ctx context.Context) (*[]model.Category, error)
	GetById(ctx context.Context, id int64) (*model.Category, error)
	GetByParams(ctx context.Context, entity model.Category) (*model.Category, error)
	Update(ctx context.Context, id int64, entity model.Category) (*int64, error)
	Delete(ctx context.Context, id int64) error
	GetPosts(ctx context.Context, entity model.Category) (*[]model.Category, error)
}
