// internal/usecase/card_usecase.go
package usecase

import (
	"context"
	"database/sql"

	"github.com/MCPutro/go-management-project/internal/model"
	"github.com/MCPutro/go-management-project/internal/repository"
)

type CardUsecase interface {
	CreateCard(ctx context.Context, card *model.Card) error
	GetCardsByListID(ctx context.Context, listID int64) ([]*model.Card, error)
	GetCardByID(ctx context.Context, id int64) (*model.Card, error)
	UpdateCard(ctx context.Context, card *model.Card) error
	DeleteCard(ctx context.Context, id int64, deletedBy int64) error
}

type cardUsecase struct {
	db       *sql.DB
	cardRepo repository.CardRepository
}

func NewCardUsecase(db *sql.DB, cardRepository repository.CardRepository) CardUsecase {
	return &cardUsecase{
		db:       db,
		cardRepo: cardRepository,
	}
}

func (c *cardUsecase) CreateCard(ctx context.Context, card *model.Card) error {
	tx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	err = c.cardRepo.Create(ctx, tx, card)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (c *cardUsecase) GetCardsByListID(ctx context.Context, listID int64) ([]*model.Card, error) {
	tx, err := c.db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	cards, err := c.cardRepo.GetByListID(ctx, tx, listID)
	if err != nil {
		return nil, err
	}

	return cards, nil
}

func (c *cardUsecase) GetCardByID(ctx context.Context, id int64) (*model.Card, error) {
	tx, err := c.db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	query := `SELECT id, list_id, title, content, position, created_at, created_by, updated_at, updated_by, deleted_at FROM cards WHERE id = $1 AND deleted_at IS NULL`
	row := tx.QueryRowContext(ctx, query, id)

	var card model.Card
	err = row.Scan(
		&card.ID, &card.ListID, &card.Title, &card.Content, &card.Position,
		&card.CreatedAt, &card.CreatedBy, &card.UpdatedAt, &card.UpdatedBy, &card.DeletedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &card, nil
}

func (c *cardUsecase) UpdateCard(ctx context.Context, card *model.Card) error {
	tx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	err = c.cardRepo.Update(ctx, tx, card)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (c *cardUsecase) DeleteCard(ctx context.Context, id int64, deletedBy int64) error {
	tx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	err = c.cardRepo.Delete(ctx, tx, id, deletedBy)
	if err != nil {
		return err
	}

	return tx.Commit()
}
