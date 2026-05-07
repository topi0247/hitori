package domain_test

import (
	"strings"
	"testing"

	"github.com/topi0247/hitori/domain"
)

func TestNewProfile(t *testing.T) {
	tests := []struct {
		name     string
		userName string
		wantErr  error
	}{
		{"正常", "テスト", nil},
		{"最大文字数", strings.Repeat("あ", domain.MaxUserNameLength), nil},
		{"最大文字数超過", strings.Repeat("あ", domain.MaxUserNameLength+1), domain.ErrInvalidUserName},
		{"空文字", "", domain.ErrInvalidUserName},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, err := domain.NewProfile("auth-user-id", tt.userName)
			if err != tt.wantErr {
				t.Errorf("wantErr=%v, got=%v", tt.wantErr, err)
			}
			if err == nil && p.UserName != tt.userName {
				t.Errorf("UserName=%v, got=%v", tt.userName, p.UserName)
			}
		})
	}
}

func TestUpdateUserName(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr error
	}{
		{"最大文字数", strings.Repeat("あ", domain.MaxUserNameLength), nil},
		{"最大文字数超過", strings.Repeat("あ", domain.MaxUserNameLength+1), domain.ErrInvalidUserName},
		{"空文字", "", domain.ErrInvalidUserName},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &domain.Profile{UserName: "既存"}
			err := p.UpdateUserName(tt.input)
			if err != tt.wantErr {
				t.Errorf("wantErr=%v, got=%v", tt.wantErr, err)
			}
		})
	}
}
