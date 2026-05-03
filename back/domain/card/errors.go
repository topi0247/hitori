package card

import "errors"

var (
	ErrInvalidCardNumber = errors.New("invalid_card_number")
	ErrInvalidWord       = errors.New("invalid_word")
	ErrInvalidGuestName  = errors.New("invalid_guest_name")
	ErrOwnerRequired     = errors.New("owner_required")
)
