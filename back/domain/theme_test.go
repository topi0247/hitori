package domain_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/topi0247/hitori/domain"
)

func TestNewTheme_Title(t *testing.T) {
	tests := []struct {
		name    string
		title   string
		wantErr error
	}{
		{"通常", "大きさ", nil},
		{"最大文字数", strings.Repeat("あ", domain.MaxTitleLength), nil},
		{"最大文字数超過", strings.Repeat("あ", domain.MaxTitleLength+1), domain.ErrInvalidTitle},
		{"空文字", "", domain.ErrInvalidTitle},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := domain.NewTheme(tt.title)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("wantErr=%v, got err=%v", tt.wantErr, err)
			}
		})
	}
}
