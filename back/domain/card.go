package domain

import (
	"time"
	"unicode/utf8"
)

const (
	MinCardNumber      = 1
	MaxCardNumber      = 100
	MaxWordLength      = 25
	MaxGuestNameLength = 10
	ProvisionalTTL     = 24 * time.Hour
)

// TODO: CreatedAt, UpdatedAt を追加予定
// NewInput は Card 生成時の入力値。Card とフィールドが異なるため分離している
type Card struct {
	ID          int64
	UUID        string
	ThemeID     int64
	Number      int
	Word        string
	ProfileID   *int64
	GuestName   *string
	IsConfirmed bool
	MatchPoints int
	ExpiresAt   *time.Time
}

type NewCardInput struct {
	ThemeID   int64
	Number    int
	Word      string
	ProfileID *int64
	GuestName *string
	Now       time.Time
}

func NewCard(input NewCardInput) (*Card, error) {
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

	expiresAt := input.Now.Add(ProvisionalTTL)
	return &Card{
		ThemeID:     input.ThemeID,
		Number:      input.Number,
		Word:        input.Word,
		ProfileID:   input.ProfileID,
		GuestName:   input.GuestName,
		IsConfirmed: false,
		MatchPoints: 0,
		ExpiresAt:   &expiresAt,
	}, nil
}

func (c *Card) VerifyOwner(profileID *int64, guestName *string) error {
	if c.ProfileID != nil {
		if profileID == nil || *c.ProfileID != *profileID {
			return ErrForbidden
		}
		return nil
	}
	if c.GuestName != nil {
		if guestName == nil || *c.GuestName != *guestName {
			return ErrForbidden
		}
		return nil
	}
	return ErrForbidden
}

func (c *Card) Confirm(word string) error {
	if c.IsConfirmed {
		return ErrAlreadyConfirmed
	}
	if utf8.RuneCountInString(word) == 0 || utf8.RuneCountInString(word) > MaxWordLength {
		return ErrInvalidWord
	}
	c.Word = word
	c.IsConfirmed = true
	c.ExpiresAt = nil
	return nil
}
