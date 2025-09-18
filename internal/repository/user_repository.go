package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/MCPutro/go-management-project/utils"
	"time"

	"github.com/MCPutro/go-management-project/internal/model"
)

type UserRepository interface {
	Create(ctx context.Context, tx *sql.Tx, user *model.User) error
	GetByID(ctx context.Context, tx *sql.Tx, id int64) (*model.User, error)
	GetByEmail(ctx context.Context, tx *sql.Tx, email string) (*model.User, error)
	Update(ctx context.Context, tx *sql.Tx, user *model.User) error
	Delete(ctx context.Context, tx *sql.Tx, id, deletedBy int64) error
	GetAll(ctx context.Context, tx *sql.Tx) ([]*model.User, error)
}

type userRepository struct {
}

func NewUserRepository() UserRepository {
	return &userRepository{}
}

func (r *userRepository) Create(ctx context.Context, tx *sql.Tx, user *model.User) error {
	query := `
		INSERT INTO users (name, email, password, created_at, created_by, updated_at, updated_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id
	`
	
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	return tx.QueryRowContext(ctx, query,
		user.Name, user.Email, user.Password,
		now, user.CreatedBy,
		now, user.UpdatedBy,
	).Scan(&user.ID)
}

func (r *userRepository) GetByID(ctx context.Context, tx *sql.Tx, id int64) (*model.User, error) {
	query := `SELECT id, name, email, created_at, created_by, updated_at, updated_by, deleted_at FROM users WHERE id = $1 AND deleted_at IS NULL`
	row := tx.QueryRowContext(ctx, query, id)

	var user model.User
	var deletedAt sql.NullTime

	err := row.Scan(
		&user.ID, &user.Name, &user.Email,
		&user.CreatedAt, &user.CreatedBy, &user.UpdatedAt, &user.UpdatedBy,
		&deletedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, utils.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	if deletedAt.Valid {
		user.DeletedAt = &deletedAt.Time
	}

	return &user, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, tx *sql.Tx, email string) (*model.User, error) {
	query := `SELECT id, name, email, password, created_at, created_by, updated_at, updated_by, deleted_at FROM users WHERE email = $1 AND deleted_at IS NULL`
	row := tx.QueryRowContext(ctx, query, email)

	var user model.User
	var deletedAt sql.NullTime

	err := row.Scan(
		&user.ID, &user.Name, &user.Email, &user.Password,
		&user.CreatedAt, &user.CreatedBy, &user.UpdatedAt, &user.UpdatedBy,
		&deletedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, utils.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	if deletedAt.Valid {
		user.DeletedAt = &deletedAt.Time
	}

	return &user, nil
}

func (r *userRepository) Update(ctx context.Context, tx *sql.Tx, user *model.User) error {
	query := `
		UPDATE users SET name = $1, email = $2, updated_at = $3, updated_by = $4
		WHERE id = $5 AND deleted_at IS NULL
	`
	now := time.Now()
	_, err := tx.ExecContext(ctx, query,
		user.Name, user.Email, now, user.UpdatedBy, user.ID,
	)
	return err
}

func (r *userRepository) Delete(ctx context.Context, tx *sql.Tx, id, deletedBy int64) error {
	query := `UPDATE users SET deleted_at = $1, updated_at = $2, updated_by = $3 WHERE id = $4 AND deleted_at IS NULL`
	now := time.Now()
	_, err := tx.ExecContext(ctx, query, now, now, deletedBy, id)
	return err
}

func (r *userRepository) GetAll(ctx context.Context, tx *sql.Tx) ([]*model.User, error) {
	query := `SELECT id, name, email, created_at, created_by, updated_at, updated_by, deleted_at FROM users WHERE deleted_at IS NULL`
	rows, err := tx.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*model.User
	for rows.Next() {
		var user model.User
		var deletedAt sql.NullTime

		err := rows.Scan(
			&user.ID, &user.Name, &user.Email,
			&user.CreatedAt, &user.CreatedBy, &user.UpdatedAt, &user.UpdatedBy,
			&deletedAt,
		)
		if err != nil {
			return nil, err
		}
		if deletedAt.Valid {
			user.DeletedAt = &deletedAt.Time
		}
		users = append(users, &user)
	}

	if len(users) == 0 {
		return nil, utils.ErrNotFound
	}

	return users, nil
}
