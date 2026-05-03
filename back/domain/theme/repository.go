package theme

//go:generate mockgen -source=repository.go -destination=mock/repository_mock.go -package=thememock

import "context"

type Repository interface {
	FindAll(ctx context.Context) ([]*Theme, error)
	FindByID(ctx context.Context, id int64) (*Theme, error)
}
