package repository

//go:generate mockgen -source=profile.repository.go -destination=mock/profile_mock.go -package=repositorymock

import (
	"context"

	"github.com/topi0247/hitori/domain"
)

type ProfileRepository interface {
	Create(ctx context.Context, authUserID string, userName string) (*domain.Profile, error)
	FetchByAuthUserID(ctx context.Context, authUserID string) (*domain.Profile, error)
	UpdateUserName(ctx context.Context, authUserID string, userName string) error
	Delete(ctx context.Context, authUserID string) error
}
