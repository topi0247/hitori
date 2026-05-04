package card

import "errors"

var (
	ErrInvalidCardNumber     = errors.New("invalid_card_number")
	ErrInvalidWord           = errors.New("invalid_word")
	ErrInvalidGuestName      = errors.New("invalid_guest_name")
	ErrOwnerRequired         = errors.New("owner_required")
	ErrAlreadyConfirmed      = errors.New("already_confirmed")
	ErrThemeCardLimitReached = errors.New("theme_card_limit_reached")
	ErrNotFound              = errors.New("not_found")
	ErrForbidden             = errors.New("forbidden")
)
