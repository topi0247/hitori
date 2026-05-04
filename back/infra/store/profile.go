package store

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	"github.com/topi0247/hitori/infra/db/sqlcgen"
	domainProfile "github.com/topi0247/hitori/domain/profile"
)

type ProfileRepository struct {
	q *sqlcgen.Queries
}

func NewProfileRepository(q *sqlcgen.Queries) *ProfileRepository {
	return &ProfileRepository{q: q}
}

func (r *ProfileRepository) FetchByAuthUserID(ctx context.Context, authUserID string) (*domainProfile.Profile, error) {
	row, err := r.q.GetProfileByAuthUserID(ctx, authUserID)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domainProfile.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &domainProfile.Profile{
		ID:         row.ID,
		AuthUserID: row.AuthUserID,
		UserName:   row.UserName,
	}, nil
}

func (r *ProfileRepository) UpdateUserName(ctx context.Context, authUserID string, userName string) error {
	return r.q.UpdateProfileUserName(ctx, sqlcgen.UpdateProfileUserNameParams{
		AuthUserID: authUserID,
		UserName:   userName,
	})
}

func (r *ProfileRepository) Delete(ctx context.Context, authUserID string) error {
	return r.q.DeleteProfile(ctx, authUserID)
}
