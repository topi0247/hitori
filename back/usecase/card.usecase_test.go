package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"go.uber.org/mock/gomock"

	domainCard "github.com/topi0247/hitori/domain/card"
	domainTheme "github.com/topi0247/hitori/domain/theme"
	"github.com/topi0247/hitori/usecase"
	repositorymock "github.com/topi0247/hitori/usecase/repository/mock"
)

func strPtr(s string) *string { return &s }

func newCardUsecase(ctrl *gomock.Controller) (*usecase.CardUsecase, *repositorymock.MockCardRepository, *repositorymock.MockThemeRepository) {
	cardRepository := repositorymock.NewMockCardRepository(ctrl)
	themeRepository := repositorymock.NewMockThemeRepository(ctrl)
	uc := usecase.NewCardUsecase(cardRepository, themeRepository)
	return uc, cardRepository, themeRepository
}

// --- Create ---

func TestCreate_Success_Guest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	uc, cardRepository, themeRepository := newCardUsecase(ctrl)
	ctx := context.Background()

	themeRepository.EXPECT().FetchByID(ctx, int64(1)).Return(&domainTheme.Theme{Title: "大きさ"}, nil)
	cardRepository.EXPECT().CountByThemeID(ctx, int64(1)).Return(0, nil)
	cardRepository.EXPECT().Save(ctx, gomock.Any()).Return(nil)

	_, err := uc.Create(ctx, usecase.CreateInput{
		ThemeID:    1,
		CardNumber: 42,
		Word:       "アリ",
		GuestName:  strPtr("ゲスト"),
		Now:        time.Now(),
	})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestCreate_ThemeNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	uc, _, themeRepository := newCardUsecase(ctrl)
	ctx := context.Background()

	themeRepository.EXPECT().FetchByID(ctx, int64(1)).Return(nil, domainTheme.ErrNotFound)

	_, err := uc.Create(ctx, usecase.CreateInput{
		ThemeID:    1,
		CardNumber: 42,
		Word:       "アリ",
		GuestName:  strPtr("ゲスト"),
		Now:        time.Now(),
	})
	if err == nil {
		t.Error("テーマが存在しない場合はエラーを返すべき")
	}
}

func TestCreate_CardLimitReached(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	uc, cardRepository, themeRepository := newCardUsecase(ctrl)
	ctx := context.Background()

	themeRepository.EXPECT().FetchByID(ctx, int64(1)).Return(&domainTheme.Theme{Title: "大きさ"}, nil)
	cardRepository.EXPECT().CountByThemeID(ctx, int64(1)).Return(domainCard.MaxCardsPerTheme, nil)

	_, err := uc.Create(ctx, usecase.CreateInput{
		ThemeID:    1,
		CardNumber: 42,
		Word:       "アリ",
		GuestName:  strPtr("ゲスト"),
		Now:        time.Now(),
	})
	if err == nil {
		t.Error("上限に達している場合はエラーを返すべき")
	}
}

func TestCreate_InvalidCard(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	uc, cardRepository, themeRepository := newCardUsecase(ctrl)
	ctx := context.Background()

	themeRepository.EXPECT().FetchByID(ctx, int64(1)).Return(&domainTheme.Theme{Title: "大きさ"}, nil)
	cardRepository.EXPECT().CountByThemeID(ctx, int64(1)).Return(0, nil)

	_, err := uc.Create(ctx, usecase.CreateInput{
		ThemeID:    1,
		CardNumber: 0,
		Word:       "アリ",
		GuestName:  strPtr("ゲスト"),
		Now:        time.Now(),
	})
	if err == nil {
		t.Error("無効なカードの場合はエラーを返すべき")
	}
}

// --- Confirm ---

func TestConfirm_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	uc, cardRepository, _ := newCardUsecase(ctrl)
	ctx := context.Background()

	expiresAt := time.Now().Add(domainCard.ProvisionalTTL)
	existing := &domainCard.Card{
		ID:        1,
		Number:    42,
		Word:      "仮の言葉",
		GuestName: strPtr("ゲスト"),
		ExpiresAt: &expiresAt,
	}
	cardRepository.EXPECT().FetchByID(ctx, int64(1)).Return(existing, nil)
	cardRepository.EXPECT().Confirm(ctx, int64(1), "確定した言葉").Return(nil)

	_, err := uc.Confirm(ctx, usecase.ConfirmInput{ID: 1, Word: "確定した言葉"})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestConfirm_CardNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	uc, cardRepository, _ := newCardUsecase(ctrl)
	ctx := context.Background()

	cardRepository.EXPECT().FetchByID(ctx, int64(1)).Return(nil, errors.New("not_found"))

	_, err := uc.Confirm(ctx, usecase.ConfirmInput{ID: 1, Word: "確定した言葉"})
	if err == nil {
		t.Error("カードが存在しない場合はエラーを返すべき")
	}
}

func TestConfirm_AlreadyConfirmed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	uc, cardRepository, _ := newCardUsecase(ctrl)
	ctx := context.Background()

	existing := &domainCard.Card{ID: 1, Number: 42, Word: "確定済み", IsConfirmed: true}
	cardRepository.EXPECT().FetchByID(ctx, int64(1)).Return(existing, nil)

	_, err := uc.Confirm(ctx, usecase.ConfirmInput{ID: 1, Word: "再確定"})
	if !errors.Is(err, domainCard.ErrAlreadyConfirmed) {
		t.Errorf("wantErr=%v, got err=%v", domainCard.ErrAlreadyConfirmed, err)
	}
}

// --- Delete ---

func TestDelete_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	uc, cardRepository, _ := newCardUsecase(ctrl)
	ctx := context.Background()

	existing := &domainCard.Card{ID: 1, Number: 42, Word: "アリ", GuestName: strPtr("ゲスト")}
	cardRepository.EXPECT().FetchByID(ctx, int64(1)).Return(existing, nil)
	cardRepository.EXPECT().Delete(ctx, int64(1)).Return(nil)

	if err := uc.Delete(ctx, 1); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestDelete_CardNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	uc, cardRepository, _ := newCardUsecase(ctrl)
	ctx := context.Background()

	cardRepository.EXPECT().FetchByID(ctx, int64(1)).Return(nil, errors.New("not_found"))

	if err := uc.Delete(ctx, 1); err == nil {
		t.Error("カードが存在しない場合はエラーを返すべき")
	}
}

func TestDelete_AlreadyConfirmed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	uc, cardRepository, _ := newCardUsecase(ctrl)
	ctx := context.Background()

	existing := &domainCard.Card{ID: 1, Number: 42, Word: "確定済み", IsConfirmed: true}
	cardRepository.EXPECT().FetchByID(ctx, int64(1)).Return(existing, nil)

	err := uc.Delete(ctx, 1)
	if !errors.Is(err, domainCard.ErrAlreadyConfirmed) {
		t.Errorf("wantErr=%v, got err=%v", domainCard.ErrAlreadyConfirmed, err)
	}
}

// --- Available ---

func TestAvailable_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	uc, cardRepository, themeRepository := newCardUsecase(ctrl)
	ctx := context.Background()

	themeRepository.EXPECT().FetchByID(ctx, int64(1)).Return(&domainTheme.Theme{Title: "大きさ"}, nil)
	cardRepository.EXPECT().GetAvailableNumber(ctx, int64(1)).Return(42, nil)

	out, err := uc.Available(ctx, int64(1))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.CardNumber != 42 {
		t.Errorf("CardNumber=%d, want 42", out.CardNumber)
	}
}

func TestAvailable_ThemeNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	uc, _, themeRepository := newCardUsecase(ctrl)
	ctx := context.Background()

	themeRepository.EXPECT().FetchByID(ctx, int64(1)).Return(nil, domainTheme.ErrNotFound)

	_, err := uc.Available(ctx, int64(1))
	if err == nil {
		t.Error("テーマが存在しない場合はエラーを返すべき")
	}
}

// --- GameCards ---

func TestGameCards_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	uc, cardRepository, themeRepository := newCardUsecase(ctrl)
	ctx := context.Background()

	themeRepository.EXPECT().FetchByID(ctx, int64(1)).Return(&domainTheme.Theme{Title: "大きさ"}, nil)
	cardRepository.EXPECT().GetGameCards(ctx, int64(1), 4).Return([]*domainCard.Card{
		{UUID: "uuid-a", Word: "アリ"},
		{UUID: "uuid-b", Word: "クジラ"},
		{UUID: "uuid-c", Word: "ゾウ"},
	}, nil)

	out, err := uc.GameCards(ctx, usecase.GameCardsInput{ThemeID: 1, CardAmount: 4})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Cards) != 3 {
		t.Errorf("len=%d, want 3", len(out.Cards))
	}
}

func TestGameCards_InvalidAmount(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	uc, _, _ := newCardUsecase(ctrl)
	ctx := context.Background()

	_, err := uc.GameCards(ctx, usecase.GameCardsInput{ThemeID: 1, CardAmount: 3})
	if err == nil {
		t.Error("枚数が4未満の場合はエラーを返すべき")
	}
}
