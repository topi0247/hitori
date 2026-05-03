package card_test

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/topi0247/hitori/domain/card"
)

func TestNewCard_CardNumber(t *testing.T) {
	tests := []struct {
		name    string
		number  int
		wantErr error
	}{
		{"最小値", card.MinCardNumber, nil},
		{"最大値", card.MaxCardNumber, nil},
		{"最小値未満", card.MinCardNumber - 1, card.ErrInvalidCardNumber},
		{"最大値超過", card.MaxCardNumber + 1, card.ErrInvalidCardNumber},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := card.New(card.NewInput{
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
		{"最大文字数", strings.Repeat("あ", card.MaxWordLength), nil},
		{"最大文字数超過", strings.Repeat("あ", card.MaxWordLength+1), card.ErrInvalidWord},
		{"空文字", "", card.ErrInvalidWord},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := card.New(card.NewInput{
				Number:    card.MinCardNumber,
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
		{"最大文字数", strPtr(strings.Repeat("あ", card.MaxGuestNameLength)), nil},
		{"最大文字数超過", strPtr(strings.Repeat("あ", card.MaxGuestNameLength+1)), card.ErrInvalidGuestName},
		{"空文字", strPtr(""), card.ErrInvalidGuestName},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := card.New(card.NewInput{
				Number:    card.MinCardNumber,
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
		{"両方nil", nil, nil, card.ErrOwnerRequired},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := card.New(card.NewInput{
				Number:    card.MinCardNumber,
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
	c, err := card.New(card.NewInput{
		Number:    card.MinCardNumber,
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
	want := now.Add(card.ProvisionalTTL)
	if !c.ExpiresAt.Equal(want) {
		t.Errorf("ExpiresAt=%v, want %v", c.ExpiresAt, want)
	}
}

func TestCard_Confirm(t *testing.T) {
	c, _ := card.New(card.NewInput{
		Number:    card.MinCardNumber,
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
	c, _ := card.New(card.NewInput{
		Number:    card.MinCardNumber,
		Word:      "テスト",
		GuestName: strPtr("ゲスト"),
		Now:       time.Now(),
	})
	_ = c.Confirm("確定した言葉")

	err := c.Confirm("再度確定")
	if !errors.Is(err, card.ErrAlreadyConfirmed) {
		t.Errorf("wantErr=%v, got err=%v", card.ErrAlreadyConfirmed, err)
	}
}

func TestCard_Confirm_InvalidWord(t *testing.T) {
	tests := []struct {
		name    string
		word    string
		wantErr error
	}{
		{"最大文字数超過", strings.Repeat("あ", card.MaxWordLength+1), card.ErrInvalidWord},
		{"空文字", "", card.ErrInvalidWord},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := card.New(card.NewInput{
				Number:    card.MinCardNumber,
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

func strPtr(s string) *string {
	return &s
}
