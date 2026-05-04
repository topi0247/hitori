package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	domainCard "github.com/topi0247/hitori/domain/card"
)

type CardRepository struct {
	pool *pgxpool.Pool
}

func NewCardRepository(pool *pgxpool.Pool) *CardRepository {
	return &CardRepository{pool: pool}
}

func (r *CardRepository) Save(ctx context.Context, c *domainCard.Card) error {
	return r.pool.QueryRow(ctx, `
		INSERT INTO cards (theme_id, profile_id, guest_name, card_number, word, expires_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, uuid`,
		c.ThemeID, c.ProfileID, c.GuestName, c.Number, c.Word, c.ExpiresAt,
	).Scan(&c.ID, &c.UUID)
}

func (r *CardRepository) FetchByID(ctx context.Context, id int64) (*domainCard.Card, error) {
	c := &domainCard.Card{}
	err := r.pool.QueryRow(ctx, `
		SELECT id, uuid, theme_id, card_number, word, profile_id, guest_name,
		       is_confirmed, match_points, expires_at
		FROM cards WHERE id = $1`, id).
		Scan(&c.ID, &c.UUID, &c.ThemeID, &c.Number, &c.Word,
			&c.ProfileID, &c.GuestName, &c.IsConfirmed, &c.MatchPoints, &c.ExpiresAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domainCard.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (r *CardRepository) CountByThemeID(ctx context.Context, themeID int64) (int, error) {
	var count int
	err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM cards WHERE theme_id = $1`, themeID).
		Scan(&count)
	return count, err
}

func (r *CardRepository) GetAvailableNumber(ctx context.Context, themeID int64) (int, error) {
	var number int
	err := r.pool.QueryRow(ctx, `
		SELECT n FROM generate_series($1::int, $2::int) AS n
		WHERE n NOT IN (SELECT card_number FROM cards WHERE theme_id = $3)
		ORDER BY random()
		LIMIT 1`,
		domainCard.MinCardNumber, domainCard.MaxCardNumber, themeID,
	).Scan(&number)
	if errors.Is(err, pgx.ErrNoRows) {
		return 0, domainCard.ErrThemeCardLimitReached
	}
	return number, err
}

func (r *CardRepository) GetGameCards(ctx context.Context, themeID int64, amount int) ([]*domainCard.Card, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, uuid, theme_id, card_number, word, profile_id, guest_name,
		       is_confirmed, match_points, expires_at
		FROM cards
		WHERE theme_id = $1 AND is_confirmed = true
		ORDER BY random()
		LIMIT $2`, themeID, amount)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cards []*domainCard.Card
	for rows.Next() {
		c := &domainCard.Card{}
		if err := rows.Scan(&c.ID, &c.UUID, &c.ThemeID, &c.Number, &c.Word,
			&c.ProfileID, &c.GuestName, &c.IsConfirmed, &c.MatchPoints, &c.ExpiresAt); err != nil {
			return nil, err
		}
		cards = append(cards, c)
	}
	return cards, rows.Err()
}

func (r *CardRepository) Confirm(ctx context.Context, id int64, word string) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE cards SET word = $1, is_confirmed = true, expires_at = NULL, updated_at = now()
		WHERE id = $2`, word, id)
	return err
}

func (r *CardRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM cards WHERE id = $1`, id)
	return err
}

func (r *CardRepository) FetchGameCardsByUUIDs(ctx context.Context, uuids []string) ([]*domainCard.Card, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, uuid, theme_id, card_number, word, profile_id, guest_name,
		       is_confirmed, match_points, expires_at
		FROM cards WHERE uuid = ANY($1)`, uuids)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cards []*domainCard.Card
	for rows.Next() {
		c := &domainCard.Card{}
		if err := rows.Scan(&c.ID, &c.UUID, &c.ThemeID, &c.Number, &c.Word,
			&c.ProfileID, &c.GuestName, &c.IsConfirmed, &c.MatchPoints, &c.ExpiresAt); err != nil {
			return nil, err
		}
		cards = append(cards, c)
	}
	return cards, rows.Err()
}

func (r *CardRepository) AddMatchPoints(ctx context.Context, uuid string, points int) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE cards SET match_points = match_points + $1, updated_at = now()
		WHERE uuid = $2`, points, uuid)
	return err
}
