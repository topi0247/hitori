package profile_test

import (
	"strings"
	"testing"

	"github.com/topi0247/hitori/domain/profile"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name     string
		userName string
		wantErr  error
	}{
		{"正常", "テスト", nil},
		{"最大文字数", strings.Repeat("あ", profile.MaxUserNameLength), nil},
		{"最大文字数超過", strings.Repeat("あ", profile.MaxUserNameLength+1), profile.ErrInvalidUserName},
		{"空文字", "", profile.ErrInvalidUserName},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, err := profile.New("auth-user-id", tt.userName)
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
		{"最大文字数", strings.Repeat("あ", profile.MaxUserNameLength), nil},
		{"最大文字数超過", strings.Repeat("あ", profile.MaxUserNameLength+1), profile.ErrInvalidUserName},
		{"空文字", "", profile.ErrInvalidUserName},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &profile.Profile{UserName: "既存"}
			err := p.UpdateUserName(tt.input)
			if err != tt.wantErr {
				t.Errorf("wantErr=%v, got=%v", tt.wantErr, err)
			}
		})
	}
}
