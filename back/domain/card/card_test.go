package card_test

import (
	"errors"
	"strings"
	"testing"

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
			})
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("wantErr=%v, got err=%v", tt.wantErr, err)
			}
		})
	}
}

func strPtr(s string) *string {
	return &s
}
