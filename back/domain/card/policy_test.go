package card_test

import (
	"errors"
	"testing"

	"github.com/topi0247/hitori/domain/card"
)

func TestCardCountPolicy_CanAdd(t *testing.T) {
	tests := []struct {
		name    string
		count   int
		wantErr error
	}{
		{"0枚", 0, nil},
		{"上限未満", card.MaxCardsPerTheme - 1, nil},
		{"上限ちょうど", card.MaxCardsPerTheme, card.ErrThemeCardLimitReached},
		{"上限超過", card.MaxCardsPerTheme + 1, card.ErrThemeCardLimitReached},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := card.CanAddCard(tt.count)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("wantErr=%v, got err=%v", tt.wantErr, err)
			}
		})
	}
}
