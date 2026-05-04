package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/topi0247/hitori/usecase/repository"
)

type PlayRecordRepository struct {
	pool *pgxpool.Pool
}

func NewPlayRecordRepository(pool *pgxpool.Pool) *PlayRecordRepository {
	return &PlayRecordRepository{pool: pool}
}

func (r *PlayRecordRepository) Save(ctx context.Context, record *repository.PlayRecord) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO play_records (theme_id, profile_id, card_amount, correct_rate)
		VALUES ($1, $2, $3, $4)`,
		record.ThemeID, record.ProfileID, record.CardAmount, record.CorrectRate,
	)
	return err
}
