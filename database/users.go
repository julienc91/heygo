package database

import (
	"crypto/sha512"
	"encoding/base64"
	"errors"
	"gomet/globals"
	"math/rand"
	"time"
)

// Check if the password is correct
func AuthenticateUser(id int64, password string) error {

	user, err := PrepareGetFromId(id, TableUsers)
	if err != nil {
		return errors.New("user does not exist")
	}

	if hashPassword(password, user["salt"].(string)) != user["password"].(string) {
		return errors.New("invalid password")
	}

	return nil
}

// Hashing function
func hashPassword(password, salt string) string {

	var hash = sha512.Sum512([]byte(salt + password))
	return base64.URLEncoding.EncodeToString(hash[:])
}

// Salt generating function
func saltGenerator() string {

	rand.Seed(time.Now().UnixNano())
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, globals.SALT_LENGTH)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
