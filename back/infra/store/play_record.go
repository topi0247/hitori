package store

import (
	"context"

	"github.com/topi0247/hitori/domain"
	"github.com/topi0247/hitori/infra/db/sqlcgen"
)

type PlayRecordRepository struct {
	q *sqlcgen.Queries
}

func NewPlayRecordRepository(q *sqlcgen.Queries) *PlayRecordRepository {
	return &PlayRecordRepository{q: q}
}

func (r *PlayRecordRepository) Save(ctx context.Context, record *domain.PlayRecord) error {
	return r.q.InsertPlayRecord(ctx, sqlcgen.InsertPlayRecordParams{
		ThemeID:     record.ThemeID,
		ProfileID:   record.ProfileID,
		CardAmount:  int16(record.CardAmount),
		CorrectRate: record.CorrectRate,
	})
}
