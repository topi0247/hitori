package store

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	domainTheme "github.com/topi0247/hitori/domain/theme"
	"github.com/topi0247/hitori/infra/db/sqlcgen"
)

type ThemeRepository struct {
	q *sqlcgen.Queries
}

func NewThemeRepository(q *sqlcgen.Queries) *ThemeRepository {
	return &ThemeRepository{q: q}
}

func (r *ThemeRepository) FetchAll(ctx context.Context) ([]*domainTheme.Theme, error) {
	rows, err := r.q.GetAllThemes(ctx)
	if err != nil {
		return nil, err
	}
	themes := make([]*domainTheme.Theme, len(rows))
	for i, row := range rows {
		themes[i] = &domainTheme.Theme{ID: row.ID, Title: row.Title}
	}
	return themes, nil
}

func (r *ThemeRepository) FetchByID(ctx context.Context, id int64) (*domainTheme.Theme, error) {
	row, err := r.q.GetThemeByID(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domainTheme.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &domainTheme.Theme{ID: row.ID, Title: row.Title}, nil
}
