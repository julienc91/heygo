package database

import (
	"errors"
	"heygo/tools"
)

// Check if the password is correct
func AuthenticateUser(id int64, password string) error {

	user, err := PrepareGetFromId(id, TableUsers)
	if err != nil {
		return errors.New("user does not exist")
	}

	if tools.Hash(password, user["salt"].(string)) != user["password"].(string) {
		return errors.New("invalid password")
	}

	return nil
}
