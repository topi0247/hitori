package repository

//go:generate mockgen -source=card.repository.go -destination=mock/card_mock.go -package=repositorymock

import (
	"context"

	"github.com/topi0247/hitori/domain"
)

type CardRepository interface {
	Save(ctx context.Context, c *domain.Card) error
	FetchByID(ctx context.Context, id int64) (*domain.Card, error)
	CountByThemeID(ctx context.Context, themeID int64) (int, error)
	GetAvailableNumber(ctx context.Context, themeID int64) (int, error)
	GetGameCards(ctx context.Context, themeID int64, amount int) ([]*domain.Card, error)
	Confirm(ctx context.Context, id int64, word string) error
	Delete(ctx context.Context, id int64) error
	FetchGameCardsByUUIDs(ctx context.Context, uuids []string) ([]*domain.Card, error)
	AddMatchPoints(ctx context.Context, uuid string, points int) error
}
