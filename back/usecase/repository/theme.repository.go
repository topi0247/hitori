package repository

//go:generate mockgen -source=theme.repository.go -destination=mock/theme_mock.go -package=repositorymock

import (
	"context"

	"github.com/topi0247/hitori/domain"
)

type ThemeRepository interface {
	FetchAll(ctx context.Context) ([]*domain.Theme, error)
	FetchByID(ctx context.Context, id int64) (*domain.Theme, error)
}
