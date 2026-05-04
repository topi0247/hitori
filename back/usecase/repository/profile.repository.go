package repository

//go:generate mockgen -source=profile.repository.go -destination=mock/profile_mock.go -package=repositorymock

import (
	"context"

	domainProfile "github.com/topi0247/hitori/domain/profile"
)

type ProfileRepository interface {
	Create(ctx context.Context, authUserID string, userName string) (*domainProfile.Profile, error)
	FetchByAuthUserID(ctx context.Context, authUserID string) (*domainProfile.Profile, error)
	UpdateUserName(ctx context.Context, authUserID string, userName string) error
	Delete(ctx context.Context, authUserID string) error
}
