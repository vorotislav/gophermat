package models

import (
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	// stringShortLen defines length limit for string fields with short values.
	stringShortLen = 64

	// passwordMinLen defines minimum length limit for password fields.
	passwordMinLen = 5
)

type User struct {
	ID       int    `json:"id"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (u *User) Validate() error {
	err := validation.ValidateStruct(u,
		validation.Field(&u.Login,
			validation.Required),
		validation.Field(&u.Password,
			validation.Required,
			validation.Length(passwordMinLen, stringShortLen)),
	)

	if err != nil {
		return fmt.Errorf("validation error: %w", err)
	}

	return nil
}
