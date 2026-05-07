package store

import (
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/topi0247/hitori/domain"
	"github.com/topi0247/hitori/infra/db/sqlcgen"
)

func NewRepositories(pool *pgxpool.Pool) *domain.Repositories {
	q := sqlcgen.New(pool)
	return &domain.Repositories{
		Card:       NewCardRepository(q),
		Theme:      NewThemeRepository(q),
		PlayRecord: NewPlayRecordRepository(q),
		Profile:    NewProfileRepository(q),
	}
}
