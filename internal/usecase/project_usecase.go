package usecase

import (
	"context"
	"database/sql"

	"github.com/MCPutro/go-management-project/internal/model"
	"github.com/MCPutro/go-management-project/internal/repository"
)

type ProjectUsecase interface {
	CreateProject(ctx context.Context, project *model.Project) error
	GetProjectByID(ctx context.Context, id int64) (*model.Project, error)
	UpdateProject(ctx context.Context, project *model.Project) error
	DeleteProject(ctx context.Context, id int64, deletedBy int64) error
}

type projectUsecase struct {
	db          *sql.DB
	projectRepo repository.ProjectRepository
	listRepo    repository.ListRepository
}

func NewProjectUsecase(db *sql.DB, projectRepository repository.ProjectRepository, listRepository repository.ListRepository) ProjectUsecase {
	return &projectUsecase{
		db:          db,
		projectRepo: projectRepository,
		listRepo:    listRepository,
	}
}

func (p *projectUsecase) CreateProject(ctx context.Context, project *model.Project) error {
	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	err = p.projectRepo.Create(ctx, tx, project)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (p *projectUsecase) GetProjectByID(ctx context.Context, id int64) (*model.Project, error) {
	tx, err := p.db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	project, err := p.projectRepo.GetByID(ctx, tx, id)
	if err != nil {
		return nil, err
	}

	return project, nil
}

func (p *projectUsecase) UpdateProject(ctx context.Context, project *model.Project) error {
	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	err = p.projectRepo.Update(ctx, tx, project)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (p *projectUsecase) DeleteProject(ctx context.Context, id int64, deletedBy int64) error {
	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	err = p.projectRepo.Delete(ctx, tx, id, deletedBy)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// Method baru:
func (p *projectUsecase) CreateProjectWithDefaultList(ctx context.Context, project *model.Project, defaultListName string) error {
	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// 1. Simpan project
	err = p.projectRepo.Create(ctx, tx, project)
	if err != nil {
		return err
	}

	// 2. Buat list default
	defaultList := &model.List{
		ProjectID: project.ID,
		Name:      defaultListName,
		Position:  1,
		Audit: model.Audit{
			CreatedBy: project.CreatedBy,
			UpdatedBy: project.UpdatedBy,
		},
	}

	err = p.listRepo.Create(ctx, tx, defaultList)
	if err != nil {
		return err
	}

	return tx.Commit()
}
