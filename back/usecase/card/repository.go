package card

//go:generate mockgen -source=repository.go -destination=mock/repository_mock.go -package=cardmock

import (
	"context"

	"github.com/topi0247/hitori/domain/card"
)

type Repository interface {
	Save(ctx context.Context, c *card.Card) error
	FetchByID(ctx context.Context, id int64) (*card.Card, error)
	CountByThemeID(ctx context.Context, themeID int64) (int, error)
	GetAvailableNumber(ctx context.Context, themeID int64) (int, error)
	GetGameCards(ctx context.Context, themeID int64, amount int) ([]*card.Card, error)
	Confirm(ctx context.Context, id int64, word string) error
	Delete(ctx context.Context, id int64) error
}
