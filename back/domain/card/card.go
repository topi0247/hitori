package card

import "unicode/utf8"

const (
	MinCardNumber  = 1
	MaxCardNumber  = 100
	MaxWordLength  = 25
	MaxGuestNameLength = 10
)

// TODO: IsConfirmed, MatchPoints, ExpiresAt, CreatedAt, UpdatedAt を追加予定
type Card struct {
	Number    int
	Word      string
	ProfileID *int64
	GuestName *string
}

// NewInput は Card 生成時の入力値。Card とフィールドが異なるため分離している
type NewInput struct {
	Number    int
	Word      string
	ProfileID *int64
	GuestName *string
}

func New(input NewInput) (*Card, error) {
	if input.Number < MinCardNumber || input.Number > MaxCardNumber {
		return nil, ErrInvalidCardNumber
	}
	if utf8.RuneCountInString(input.Word) == 0 || utf8.RuneCountInString(input.Word) > MaxWordLength {
		return nil, ErrInvalidWord
	}
	if input.ProfileID == nil && input.GuestName == nil {
		return nil, ErrOwnerRequired
	}
	if input.GuestName != nil {
		if utf8.RuneCountInString(*input.GuestName) == 0 || utf8.RuneCountInString(*input.GuestName) > MaxGuestNameLength {
			return nil, ErrInvalidGuestName
		}
	}

	return &Card{
		Number:    input.Number,
		Word:      input.Word,
		ProfileID: input.ProfileID,
		GuestName: input.GuestName,
	}, nil
}
