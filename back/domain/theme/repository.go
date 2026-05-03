package theme

import "context"

type Repository interface {
	FindAll(ctx context.Context) ([]*Theme, error)
	FindByID(ctx context.Context, id int64) (*Theme, error)
}
