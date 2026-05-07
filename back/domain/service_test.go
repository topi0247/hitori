package domain_test

import (
	"errors"
	"testing"

	"github.com/topi0247/hitori/domain"
)

func TestCardCountPolicy_CanAdd(t *testing.T) {
	tests := []struct {
		name    string
		count   int
		wantErr error
	}{
		{"0枚", 0, nil},
		{"上限未満", domain.MaxCardsPerTheme - 1, nil},
		{"上限ちょうど", domain.MaxCardsPerTheme, domain.ErrThemeCardLimitReached},
		{"上限超過", domain.MaxCardsPerTheme + 1, domain.ErrThemeCardLimitReached},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := domain.CanAddCard(tt.count)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("wantErr=%v, got err=%v", tt.wantErr, err)
			}
		})
	}
}
