package usecase

import (
	"context"

	"github.com/topi0247/hitori/domain"
)

type ThemeUsecase struct {
	repository domain.ThemeRepository
}

func NewThemeUsecase(repository domain.ThemeRepository) *ThemeUsecase {
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
