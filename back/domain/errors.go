package domain

import "errors"

var (
	ErrInvalidCardNumber     = errors.New("invalid_card_number")
	ErrInvalidWord           = errors.New("invalid_word")
	ErrInvalidGuestName      = errors.New("invalid_guest_name")
	ErrOwnerRequired         = errors.New("owner_required")
	ErrAlreadyConfirmed      = errors.New("already_confirmed")
	ErrThemeCardLimitReached = errors.New("theme_card_limit_reached")
	ErrForbidden             = errors.New("forbidden")
	ErrInvalidUserName       = errors.New("invalid_user_name")
	ErrAlreadyExists         = errors.New("already_exists")
	ErrInvalidTitle          = errors.New("invalid_title")
	ErrNotFound              = errors.New("not_found")
)
