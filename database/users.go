package database

import (
	"crypto/sha512"
	"encoding/base64"
	"errors"
	"fmt"
	"gomet/globals"
	"math/rand"
	"time"
)

type User struct {
	Id       int64
	Login    string
	Password string
	Salt     string
}

func AuthenticateUser(id int64, password string) error {

	stmt, err := db.Prepare(`SELECT password, salt
FROM users
WHERE id = ?;`)
	if err != nil {
		return err
	}

	var storedPassword, salt string
	err = stmt.QueryRow(id).Scan(&storedPassword, &salt)
	if err != nil {
		return errors.New("User does not exist")
	}

	if hashPassword(password, salt) != storedPassword {
		return errors.New("Invalid password")
	}

	return nil
}

func GetUserIdFromLogin(login string) (int64, error) {

	stmt, err := db.Prepare(`SELECT id
FROM users
WHERE login = ?;`)
	if err != nil {
		return 0, err
	}

	var id int64
	err = stmt.QueryRow(login).Scan(&id)
	if err != nil {
		return 0, errors.New("Unknown login")
	}

	return id, nil
}

func hashPassword(password, salt string) string {

	var hash = sha512.Sum512([]byte(salt + password))
	return base64.URLEncoding.EncodeToString(hash[:])
}

func saltGenerator() string {
	rand.Seed(time.Now().UnixNano())
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, globals.SALT_LENGTH)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func AddUser(login, password, invitation string) error {

	stmt, err := db.Prepare(`INSERT INTO users (login, password, salt)
VALUES (?, ?, ?);`)
	if err != nil {
		return err
	}

	var salt = saltGenerator()
	_, err = stmt.Exec(login, hashPassword(password, salt), salt)
	if err != nil {
		return err
	}
	return RevokeInvitation(invitation)
}

func CheckInvitation(invitation string) (bool, error) {

	stmt, err := db.Prepare(`SELECT count(*)
FROM invitations
WHERE value = ?;`)
	if err != nil {
		return false, err
	}

	var nb int
	err = stmt.QueryRow(invitation).Scan(&nb)
	if err != nil {
		return false, errors.New("Unknown invitation")
	}

	return nb > 0, nil
}

func RevokeInvitation(invitation string) error {

	stmt, err := db.Prepare(`DELETE FROM invitations
WHERE value = ?;`)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(invitation)
	return err
}

func GetAllUsers() ([]User, error) {

	stmt, err := db.Prepare(`SELECT id, login, password, salt FROM users;`)
	if err != nil {
		return nil, err
	}

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []User
	for rows.Next() {
		var user User
		err = rows.Scan(&user.Id, &user.Login, &user.Password, &user.Salt)
		if err != nil {
			return nil, err
		}

		res = append(res, user)
	}

	return res, nil
}

func UpdateUser(userId int64, values map[string]interface{}) (map[string]interface{}, error) {
	fmt.Println(values)
	if _, ok := values["password"]; ok {

		var salt = saltGenerator()
		var password = hashPassword(values["password"].(string), salt)

		values["password"] = password
		values["salt"] = salt
	}

	return UpdateRow(userId, values, []string{"login", "password", "salt"}, "users")
}

func InsertUser(values map[string]interface{}) (map[string]interface{}, error) {

	if _, ok := values["password"]; ok {

		var salt = saltGenerator()
		var password = hashPassword(values["password"].(string), salt)

		values["password"] = password
		values["salt"] = salt
	}

	return InsertRow(values, []string{"login", "password", "salt"}, "users")
}
