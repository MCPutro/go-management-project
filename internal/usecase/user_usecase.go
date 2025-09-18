package usecase

import (
	"context"
	"database/sql"

	"github.com/MCPutro/go-management-project/internal/model"
	"github.com/MCPutro/go-management-project/internal/repository"
)

type UserUsecase interface {
	CreateUser(ctx context.Context, user *model.User) error
	GetUserByID(ctx context.Context, id int64) (*model.User, error)
	UpdateUser(ctx context.Context, user *model.User) error
	DeleteUser(ctx context.Context, id int64, deletedBy int64) error
}

type userUsecase struct {
	db       *sql.DB
	userRepo repository.UserRepository
}

func NewUserUsecase(db *sql.DB, userRepository repository.UserRepository) UserUsecase {
	return &userUsecase{db: db, userRepo: userRepository}
}

func (u *userUsecase) CreateUser(ctx context.Context, user *model.User) error {
	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	err = u.userRepo.Create(ctx, tx, user)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (u *userUsecase) GetUserByID(ctx context.Context, id int64) (*model.User, error) {
	tx, err := u.db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback() // tidak perlu commit untuk read-only

	user, err := u.userRepo.GetByID(ctx, tx, id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *userUsecase) UpdateUser(ctx context.Context, user *model.User) error {
	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	err = u.userRepo.Update(ctx, tx, user)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (u *userUsecase) DeleteUser(ctx context.Context, id int64, deletedBy int64) error {
	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	err = u.userRepo.Delete(ctx, tx, id, deletedBy)
	if err != nil {
		return err
	}

	return tx.Commit()
}
