package store

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	domainCard "github.com/topi0247/hitori/domain/card"
	"github.com/topi0247/hitori/infra/db/sqlcgen"
)

type CardRepository struct {
	q *sqlcgen.Queries
}

func NewCardRepository(q *sqlcgen.Queries) *CardRepository {
	return &CardRepository{q: q}
}

func toCard(m sqlcgen.Card) *domainCard.Card {
	return &domainCard.Card{
		ID:          m.ID,
		UUID:        m.Uuid,
		ThemeID:     m.ThemeID,
		Number:      int(m.CardNumber),
		Word:        m.Word,
		ProfileID:   fromPgtypeInt8(m.ProfileID),
		GuestName:   fromPgtypeText(m.GuestName),
		IsConfirmed: m.IsConfirmed,
		MatchPoints: int(m.MatchPoints),
		ExpiresAt:   fromPgtypeTimestamptz(m.ExpiresAt),
	}
}

func (r *CardRepository) Save(ctx context.Context, c *domainCard.Card) error {
	row, err := r.q.InsertCard(ctx, sqlcgen.InsertCardParams{
		ThemeID:    c.ThemeID,
		ProfileID:  toPgtypeInt8(c.ProfileID),
		GuestName:  toPgtypeText(c.GuestName),
		CardNumber: int16(c.Number),
		Word:       c.Word,
		ExpiresAt:  toPgtypeTimestamptz(c.ExpiresAt),
	})
	if err != nil {
		return err
	}
	c.ID = row.ID
	c.UUID = row.Uuid
	return nil
}

func (r *CardRepository) FetchByID(ctx context.Context, id int64) (*domainCard.Card, error) {
	m, err := r.q.GetCardByID(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domainCard.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return toCard(m), nil
}

func (r *CardRepository) CountByThemeID(ctx context.Context, themeID int64) (int, error) {
	count, err := r.q.CountCardsByThemeID(ctx, themeID)
	return int(count), err
}

func (r *CardRepository) GetAvailableNumber(ctx context.Context, themeID int64) (int, error) {
	n, err := r.q.GetAvailableCardNumber(ctx, sqlcgen.GetAvailableCardNumberParams{
		Column1: int32(domainCard.MinCardNumber),
		Column2: int32(domainCard.MaxCardNumber),
		ThemeID: themeID,
	})
	if errors.Is(err, pgx.ErrNoRows) {
		return 0, domainCard.ErrThemeCardLimitReached
	}
	return int(n), err
}

func (r *CardRepository) GetGameCards(ctx context.Context, themeID int64, amount int) ([]*domainCard.Card, error) {
	rows, err := r.q.GetGameCards(ctx, sqlcgen.GetGameCardsParams{
		ThemeID: themeID,
		Column2: int32(amount),
	})
	if err != nil {
		return nil, err
	}
	cards := make([]*domainCard.Card, len(rows))
	for i, m := range rows {
		cards[i] = toCard(m)
	}
	return cards, nil
}

func (r *CardRepository) Confirm(ctx context.Context, id int64, word string) error {
	return r.q.ConfirmCard(ctx, sqlcgen.ConfirmCardParams{ID: id, Word: word})
}

func (r *CardRepository) Delete(ctx context.Context, id int64) error {
	return r.q.DeleteCard(ctx, id)
}

func (r *CardRepository) FetchGameCardsByUUIDs(ctx context.Context, uuids []string) ([]*domainCard.Card, error) {
	rows, err := r.q.GetCardsByUUIDs(ctx, uuids)
	if err != nil {
		return nil, err
	}
	cards := make([]*domainCard.Card, len(rows))
	for i, m := range rows {
		cards[i] = toCard(m)
	}
	return cards, nil
}

func (r *CardRepository) AddMatchPoints(ctx context.Context, uuid string, points int) error {
	return r.q.AddMatchPoints(ctx, sqlcgen.AddMatchPointsParams{
		Column1: int32(points),
		Uuid:    uuid,
	})
}
