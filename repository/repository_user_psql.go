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

type repositoryUserPsql struct {
	db *sql.DB
}

// GetByIdWithPosts implements RepositoryUser.
func (r *repositoryUserPsql) GetByIdWithPosts(ctx context.Context, id int64) (*model.User, error) {
	//userSqlStatement, _, _ := goqu.From("users").Where(goqu.Ex{"id": id}).ToSQL()
	row := r.db.QueryRowContext(ctx, "SELECT * FROM users WHERE id = 17")

	var user model.User
	if err := row.Scan(
		&user.Id,
		&user.FirstName,
		&user.MiddleName,
		&user.LastName,
		&user.Mobile,
		&user.Email,
		&user.PasswordHash,
		&user.RegisteredAt,
		&user.LastLogin,
		&user.Intro,
		&user.Profile); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotExists()
		}
		return nil, err
	}

	postsSqlStatement, _, _ := goqu.From("posts").Where(goqu.Ex{"authorid": id}).ToSQL()
	postRows, postsErr := r.db.QueryContext(ctx, postsSqlStatement)

	if postsErr != nil {
		return nil, postsErr
	}
	defer postRows.Close()

	var posts []model.Post
	for postRows.Next() {
		var post model.Post
		if err := postRows.Scan(
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
		posts = append(posts, post)
	}

	user.Posts = &posts

	return &user, nil
}

func NewRepositoryUserPsql(db *sql.DB) RepositoryUser {
	return &repositoryUserPsql{
		db: db,
	}
}

// All implements Repository.
func (r *repositoryUserPsql) All(ctx context.Context) (*[]model.User, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT * FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var all []model.User
	for rows.Next() {
		var user model.User
		if err := rows.Scan(
			&user.Id,
			&user.FirstName,
			&user.MiddleName,
			&user.LastName,
			&user.Mobile,
			&user.Email,
			&user.PasswordHash,
			&user.RegisteredAt,
			&user.LastLogin,
			&user.Intro,
			&user.Profile); err != nil {
			return nil, err
		}
		all = append(all, user)
	}

	return &all, nil
}

// Create implements Repository.
func (r *repositoryUserPsql) Create(ctx context.Context, entity model.User) (*model.User, error) {
	var id int64

	entityByte, _ := json.Marshal(&entity)
	var arguments map[string]interface{}
	_ = json.Unmarshal(entityByte, &arguments)

	sqlStatement, _, _ := goqu.Insert("users").Rows(arguments).Returning(goqu.T("id")).ToSQL()
	row := r.db.QueryRowContext(ctx, sqlStatement)

	if err := row.Scan(&id); err != nil {
		return nil, err
	}
	entity.Id = id

	return &entity, nil
}

// Delete implements Repository.
func (r *repositoryUserPsql) Delete(ctx context.Context, id int64) error {
	sqlStatement, _, _ := goqu.Delete("users").Where(goqu.Ex{"id": id}).ToSQL()

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

// GetById implements Repository.
func (r *repositoryUserPsql) GetById(ctx context.Context, id int64) (*model.User, error) {
	sqlStatement, _, _ := goqu.From("users").Where(goqu.Ex{"id": id}).ToSQL()
	row := r.db.QueryRowContext(ctx, sqlStatement)

	var user model.User
	if err := row.Scan(
		&user.Id,
		&user.FirstName,
		&user.MiddleName,
		&user.LastName,
		&user.Mobile,
		&user.Email,
		&user.PasswordHash,
		&user.RegisteredAt,
		&user.LastLogin,
		&user.Intro,
		&user.Profile); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotExists()
		}
		return nil, err
	}

	return &user, nil
}

// GetByParams implements Repository.
func (r *repositoryUserPsql) GetByParams(ctx context.Context, entity model.User) (*model.User, error) {
	var user model.User

	v := reflect.ValueOf(entity)
	var arguments goqu.ExOr = make(exp.ExOr)
	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).Interface() != nil || v.Field(i).Interface() != "" {
			arguments[v.Type().Field(i).Name] = v.Field(i).Interface()
		}
	}
	sqlStatement, _, _ := goqu.From("users").Where(arguments).ToSQL()

	row := r.db.QueryRowContext(ctx, sqlStatement)
	if err := row.Scan(
		&user.Id,
		&user.FirstName,
		&user.MiddleName,
		&user.LastName,
		&user.Mobile,
		&user.Email,
		&user.PasswordHash,
		&user.RegisteredAt,
		&user.LastLogin,
		&user.Intro,
		&user.Profile); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotExists()
		}
		return nil, err
	}

	return &user, nil
}

// Migrate implements Repository.
func (r *repositoryUserPsql) Migrate(ctx context.Context) error {
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

// Update implements Repository.
func (r *repositoryUserPsql) Update(ctx context.Context, id int64, entity model.User) (*int64, error) {
	sqlStatement, _, _ := goqu.Update("users").Set(entity).Where(goqu.Ex{"id": id}).ToSQL()
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
