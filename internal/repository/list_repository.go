package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/MCPutro/go-management-project/utils"
	"time"

	"github.com/MCPutro/go-management-project/internal/model"
)

type ListRepository interface {
	Create(ctx context.Context, tx *sql.Tx, list *model.List) error
	GetByID(ctx context.Context, tx *sql.Tx, id int64) (*model.List, error)
	GetByProjectID(ctx context.Context, tx *sql.Tx, projectID int64) ([]*model.List, error)
	Update(ctx context.Context, tx *sql.Tx, list *model.List) error
	Delete(ctx context.Context, tx *sql.Tx, id, deletedBy int64) error
}

type listRepository struct {
}

func NewListRepository() ListRepository {
	return &listRepository{}
}

func (r *listRepository) Create(ctx context.Context, tx *sql.Tx, list *model.List) error {
	query := `
		INSERT INTO lists (project_id, name, position, created_at, created_by, updated_at, updated_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id
	`
	now := time.Now()
	return tx.QueryRowContext(ctx, query,
		list.ProjectID, list.Name, list.Position,
		now, list.CreatedBy,
		now, list.UpdatedBy,
	).Scan(&list.ID)
}

func (r *listRepository) GetByID(ctx context.Context, tx *sql.Tx, id int64) (*model.List, error) {
	query := `SELECT id, project_id, name, position, created_at, created_by, updated_at, updated_by, deleted_at FROM lists WHERE id = $1 AND deleted_at IS NULL`
	row := tx.QueryRowContext(ctx, query, id)

	var list model.List
	var deletedAt sql.NullTime

	err := row.Scan(
		&list.ID, &list.ProjectID, &list.Name, &list.Position,
		&list.CreatedAt, &list.CreatedBy, &list.UpdatedAt, &list.UpdatedBy,
		&deletedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, utils.ErrNotFound
	}

	if err != nil {
		return nil, err
	}

	if deletedAt.Valid {
		list.DeletedAt = &deletedAt.Time
	}

	return &list, nil
}

func (r *listRepository) GetByProjectID(ctx context.Context, tx *sql.Tx, projectID int64) ([]*model.List, error) {
	query := `SELECT id, project_id, name, position, created_at, created_by, updated_at, updated_by, deleted_at FROM lists WHERE project_id = $1 AND deleted_at IS NULL ORDER BY position ASC`
	rows, err := tx.QueryContext(ctx, query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lists []*model.List
	for rows.Next() {
		var list model.List
		var deletedAt sql.NullTime

		err := rows.Scan(
			&list.ID, &list.ProjectID, &list.Name, &list.Position,
			&list.CreatedAt, &list.CreatedBy, &list.UpdatedAt, &list.UpdatedBy,
			&deletedAt,
		)
		if err != nil {
			return nil, err
		}
		if deletedAt.Valid {
			list.DeletedAt = &deletedAt.Time
		}
		lists = append(lists, &list)
	}

	if len(lists) == 0 {
		return nil, utils.ErrNotFound
	}

	return lists, nil
}

func (r *listRepository) Update(ctx context.Context, tx *sql.Tx, list *model.List) error {
	query := `
		UPDATE lists SET name = $1, position = $2, updated_at = $3, updated_by = $4
		WHERE id = $5 AND deleted_at IS NULL
	`
	now := time.Now()
	_, err := tx.ExecContext(ctx, query,
		list.Name, list.Position, now, list.UpdatedBy, list.ID,
	)
	return err
}

func (r *listRepository) Delete(ctx context.Context, tx *sql.Tx, id, deletedBy int64) error {
	query := `UPDATE lists SET deleted_at = $1, updated_at = $2, updated_by = $3 WHERE id = $4 AND deleted_at IS NULL`
	now := time.Now()
	_, err := tx.ExecContext(ctx, query, now, now, deletedBy, id)
	return err
}
