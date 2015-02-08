package database

import (
	"errors"
	"github.com/julienc91/heygo/tools"
)

// Check if the password is correct
func AuthenticateUser(id int64, password string) error {

	user, err := GetUserFromId(id)
	if err != nil {
		return errors.New("user does not exist")
	}

	if tools.Hash(password, user.Salt) != user.Password {
		return errors.New("invalid password")
	}

	return nil
}
