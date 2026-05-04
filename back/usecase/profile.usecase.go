package usecase

import (
	"context"
	"errors"

	domainProfile "github.com/topi0247/hitori/domain/profile"
	"github.com/topi0247/hitori/usecase/repository"
)

type ProfileUsecase struct {
	repository repository.ProfileRepository
}

func NewProfileUsecase(repository repository.ProfileRepository) *ProfileUsecase {
	return &ProfileUsecase{repository: repository}
}

// --- Create ---

type CreateProfileInput struct {
	AuthUserID string
	UserName   string
}

type CreateProfileOutput struct {
	UserName string
}

func (u *ProfileUsecase) Create(ctx context.Context, input CreateProfileInput) (*CreateProfileOutput, error) {
	_, err := u.repository.FetchByAuthUserID(ctx, input.AuthUserID)
	if err == nil {
		return nil, domainProfile.ErrAlreadyExists
	}
	if !errors.Is(err, domainProfile.ErrNotFound) {
		return nil, err
	}

	p, err := domainProfile.New(input.AuthUserID, input.UserName)
	if err != nil {
		return nil, err
	}

	created, err := u.repository.Create(ctx, p.AuthUserID, p.UserName)
	if err != nil {
		return nil, err
	}
	return &CreateProfileOutput{UserName: created.UserName}, nil
}

// --- Get ---

type GetProfileOutput struct {
	UserName string
}

func (u *ProfileUsecase) Get(ctx context.Context, authUserID string) (*GetProfileOutput, error) {
	p, err := u.repository.FetchByAuthUserID(ctx, authUserID)
	if err != nil {
		return nil, err
	}
	return &GetProfileOutput{UserName: p.UserName}, nil
}

// --- Update ---

type UpdateProfileInput struct {
	AuthUserID string
	UserName   string
}

type UpdateProfileOutput struct {
	UserName string
}

func (u *ProfileUsecase) Update(ctx context.Context, input UpdateProfileInput) (*UpdateProfileOutput, error) {
	p, err := u.repository.FetchByAuthUserID(ctx, input.AuthUserID)
	if err != nil {
		return nil, err
	}

	if err := p.UpdateUserName(input.UserName); err != nil {
		return nil, err
	}

	if err := u.repository.UpdateUserName(ctx, input.AuthUserID, p.UserName); err != nil {
		return nil, err
	}

	return &UpdateProfileOutput{UserName: p.UserName}, nil
}

// --- Delete ---

func (u *ProfileUsecase) Delete(ctx context.Context, authUserID string) error {
	if _, err := u.repository.FetchByAuthUserID(ctx, authUserID); err != nil {
		return err
	}
	return u.repository.Delete(ctx, authUserID)
}

