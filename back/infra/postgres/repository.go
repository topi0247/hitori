package postgres

import (
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/topi0247/hitori/usecase/repository"
)

func NewRepositories(pool *pgxpool.Pool) *repository.Repositories {
	return &repository.Repositories{
		Card:       NewCardRepository(pool),
		Theme:      NewThemeRepository(pool),
		PlayRecord: NewPlayRecordRepository(pool),
	}
}
