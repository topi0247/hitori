package theme

import "unicode/utf8"

const MaxTitleLength = 100

type Theme struct {
	ID    int64
	Title string
}

func New(title string) (*Theme, error) {
	if utf8.RuneCountInString(title) == 0 || utf8.RuneCountInString(title) > MaxTitleLength {
		return nil, ErrInvalidTitle
	}
	return &Theme{Title: title}, nil
}
