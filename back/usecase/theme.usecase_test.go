package usecase_test

import (
	"context"
	"testing"

	"go.uber.org/mock/gomock"

	domainTheme "github.com/topi0247/hitori/domain/theme"
	"github.com/topi0247/hitori/usecase"
	repositorymock "github.com/topi0247/hitori/usecase/repository/mock"
)

func TestList_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	themeRepository := repositorymock.NewMockThemeRepository(ctrl)
	ctx := context.Background()

	themeRepository.EXPECT().FetchAll(ctx).Return([]*domainTheme.Theme{
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
