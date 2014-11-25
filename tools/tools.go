// Some useful functions
package tools

import (
	crand "crypto/rand"
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"heygo/globals"
	"io"
	mrand "math/rand"
	"os"
	"time"
)

// Check if e is in a
func InArray(a []string, e string) bool {
	for _, x := range a {
		if x == e {
			return true
		}
	}
	return false
}

// Check if a file exists
func CheckFilePath(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// Hash function
func Hash(password, salt string) string {
	var hash = sha512.Sum512([]byte(salt + password))
	return base64.URLEncoding.EncodeToString(hash[:])
}

// Salt generating function
func SaltGenerator() string {

	mrand.Seed(time.Now().UnixNano())
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, globals.SALT_LENGTH)
	for i := range b {
		b[i] = letters[mrand.Intn(len(letters))]
	}
	return string(b)
}

// UUID generating function
// From: http://play.golang.org/p/4FkNSiUDMg
func UuidGenerator() (string, error) {
	uuid := make([]byte, 16)
	n, err := io.ReadFull(crand.Reader, uuid)
	if n != len(uuid) || err != nil {
		return "", err
	}
	uuid[8] = uuid[8]&^0xc0 | 0x80
	uuid[6] = uuid[6]&^0xf0 | 0x40
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]), nil
}
