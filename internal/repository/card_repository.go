package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/MCPutro/go-management-project/utils"
	"time"

	"github.com/MCPutro/go-management-project/internal/model"
)

type CardRepository interface {
	Create(ctx context.Context, tx *sql.Tx, card *model.Card) error
	GetByID(ctx context.Context, tx *sql.Tx, id int64) (*model.Card, error)
	GetByListID(ctx context.Context, tx *sql.Tx, listID int64) ([]*model.Card, error)
	Update(ctx context.Context, tx *sql.Tx, card *model.Card) error
	Delete(ctx context.Context, tx *sql.Tx, id, deletedBy int64) error
}

type cardRepository struct {
}

func NewCardRepository() CardRepository {
	return &cardRepository{}
}

func (r *cardRepository) Create(ctx context.Context, tx *sql.Tx, card *model.Card) error {
	query := `
		INSERT INTO cards (list_id, title, content, position, created_at, created_by, updated_at, updated_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id
	`
	now := time.Now()
	return tx.QueryRowContext(ctx, query,
		card.ListID, card.Title, card.Content, card.Position,
		now, card.CreatedBy,
		now, card.UpdatedBy,
	).Scan(&card.ID)
}

func (r *cardRepository) GetByID(ctx context.Context, tx *sql.Tx, id int64) (*model.Card, error) {
	query := `SELECT id, list_id, title, content, position, created_at, created_by, updated_at, updated_by, deleted_at FROM cards WHERE id = $1 AND deleted_at IS NULL`
	row := tx.QueryRowContext(ctx, query, id)

	var card model.Card
	var deletedAt sql.NullTime

	err := row.Scan(
		&card.ID, &card.ListID, &card.Title, &card.Content, &card.Position,
		&card.CreatedAt, &card.CreatedBy, &card.UpdatedAt, &card.UpdatedBy,
		&deletedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, utils.ErrNotFound
	}

	if err != nil {
		return nil, err
	}

	if deletedAt.Valid {
		card.DeletedAt = &deletedAt.Time
	}

	return &card, nil
}

func (r *cardRepository) GetByListID(ctx context.Context, tx *sql.Tx, listID int64) ([]*model.Card, error) {
	query := `SELECT id, list_id, title, content, position, created_at, created_by, updated_at, updated_by, deleted_at FROM cards WHERE list_id = $1 AND deleted_at IS NULL ORDER BY position ASC`
	rows, err := tx.QueryContext(ctx, query, listID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cards []*model.Card
	for rows.Next() {
		var card model.Card
		var deletedAt sql.NullTime

		err := rows.Scan(
			&card.ID, &card.ListID, &card.Title, &card.Content, &card.Position,
			&card.CreatedAt, &card.CreatedBy, &card.UpdatedAt, &card.UpdatedBy,
			&deletedAt,
		)
		if err != nil {
			return nil, err
		}
		if deletedAt.Valid {
			card.DeletedAt = &deletedAt.Time
		}
		cards = append(cards, &card)
	}

	if len(cards) == 0 {
		return nil, utils.ErrNotFound
	}

	return cards, nil
}

func (r *cardRepository) Update(ctx context.Context, tx *sql.Tx, card *model.Card) error {
	query := `
		UPDATE cards SET title = $1, content = $2, position = $3, updated_at = $4, updated_by = $5
		WHERE id = $6 AND deleted_at IS NULL
	`
	now := time.Now()
	_, err := tx.ExecContext(ctx, query,
		card.Title, card.Content, card.Position, now, card.UpdatedBy, card.ID,
	)
	return err
}

func (r *cardRepository) Delete(ctx context.Context, tx *sql.Tx, id, deletedBy int64) error {
	query := `UPDATE cards SET deleted_at = $1, updated_at = $2, updated_by = $3 WHERE id = $4 AND deleted_at IS NULL`
	now := time.Now()
	_, err := tx.ExecContext(ctx, query, now, now, deletedBy, id)
	return err
}
