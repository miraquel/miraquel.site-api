package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"reflect"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/jackc/pgconn"

	"miraquel.site/api/model"
)

type repositoryCategoryPsql struct {
	db *sql.DB
}

// GetPosts implements RepositoryCategory.
func (r *repositoryCategoryPsql) GetPosts(ctx context.Context, entity model.Category) (*[]model.Category, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT * FROM categories")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var all []model.Category
	for rows.Next() {
		var category model.Category
		if err := rows.Scan(
			&category.Id,
			&category.ParentId,
			&category.Title,
			&category.MetaTitle,
			&category.Slug,
			&category.Content); err != nil {
			return nil, err
		}
		all = append(all, category)
	}

	return &all, nil
}

// All implements RepositoryCategory.
func (r *repositoryCategoryPsql) All(ctx context.Context) (*[]model.Category, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT * FROM categories")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var all []model.Category
	for rows.Next() {
		var category model.Category
		if err := rows.Scan(
			&category.Id,
			&category.ParentId,
			&category.Title,
			&category.MetaTitle,
			&category.Slug,
			&category.Content); err != nil {
			return nil, err
		}
		all = append(all, category)
	}

	return &all, nil
}

// Create implements RepositoryCategory.
func (r *repositoryCategoryPsql) Create(ctx context.Context, entity model.Category) (*model.Category, error) {
	var id int64

	entityBtye, _ := json.Marshal(entity)
	var arguments map[string]interface{}
	_ = json.Unmarshal(entityBtye, &arguments)

	sqlstatement, _, _ := goqu.Insert("categories").Rows(arguments).ToSQL()
	row := r.db.QueryRowContext(ctx, sqlstatement)

	if err := row.Scan(&id); err != nil {
		return nil, err
	}
	entity.Id = id

	return &entity, nil
}

// Delete implements RepositoryCategory.
func (r *repositoryCategoryPsql) Delete(ctx context.Context, id int64) error {
	sqlStatement, _, _ := goqu.Delete("categories").Where(goqu.Ex{"id": id}).ToSQL()

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

// GetById implements RepositoryCategory.
func (r *repositoryCategoryPsql) GetById(ctx context.Context, id int64) (*model.Category, error) {
	sqlStatement, _, _ := goqu.From("categories").Where(goqu.Ex{"id": id}).ToSQL()
	row := r.db.QueryRowContext(ctx, sqlStatement)

	var category model.Category
	if err := row.Scan(
		&category.Id,
		&category.ParentId,
		&category.Title,
		&category.MetaTitle,
		&category.Slug,
		&category.Content); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotExists()
		}
		return nil, err
	}

	return &category, nil
}

// GetByParams implements RepositoryCategory.
func (r *repositoryCategoryPsql) GetByParams(ctx context.Context, entity model.Category) (*model.Category, error) {
	var category model.Category

	v := reflect.ValueOf(entity)
	var arguments goqu.ExOr = make(exp.ExOr)
	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).Interface() != nil || v.Field(i).Interface() != "" {
			arguments[v.Type().Field(i).Name] = v.Field(i).Interface()
		}
	}
	sqlStatement, _, _ := goqu.From("categories").Where(arguments).ToSQL()

	row := r.db.QueryRowContext(ctx, sqlStatement)
	if err := row.Scan(
		&category.Id,
		&category.ParentId,
		&category.Title,
		&category.MetaTitle,
		&category.Slug,
		&category.Content); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotExists()
		}
		return nil, err
	}

	return &category, nil
}

// Migrate implements RepositoryCategory.
func (r *repositoryCategoryPsql) Migrate(ctx context.Context) error {
	query := `
  CREATE TABLE categories (
  id BIGSERIAL PRIMARY KEY,
  parentId BIGINT NULL DEFAULT NULL,
  title VARCHAR(75) NOT NULL,
  metaTitle VARCHAR(100) NULL DEFAULT NULL,
  slug VARCHAR(100) NOT NULL,
  content TEXT NULL DEFAULT NULL);

  CREATE INDEX idx_categories_parent ON categories (parentId ASC);
  ALTER TABLE categories 
  ADD CONSTRAINT fk_categories_parent
    FOREIGN KEY (parentId)
    REFERENCES categories (id)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION;
  `

	_, err := r.db.ExecContext(ctx, query)
	return err
}

// Update implements RepositoryCategory.
func (r *repositoryCategoryPsql) Update(ctx context.Context, id int64, entity model.Category) (*int64, error) {
	sqlStatement, _, _ := goqu.Update("categories").Set(entity).Where(goqu.Ex{"id": id}).ToSQL()
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

func NewRepositoryCategoryPsql(db *sql.DB) RepositoryCategory {
	return &repositoryCategoryPsql{
		db: db,
	}
}
