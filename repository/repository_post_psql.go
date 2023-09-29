package repository

import (
	"context"
	"database/sql"
	"errors"
	"reflect"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/goccy/go-json"
	"github.com/jackc/pgconn"

	"miraquel.site/api/model"
)

type repositoryPostPsql struct {
	db *sql.DB
}

// All implements RepositoryPost.
func (r *repositoryPostPsql) All(ctx context.Context) (*[]model.Post, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT * FROM posts")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var all []model.Post
	for rows.Next() {
		var post model.Post
		if err := rows.Scan(
			&post.Id,
			&post.AuthorId,
			&post.ParentId,
			&post.Title,
			&post.MetaTitle,
			&post.Slug,
			&post.Summary,
			&post.Published,
			&post.CreatedAt,
			&post.UpdatedAt,
			&post.PublishedAt,
			&post.Content); err != nil {
			return nil, err
		}
		all = append(all, post)
	}

	return &all, nil
}

// Create implements RepositoryPost.
func (r *repositoryPostPsql) Create(ctx context.Context, entity model.Post) (*model.Post, error) {
	var id int64

	entityBtye, _ := json.Marshal(&entity)
	var arguments map[string]interface{}
	_ = json.Unmarshal(entityBtye, &arguments)

	postSqlStatement, _, _ := goqu.Insert("posts").Rows(arguments).Returning(goqu.T("id")).ToSQL()
	row := r.db.QueryRowContext(ctx, postSqlStatement)

	if err := row.Scan(&id); err != nil {
		return nil, err
	}
	entity.Id = id

	return &entity, nil
}

// Delete implements RepositoryPost.
func (r *repositoryPostPsql) Delete(ctx context.Context, id int64) error {
	sqlStatement, _, _ := goqu.Delete("posts").Where(goqu.Ex{"id": id}).ToSQL()

	res, err := r.db.ExecContext(ctx, sqlStatement)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrDeleteFailed()
	}

	return err
}

// GetById implements RepositoryPost.
func (r *repositoryPostPsql) GetById(ctx context.Context, id int64) (*model.Post, error) {
	sqlStatement, _, _ := goqu.From("posts").Where(goqu.Ex{"id": id}).ToSQL()
	row := r.db.QueryRowContext(ctx, sqlStatement)

	var post model.Post
	if err := row.Scan(
		&post.Id,
		&post.AuthorId,
		&post.ParentId,
		&post.Title,
		&post.MetaTitle,
		&post.Slug,
		&post.Summary,
		&post.Published,
		&post.CreatedAt,
		&post.UpdatedAt,
		&post.PublishedAt,
		&post.Content); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotExists()
		}
		return nil, err
	}

	return &post, nil
}

// GetByParams implements RepositoryPost.
func (r *repositoryPostPsql) GetByParams(ctx context.Context, entity model.Post) (*model.Post, error) {
	var post model.Post

	v := reflect.ValueOf(entity)
	var arguments goqu.ExOr = make(exp.ExOr)
	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).Interface() != nil || v.Field(i).Interface() != "" {
			arguments[v.Type().Field(i).Name] = v.Field(i).Interface()
		}
	}
	sqlStatement, _, _ := goqu.From("posts").Where(arguments).ToSQL()

	row := r.db.QueryRowContext(ctx, sqlStatement)
	if err := row.Scan(
		&post.Id,
		&post.AuthorId,
		&post.ParentId,
		&post.Title,
		&post.MetaTitle,
		&post.Slug,
		&post.Summary,
		&post.Published,
		&post.CreatedAt,
		&post.UpdatedAt,
		&post.PublishedAt,
		&post.Content); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotExists()
		}
		return nil, err
	}

	return &post, nil
}

// Migrate implements RepositoryPost.
func (r *repositoryPostPsql) Migrate(ctx context.Context) error {
	query := `
    CREATE TABLE users (
		id BIGSERIAL PRIMARY KEY,
		firstName VARCHAR(50) NULL DEFAULT NULL,
		middleName VARCHAR(50) NULL DEFAULT NULL,
		lastName VARCHAR(50) NULL DEFAULT NULL,
		mobile VARCHAR(15) NULL,
		email VARCHAR(50) NULL,
		passwordHash VARCHAR(32) NOT NULL,
		registeredAt TIMESTAMP NOT NULL,
		lastLogin TIMESTAMP NULL DEFAULT NULL,
		intro TEXT NULL DEFAULT NULL,
		profile TEXT NULL DEFAULT NULL);

	  CREATE UNIQUE INDEX uq_mobile ON users (mobile ASC);
	  CREATE UNIQUE INDEX uq_email ON users (email ASC);
  `

	_, err := r.db.ExecContext(ctx, query)
	return err
}

// Update implements RepositoryPost.
func (r *repositoryPostPsql) Update(ctx context.Context, id int64, entity model.Post) (*int64, error) {
	sqlStatement, _, _ := goqu.Update("posts").Set(entity).Where(goqu.Ex{"id": id}).ToSQL()
	res, err := r.db.ExecContext(ctx, sqlStatement)
	if err != nil {
		var pgxError *pgconn.PgError
		if errors.As(err, &pgxError) {
			if pgxError.Code == "23505" {
				return nil, ErrDuplicate()
			}
		}
		return nil, err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}

	if rowsAffected == 0 {
		return nil, ErrUpdateFailed()
	}

	return &rowsAffected, nil
}

func NewRepositoryPostPsql(db *sql.DB) RepositoryPost {
	return &repositoryPostPsql{
		db: db,
	}
}
