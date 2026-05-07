package domain_test

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/topi0247/hitori/domain"
)

func TestNewCard_CardNumber(t *testing.T) {
	tests := []struct {
		name    string
		number  int
		wantErr error
	}{
		{"最小値", domain.MinCardNumber, nil},
		{"最大値", domain.MaxCardNumber, nil},
		{"最小値未満", domain.MinCardNumber - 1, domain.ErrInvalidCardNumber},
		{"最大値超過", domain.MaxCardNumber + 1, domain.ErrInvalidCardNumber},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := domain.NewCard(domain.NewCardInput{
				Number:    tt.number,
				Word:      "テスト",
				GuestName: strPtr("ゲスト"),
				Now:       time.Now(),
			})
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("wantErr=%v, got err=%v", tt.wantErr, err)
			}
		})
	}
}

func TestNewCard_Word(t *testing.T) {
	tests := []struct {
		name    string
		word    string
		wantErr error
	}{
		{"最大文字数", strings.Repeat("あ", domain.MaxWordLength), nil},
		{"最大文字数超過", strings.Repeat("あ", domain.MaxWordLength+1), domain.ErrInvalidWord},
		{"空文字", "", domain.ErrInvalidWord},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := domain.NewCard(domain.NewCardInput{
				Number:    domain.MinCardNumber,
				Word:      tt.word,
				GuestName: strPtr("ゲスト"),
				Now:       time.Now(),
			})
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("wantErr=%v, got err=%v", tt.wantErr, err)
			}
		})
	}
}

func TestNewCard_GuestName(t *testing.T) {
	tests := []struct {
		name      string
		guestName *string
		wantErr   error
	}{
		{"最大文字数", strPtr(strings.Repeat("あ", domain.MaxGuestNameLength)), nil},
		{"最大文字数超過", strPtr(strings.Repeat("あ", domain.MaxGuestNameLength+1)), domain.ErrInvalidGuestName},
		{"空文字", strPtr(""), domain.ErrInvalidGuestName},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := domain.NewCard(domain.NewCardInput{
				Number:    domain.MinCardNumber,
				Word:      "テスト",
				GuestName: tt.guestName,
				Now:       time.Now(),
			})
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("wantErr=%v, got err=%v", tt.wantErr, err)
			}
		})
	}
}

func TestNewCard_Owner(t *testing.T) {
	profileID := int64(1)

	tests := []struct {
		name      string
		profileID *int64
		guestName *string
		wantErr   error
	}{
		{"profileIDあり", &profileID, nil, nil},
		{"guestNameあり", nil, strPtr("ゲスト"), nil},
		{"両方nil", nil, nil, domain.ErrOwnerRequired},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := domain.NewCard(domain.NewCardInput{
				Number:    domain.MinCardNumber,
				Word:      "テスト",
				ProfileID: tt.profileID,
				GuestName: tt.guestName,
				Now:       time.Now(),
			})
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("wantErr=%v, got err=%v", tt.wantErr, err)
			}
		})
	}
}

func TestNewCard_InitialState(t *testing.T) {
	now := time.Now()
	c, err := domain.NewCard(domain.NewCardInput{
		Number:    domain.MinCardNumber,
		Word:      "テスト",
		GuestName: strPtr("ゲスト"),
		Now:       now,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if c.IsConfirmed {
		t.Error("新規カードは仮登録状態であるべき")
	}
	if c.MatchPoints != 0 {
		t.Errorf("初期マッチポイントは0であるべき, got %d", c.MatchPoints)
	}
	if c.ExpiresAt == nil {
		t.Error("仮登録カードはExpiresAtが設定されるべき")
	}
	want := now.Add(domain.ProvisionalTTL)
	if !c.ExpiresAt.Equal(want) {
		t.Errorf("ExpiresAt=%v, want %v", c.ExpiresAt, want)
	}
}

func TestCard_Confirm(t *testing.T) {
	c, _ := domain.NewCard(domain.NewCardInput{
		Number:    domain.MinCardNumber,
		Word:      "仮の言葉",
		GuestName: strPtr("ゲスト"),
		Now:       time.Now(),
	})

	err := c.Confirm("確定した言葉")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !c.IsConfirmed {
		t.Error("本登録後はIsConfirmedがtrueであるべき")
	}
	if c.ExpiresAt != nil {
		t.Error("本登録後はExpiresAtがnilであるべき")
	}
	if c.Word != "確定した言葉" {
		t.Errorf("Word=%v, want 確定した言葉", c.Word)
	}
}

func TestCard_Confirm_AlreadyConfirmed(t *testing.T) {
	c, _ := domain.NewCard(domain.NewCardInput{
		Number:    domain.MinCardNumber,
		Word:      "テスト",
		GuestName: strPtr("ゲスト"),
		Now:       time.Now(),
	})
	_ = c.Confirm("確定した言葉")

	err := c.Confirm("再度確定")
	if !errors.Is(err, domain.ErrAlreadyConfirmed) {
		t.Errorf("wantErr=%v, got err=%v", domain.ErrAlreadyConfirmed, err)
	}
}

func TestCard_Confirm_InvalidWord(t *testing.T) {
	tests := []struct {
		name    string
		word    string
		wantErr error
	}{
		{"最大文字数超過", strings.Repeat("あ", domain.MaxWordLength+1), domain.ErrInvalidWord},
		{"空文字", "", domain.ErrInvalidWord},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := domain.NewCard(domain.NewCardInput{
				Number:    domain.MinCardNumber,
				Word:      "テスト",
				GuestName: strPtr("ゲスト"),
				Now:       time.Now(),
			})
			err := c.Confirm(tt.word)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("wantErr=%v, got err=%v", tt.wantErr, err)
			}
		})
	}
}

func TestCard_VerifyOwner(t *testing.T) {
	profileID1 := int64(1)
	profileID2 := int64(2)

	tests := []struct {
		name      string
		c         *domain.Card
		profileID *int64
		guestName *string
		wantErr   error
	}{
		{"登録ユーザーのカード・本人", &domain.Card{ProfileID: &profileID1}, &profileID1, nil, nil},
		{"登録ユーザーのカード・他人", &domain.Card{ProfileID: &profileID1}, &profileID2, nil, domain.ErrForbidden},
		{"登録ユーザーのカード・未認証", &domain.Card{ProfileID: &profileID1}, nil, nil, domain.ErrForbidden},
		{"ゲストのカード・本人", &domain.Card{GuestName: strPtr("ゲスト")}, nil, strPtr("ゲスト"), nil},
		{"ゲストのカード・名前不一致", &domain.Card{GuestName: strPtr("ゲスト")}, nil, strPtr("別のゲスト"), domain.ErrForbidden},
		{"ゲストのカード・名前なし", &domain.Card{GuestName: strPtr("ゲスト")}, nil, nil, domain.ErrForbidden},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.c.VerifyOwner(tt.profileID, tt.guestName)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("wantErr=%v, got err=%v", tt.wantErr, err)
			}
		})
	}
}

func strPtr(s string) *string {
	return &s
}
