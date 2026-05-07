package domain

//go:generate mockgen -source=profile.go -destination=mock/profile_mock.go -package=domainmock

import (
	"context"
	"unicode/utf8"
)

const MaxUserNameLength = 10

type Profile struct {
	ID         int64
	AuthUserID string
	UserName   string
}

func NewProfile(authUserID, userName string) (*Profile, error) {
	if utf8.RuneCountInString(userName) == 0 || utf8.RuneCountInString(userName) > MaxUserNameLength {
		return nil, ErrInvalidUserName
	}
	return &Profile{AuthUserID: authUserID, UserName: userName}, nil
}

func (p *Profile) UpdateUserName(name string) error {
	if utf8.RuneCountInString(name) == 0 || utf8.RuneCountInString(name) > MaxUserNameLength {
		return ErrInvalidUserName
	}
	p.UserName = name
	return nil
}

type ProfileRepository interface {
	Create(ctx context.Context, authUserID string, userName string) (*Profile, error)
	FetchByAuthUserID(ctx context.Context, authUserID string) (*Profile, error)
	UpdateUserName(ctx context.Context, authUserID string, userName string) error
	Delete(ctx context.Context, authUserID string) error
}
