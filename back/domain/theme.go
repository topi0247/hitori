package domain

//go:generate mockgen -source=theme.go -destination=mock/theme_mock.go -package=domainmock

import (
	"context"
	"unicode/utf8"
)

const MaxTitleLength = 100

type Theme struct {
	ID    int64
	Title string
}

func NewTheme(title string) (*Theme, error) {
	if utf8.RuneCountInString(title) == 0 || utf8.RuneCountInString(title) > MaxTitleLength {
		return nil, ErrInvalidTitle
	}
	return &Theme{Title: title}, nil
}

type ThemeRepository interface {
	FetchAll(ctx context.Context) ([]*Theme, error)
	FetchByID(ctx context.Context, id int64) (*Theme, error)
}
