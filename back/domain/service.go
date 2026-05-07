package domain

const MaxCardsPerTheme = 100

func CanAddCard(currentCount int) error {
	if currentCount >= MaxCardsPerTheme {
		return ErrThemeCardLimitReached
	}
	return nil
}
