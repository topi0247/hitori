package usecase

import (
	"context"

	"github.com/topi0247/hitori/domain"
	"github.com/topi0247/hitori/usecase/repository"
)

type ThemeUsecase struct {
	repository repository.ThemeRepository
}

func NewThemeUsecase(repository repository.ThemeRepository) *ThemeUsecase {
	return &ThemeUsecase{repository: repository}
}

type ListOutput struct {
	Themes []*domain.Theme
}

func (u *ThemeUsecase) List(ctx context.Context) (*ListOutput, error) {
	themes, err := u.repository.FetchAll(ctx)
	if err != nil {
		return nil, err
	}
	return &ListOutput{Themes: themes}, nil
}
