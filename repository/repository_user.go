package repository

import (
	"context"

	"miraquel.site/api/model"
)

type RepositoryUser interface {
	Migrate(ctx context.Context) error
	Create(ctx context.Context, entity model.User) (*model.User, error)
	All(ctx context.Context) (*[]model.User, error)
	GetById(ctx context.Context, id int64) (*model.User, error)
	GetByParams(ctx context.Context, entity model.User) (*model.User, error)
	Update(ctx context.Context, id int64, entity model.User) (*int64, error)
	Delete(ctx context.Context, id int64) error
	GetByIdWithPosts(ctx context.Context, id int64) (*model.User, error)
}
