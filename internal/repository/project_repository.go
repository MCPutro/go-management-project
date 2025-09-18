package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/MCPutro/go-management-project/utils"

	"time"

	"github.com/MCPutro/go-management-project/internal/model"
)

type ProjectRepository interface {
	Create(ctx context.Context, tx *sql.Tx, project *model.Project) error
	GetByID(ctx context.Context, tx *sql.Tx, id int64) (*model.Project, error)
	Update(ctx context.Context, tx *sql.Tx, project *model.Project) error
	Delete(ctx context.Context, tx *sql.Tx, id, deletedBy int64) error
	GetAll(ctx context.Context, tx *sql.Tx) ([]*model.Project, error)
}

type projectRepository struct {
}

func NewProjectRepository() ProjectRepository {
	return &projectRepository{}
}

func (r *projectRepository) Create(ctx context.Context, tx *sql.Tx, project *model.Project) error {
	query := `
		INSERT INTO projects (name, description, created_at, created_by, updated_at, updated_by)
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING id
	`
	now := time.Now()
	return tx.QueryRowContext(ctx, query,
		project.Name, project.Description,
		now, project.CreatedBy,
		now, project.UpdatedBy,
	).Scan(&project.ID)
}

func (r *projectRepository) GetByID(ctx context.Context, tx *sql.Tx, id int64) (*model.Project, error) {
	query := `SELECT id, name, description, created_at, created_by, updated_at, updated_by, deleted_at FROM projects WHERE id = $1 AND deleted_at IS NULL`
	row := tx.QueryRowContext(ctx, query, id)

	var project model.Project
	var deletedAt sql.NullTime

	err := row.Scan(
		&project.ID, &project.Name, &project.Description,
		&project.CreatedAt, &project.CreatedBy, &project.UpdatedAt, &project.UpdatedBy,
		&deletedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, utils.ErrNotFound
	}

	if err != nil {
		return nil, err
	}

	if deletedAt.Valid {
		project.DeletedAt = &deletedAt.Time
	}

	return &project, nil
}

func (r *projectRepository) Update(ctx context.Context, tx *sql.Tx, project *model.Project) error {
	query := `
		UPDATE projects SET name = $1, description = $2, updated_at = $3, updated_by = $4
		WHERE id = $5 AND deleted_at IS NULL
	`
	now := time.Now()
	_, err := tx.ExecContext(ctx, query,
		project.Name, project.Description, now, project.UpdatedBy, project.ID,
	)
	return err
}

func (r *projectRepository) Delete(ctx context.Context, tx *sql.Tx, id, deletedBy int64) error {
	query := `UPDATE projects SET deleted_at = $1, updated_at = $2, updated_by = $3 WHERE id = $4 AND deleted_at IS NULL`
	now := time.Now()
	_, err := tx.ExecContext(ctx, query, now, now, deletedBy, id)
	return err
}

func (r *projectRepository) GetAll(ctx context.Context, tx *sql.Tx) ([]*model.Project, error) {
	query := `SELECT id, name, description, created_at, created_by, updated_at, updated_by, deleted_at FROM projects WHERE deleted_at IS NULL`
	rows, err := tx.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []*model.Project
	for rows.Next() {
		var project model.Project
		var deletedAt sql.NullTime

		err := rows.Scan(
			&project.ID, &project.Name, &project.Description,
			&project.CreatedAt, &project.CreatedBy, &project.UpdatedAt, &project.UpdatedBy,
			&deletedAt,
		)
		if err != nil {
			return nil, err
		}
		if deletedAt.Valid {
			project.DeletedAt = &deletedAt.Time
		}
		projects = append(projects, &project)
	}

	if len(projects) == 0 {
		return nil, utils.ErrNotFound
	}

	return projects, nil
}
