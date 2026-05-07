package usecase_test

import (
	"context"
	"testing"

	"go.uber.org/mock/gomock"

	"github.com/topi0247/hitori/domain"
	"github.com/topi0247/hitori/usecase"
	repositorymock "github.com/topi0247/hitori/usecase/repository/mock"
)

func TestPlay_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cardRepository := repositorymock.NewMockCardRepository(ctrl)
	playRecordRepository := repositorymock.NewMockPlayRecordRepository(ctrl)
	profileRepository := repositorymock.NewMockProfileRepository(ctrl)
	ctx := context.Background()

	profileRepository.EXPECT().FetchByAuthUserID(ctx, "user-123").Return(&domain.Profile{ID: 1, AuthUserID: "user-123", UserName: "テスト"}, nil)

	cards := []*domain.Card{
		{UUID: "a", Number: 10, Word: "アリ"},
		{UUID: "b", Number: 30, Word: "クジラ"},
		{UUID: "c", Number: 50, Word: "ゾウ"},
	}
	cardRepository.EXPECT().FetchGameCardsByUUIDs(ctx, []string{"a", "b", "c"}).Return(cards, nil)
	cardRepository.EXPECT().AddMatchPoints(ctx, "a", 3).Return(nil)
	cardRepository.EXPECT().AddMatchPoints(ctx, "b", 3).Return(nil)
	cardRepository.EXPECT().AddMatchPoints(ctx, "c", 3).Return(nil)
	playRecordRepository.EXPECT().Save(ctx, gomock.Any()).Return(nil)

	uc := usecase.NewGameUsecase(cardRepository, playRecordRepository, profileRepository)
	out, err := uc.Play(ctx, usecase.PlayInput{
		ThemeID:    1,
		AuthUserID: "user-123",
		Answers: []usecase.Answer{
			{UUID: "a", Order: 1},
			{UUID: "b", Order: 2},
			{UUID: "c", Order: 3},
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.CorrectRate != 100.0 {
		t.Errorf("CorrectRate=%v, want 100.0", out.CorrectRate)
	}
}

func TestPlay_NoAnswers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cardRepository := repositorymock.NewMockCardRepository(ctrl)
	playRecordRepository := repositorymock.NewMockPlayRecordRepository(ctrl)
	profileRepository := repositorymock.NewMockProfileRepository(ctrl)
	ctx := context.Background()

	uc := usecase.NewGameUsecase(cardRepository, playRecordRepository, profileRepository)
	_, err := uc.Play(ctx, usecase.PlayInput{
		ThemeID:    1,
		AuthUserID: "user-123",
		Answers:    []usecase.Answer{},
	})
	if err == nil {
		t.Error("回答が空の場合はエラーを返すべき")
	}
}
