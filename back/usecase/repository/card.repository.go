package repository

//go:generate mockgen -source=card.repository.go -destination=mock/card_mock.go -package=repositorymock

import (
	"context"

	domainCard "github.com/topi0247/hitori/domain/card"
)

type CardRepository interface {
	Save(ctx context.Context, c *domainCard.Card) error
	FetchByID(ctx context.Context, id int64) (*domainCard.Card, error)
	CountByThemeID(ctx context.Context, themeID int64) (int, error)
	GetAvailableNumber(ctx context.Context, themeID int64) (int, error)
	GetGameCards(ctx context.Context, themeID int64, amount int) ([]*domainCard.Card, error)
	Confirm(ctx context.Context, id int64, word string) error
	Delete(ctx context.Context, id int64) error
	FetchGameCardsByUUIDs(ctx context.Context, uuids []string) ([]*domainCard.Card, error)
	AddMatchPoints(ctx context.Context, uuid string, points int) error
}
