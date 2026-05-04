package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	domainTheme "github.com/topi0247/hitori/domain/theme"
)

type ThemeRepository struct {
	pool *pgxpool.Pool
}

func NewThemeRepository(pool *pgxpool.Pool) *ThemeRepository {
	return &ThemeRepository{pool: pool}
}

func (r *ThemeRepository) FetchAll(ctx context.Context) ([]*domainTheme.Theme, error) {
	rows, err := r.pool.Query(ctx, `SELECT id, title FROM themes ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var themes []*domainTheme.Theme
	for rows.Next() {
		t := &domainTheme.Theme{}
		if err := rows.Scan(&t.ID, &t.Title); err != nil {
			return nil, err
		}
		themes = append(themes, t)
	}
	return themes, rows.Err()
}

func (r *ThemeRepository) FetchByID(ctx context.Context, id int64) (*domainTheme.Theme, error) {
	t := &domainTheme.Theme{}
	err := r.pool.QueryRow(ctx, `SELECT id, title FROM themes WHERE id = $1`, id).
		Scan(&t.ID, &t.Title)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domainTheme.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return t, nil
}
