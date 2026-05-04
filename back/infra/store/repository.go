package store

import (
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/topi0247/hitori/infra/db/sqlcgen"
	"github.com/topi0247/hitori/usecase/repository"
)

func NewRepositories(pool *pgxpool.Pool) *repository.Repositories {
	q := sqlcgen.New(pool)
	return &repository.Repositories{
		Card:       NewCardRepository(q),
		Theme:      NewThemeRepository(q),
		PlayRecord: NewPlayRecordRepository(q),
		Profile:    NewProfileRepository(q),
	}
}
