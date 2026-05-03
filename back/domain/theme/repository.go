package theme

//go:generate mockgen -source=repository.go -destination=mock/repository_mock.go -package=thememock

import "context"

type Repository interface {
	FetchAll(ctx context.Context) ([]*Theme, error)
	FetchByID(ctx context.Context, id int64) (*Theme, error)
}
