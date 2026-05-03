package theme

//go:generate mockgen -source=repository.go -destination=mock/repository_mock.go -package=thememock

import (
	"context"

	"github.com/topi0247/hitori/domain/theme"
)

type Repository interface {
	FetchAll(ctx context.Context) ([]*theme.Theme, error)
	FetchByID(ctx context.Context, id int64) (*theme.Theme, error)
}
