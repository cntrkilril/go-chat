package gateway

import (
	"context"
	"database/sql"
	"github.com/cntrkilril/go-chat-common/pkg/errors"
	"github.com/jmoiron/sqlx"
	"user-service/internal/entity"
)

type UserRepository struct {
	db *sqlx.DB
}

func (r UserRepository) Save(ctx context.Context, params entity.CreateUserParams) (result entity.User, err error) {
	q := `
			INSERT into users
			(username, password)
			VALUES ($1,$2)
			RETURNING id, username, password
	`
	err = r.db.GetContext(ctx, &result, q, params.Username, params.Password)
	if err != nil {
		return entity.User{}, err
	}

	return result, nil
}

func (r UserRepository) UpdatePassword(ctx context.Context, params entity.UpdatePasswordParams) (result entity.User, err error) {
	q := `
			UPDATE users
			SET password=$1
			WHERE id=$2
			RETURNING id, username, password
	`
	err = r.db.GetContext(ctx, &result, q, params.NewPassword, params.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return entity.User{}, errors.ErrUserNotFound
		}
		return entity.User{}, err
	}

	return result, nil
}

func (r UserRepository) GetByID(ctx context.Context, id int64) (result entity.User, err error) {
	q := `
			SELECT id, username, password
			FROM users
			WHERE id=$1
	`
	err = r.db.GetContext(ctx, &result, q, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return entity.User{}, errors.ErrUserNotFound
		}
		return entity.User{}, err
	}

	return result, nil
}

func (r UserRepository) GetByUsername(ctx context.Context, username string) (result entity.User, err error) {
	q := `
			SELECT id, username, password
			FROM users
			WHERE username=$1
	`
	err = r.db.GetContext(ctx, &result, q, username)
	if err != nil {
		if err == sql.ErrNoRows {
			return entity.User{}, errors.ErrUserNotFound
		}
		return entity.User{}, err
	}

	return result, nil
}

func (r UserRepository) GetByIncludedUsername(ctx context.Context, params entity.GetUsersByUsernameParams) (result []entity.User, err error) {
	q := `
			SELECT id, username, password
			FROM users
			WHERE username LIKE $1
			LIMIT $2 OFFSET $3
	`
	err = r.db.SelectContext(ctx, &result, q, "%"+params.Username+"%", params.Limit, params.Offset)
	if err != nil {
		return []entity.User{}, err
	}

	return result, nil
}

func (r UserRepository) CountByIncludedUsername(ctx context.Context, username string) (result int64, err error) {
	q := `
			SELECT COUNT(*) as count
			FROM users
			WHERE username LIKE $1
	`
	err = r.db.GetContext(ctx, &result, q, "%"+username+"%")
	if err != nil {
		return 0, err
	}

	return result, nil
}

func (r UserRepository) Delete(ctx context.Context, id int64) error {
	q := `
			DELETE 
			FROM users
			WHERE id=$1
			RETURNING id
	`
	_, err := r.db.ExecContext(ctx, q, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.ErrUserNotFound
		}
		return err
	}

	return nil
}

var _ UserGateway = (*UserRepository)(nil)

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}
