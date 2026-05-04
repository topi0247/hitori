package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	domainProfile "github.com/topi0247/hitori/domain/profile"
)

type ProfileRepository struct {
	pool *pgxpool.Pool
}

func NewProfileRepository(pool *pgxpool.Pool) *ProfileRepository {
	return &ProfileRepository{pool: pool}
}

func (r *ProfileRepository) FetchByAuthUserID(ctx context.Context, authUserID string) (*domainProfile.Profile, error) {
	p := &domainProfile.Profile{}
	err := r.pool.QueryRow(ctx, `
		SELECT id, auth_user_id, user_name
		FROM profiles WHERE auth_user_id = $1`, authUserID).
		Scan(&p.ID, &p.AuthUserID, &p.UserName)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domainProfile.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (r *ProfileRepository) UpdateUserName(ctx context.Context, authUserID string, userName string) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE profiles SET user_name = $1 WHERE auth_user_id = $2`,
		userName, authUserID)
	return err
}

func (r *ProfileRepository) Delete(ctx context.Context, authUserID string) error {
	_, err := r.pool.Exec(ctx, `
		DELETE FROM profiles WHERE auth_user_id = $1`, authUserID)
	return err
}
