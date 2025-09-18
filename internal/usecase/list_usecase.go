// internal/usecase/list_usecase.go
package usecase

import (
	"context"
	"database/sql"

	"github.com/MCPutro/go-management-project/internal/model"
	"github.com/MCPutro/go-management-project/internal/repository"
)

type ListUsecase interface {
	CreateList(ctx context.Context, list *model.List) error
	GetListsByProjectID(ctx context.Context, projectID int64) ([]*model.List, error)
	GetListByID(ctx context.Context, id int64) (*model.List, error)
	UpdateList(ctx context.Context, list *model.List) error
	DeleteList(ctx context.Context, id int64, deletedBy int64) error
}

type listUsecase struct {
	db       *sql.DB
	listRepo repository.ListRepository
}

func NewListUsecase(db *sql.DB, listRepository repository.ListRepository) ListUsecase {
	return &listUsecase{
		db:       db,
		listRepo: listRepository,
	}
}

func (l *listUsecase) CreateList(ctx context.Context, list *model.List) error {
	tx, err := l.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	err = l.listRepo.Create(ctx, tx, list)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (l *listUsecase) GetListsByProjectID(ctx context.Context, projectID int64) ([]*model.List, error) {
	tx, err := l.db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	lists, err := l.listRepo.GetByProjectID(ctx, tx, projectID)
	if err != nil {
		return nil, err
	}

	return lists, nil
}

func (l *listUsecase) GetListByID(ctx context.Context, id int64) (*model.List, error) {
	tx, err := l.db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	query := `SELECT id, project_id, name, position, created_at, created_by, updated_at, updated_by, deleted_at FROM lists WHERE id = $1 AND deleted_at IS NULL`
	row := tx.QueryRowContext(ctx, query, id)

	var list model.List
	err = row.Scan(
		&list.ID, &list.ProjectID, &list.Name, &list.Position,
		&list.CreatedAt, &list.CreatedBy, &list.UpdatedAt, &list.UpdatedBy, &list.DeletedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil // atau return error jika diinginkan
	}
	if err != nil {
		return nil, err
	}

	return &list, nil
}

func (l *listUsecase) UpdateList(ctx context.Context, list *model.List) error {
	tx, err := l.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	err = l.listRepo.Update(ctx, tx, list)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (l *listUsecase) DeleteList(ctx context.Context, id int64, deletedBy int64) error {
	tx, err := l.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	err = l.listRepo.Delete(ctx, tx, id, deletedBy)
	if err != nil {
		return err
	}

	return tx.Commit()
}
