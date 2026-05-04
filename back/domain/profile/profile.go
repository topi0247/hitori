package profile

import "unicode/utf8"

const MaxUserNameLength = 10

type Profile struct {
	ID         int64
	AuthUserID string
	UserName   string
}

func (p *Profile) UpdateUserName(name string) error {
	if utf8.RuneCountInString(name) == 0 || utf8.RuneCountInString(name) > MaxUserNameLength {
		return ErrInvalidUserName
	}
	p.UserName = name
	return nil
}
