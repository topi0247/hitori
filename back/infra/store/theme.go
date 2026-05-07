package store

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	"github.com/topi0247/hitori/domain"
	"github.com/topi0247/hitori/infra/db/sqlcgen"
)

type ThemeRepository struct {
	q *sqlcgen.Queries
}

func NewThemeRepository(q *sqlcgen.Queries) *ThemeRepository {
	return &ThemeRepository{q: q}
}

func (r *ThemeRepository) FetchAll(ctx context.Context) ([]*domain.Theme, error) {
	rows, err := r.q.GetAllThemes(ctx)
	if err != nil {
		return nil, err
	}
	themes := make([]*domain.Theme, len(rows))
	for i, row := range rows {
		themes[i] = &domain.Theme{ID: row.ID, Title: row.Title}
	}
	return themes, nil
}

func (r *ThemeRepository) FetchByID(ctx context.Context, id int64) (*domain.Theme, error) {
	row, err := r.q.GetThemeByID(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &domain.Theme{ID: row.ID, Title: row.Title}, nil
}
