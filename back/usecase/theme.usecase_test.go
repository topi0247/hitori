package usecase_test

import (
	"context"
	"testing"

	"go.uber.org/mock/gomock"

	"github.com/topi0247/hitori/domain"
	"github.com/topi0247/hitori/usecase"
	domainmock "github.com/topi0247/hitori/domain/mock"
)

func TestList_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	themeRepository := domainmock.NewMockThemeRepository(ctrl)
	ctx := context.Background()

	themeRepository.EXPECT().FetchAll(ctx).Return([]*domain.Theme{
		{Title: "大きさ"},
		{Title: "速さ"},
	}, nil)

	uc := usecase.NewThemeUsecase(themeRepository)
	out, err := uc.List(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Themes) != 2 {
		t.Errorf("len=%d, want 2", len(out.Themes))
	}
}
