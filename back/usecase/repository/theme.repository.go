package repository

//go:generate mockgen -source=theme.repository.go -destination=mock/theme_mock.go -package=repositorymock

import (
	"context"

	domainTheme "github.com/topi0247/hitori/domain/theme"
)

type ThemeRepository interface {
	FetchAll(ctx context.Context) ([]*domainTheme.Theme, error)
	FetchByID(ctx context.Context, id int64) (*domainTheme.Theme, error)
}
