package database

import (
	"errors"
	"fmt"
	"github.com/julienc91/heygo/globals"
	"github.com/julienc91/heygo/tools"
	"strings"
)

// Check if the password is correct
func AuthenticateUser(id int64, password string) error {

	user, _, err := GetUserFromId(id)
	if err != nil {
		return err
	}

	if tools.Hash(password, user.Salt) != user.Password {
		return errors.New("invalid password")
	}

	return nil
}

func getGroupsFromUserId(id int64) ([]globals.Group, error) {

	var query = "SELECT groups.id, groups.title FROM groups "
	query += "INNER JOIN membership ON membership.users_id=? AND groups.id = membership.groups_id;"
	var params = []interface{}{id}
	res, err := getDb(query, params, func() interface{} { return globals.Group{} })
	if err != nil {
		return nil, err
	}
	var groups []globals.Group
	for _, g := range res {
		groups = append(groups, g.(globals.Group))
	}
	return groups, err
}

func getGroupFromId(id int64) (globals.Group, error) {

	var query = "SELECT id, title FROM groups WHERE id=?;"
	var params = []interface{}{id}
	res, err := getDb(query, params, func() interface{} { return globals.Group{} })
	if err != nil {
		return globals.Group{}, err
	} else if len(res) != 1 {
		return globals.Group{}, errors.New("Unexpected result from database")
	}
	return res[0].(globals.Group), nil
}

func getUsersFromGroupId(id int64) ([]globals.User, error) {

	var query = "SELECT users.id, users.login, users.password, users.salt FROM users "
	query += "INNER JOIN membership ON membership.users_id=users.id AND membership.groups_id=?;"
	var params = []interface{}{id}
	res, err := getDb(query, params, func() interface{} { return globals.User{} })
	if err != nil {
		return nil, err
	}
	var users []globals.User
	for _, u := range res {
		users = append(users, u.(globals.User))
	}
	return users, nil
}

func getUserFromId(id int64) (globals.User, error) {

	var query = "SELECT id, login, password, salt FROM users WHERE id=?;"
	var params = []interface{}{id}
	res, err := getDb(query, params, func() interface{} { return globals.User{} })
	if err != nil {
		return globals.User{}, err
	} else if len(res) != 1 {
		return globals.User{}, errors.New("Unexpected result from database")
	}
	return res[0].(globals.User), nil
}

func getUserFromLogin(login string) (globals.User, error) {

	var query = "SELECT id, login, password, salt FROM users WHERE login=?;"
	var params = []interface{}{login}
	res, err := getDb(query, params, func() interface{} { return globals.User{} })
	if err != nil {
		return globals.User{}, err
	} else if len(res) != 1 {
		fmt.Println(res)
		return globals.User{}, errors.New("Unexpected result from database")
	}
	return res[0].(globals.User), nil
}

func getAllUsers() ([]globals.User, error) {

	var query = "SELECT id, login, password, salt FROM users;"
	res, err := getDb(query, nil, func() interface{} { return globals.User{} })
	if err != nil {
		return nil, err
	}

	var users []globals.User
	for _, u := range res {
		users = append(users, u.(globals.User))
	}
	return users, err
}

func getAllGroups() ([]globals.Group, error) {

	var query = "SELECT id, title FROM groups;"
	res, err := getDb(query, nil, func() interface{} { return globals.Group{} })
	if err != nil {
		return nil, err
	}

	var groups []globals.Group
	for _, g := range res {
		groups = append(groups, g.(globals.Group))
	}
	return groups, err
}

func updateMembershipFromUserId(id int64, groups []globals.Group) ([]globals.Group, error) {

	if err := deleteMembershipFromUserId(id); err != nil {
		return nil, err
	}

	var query = "INSERT INTO membership (users_id, groups_id) VALUES "
	var params = []interface{}{}
	var values []string
	for _, group := range groups {
		values = append(values, "(?, ?)")
		params = append(params, id, group.Id)
	}
	query += strings.Join(values, ", ") + ";"

	if err := insertDb(query, params); err != nil {
		return nil, err
	}
	return getGroupsFromUserId(id)
}

func updateMembershipFromGroupId(id int64, users []globals.User) ([]globals.User, error) {

	if err := deleteMembershipFromGroupId(id); err != nil {
		return nil, err
	}

	var query = "INSERT INTO membership (users_id, groups_id) VALUES "
	var params = []interface{}{}
	var values []string
	for _, user := range users {
		values = append(values, "(?, ?)")
		params = append(params, user.Id, id)
	}
	query += strings.Join(values, ", ") + ";"

	if err := insertDb(query, params); err != nil {
		return nil, err
	}
	return getUsersFromGroupId(id)
}

func updateGroup(group globals.Group) (globals.Group, error) {

	var query = "UPDATE video_groups SET title=? WHERE id=?;"
	var params = []interface{}{group.Title, group.Id}
	if err := updateDb(query, params); err != nil {
		return globals.Group{}, err
	}

	return getGroupFromId(group.Id)
}

func updateUser(user globals.User) (globals.User, error) {

	var query = "UPDATE users SET login=?, password=?, salt=? WHERE id=?;"
	var params = []interface{}{user.Login, user.Password, user.Salt, user.Id}
	if err := updateDb(query, params); err != nil {
		return globals.User{}, err
	}

	return getUserFromId(user.Id)
}

func insertGroup(group globals.Group) (globals.Group, error) {

	var query = "INSERT INTO groups (title) VALUES (?);"
	var params = []interface{}{group.Title}
	id, err := insertAndGetId(query, params)
	if err != nil {
		return globals.Group{}, err
	}

	return getGroupFromId(id)
}

func insertUser(user globals.User) (globals.User, error) {

	var query = "INSERT INTO videos (login, password, salt) VALUES (?, ?, ?, ?);"
	var params = []interface{}{user.Login, user.Password, user.Salt}
	id, err := insertAndGetId(query, params)
	if err != nil {
		return globals.User{}, err
	}

	return getUserFromId(id)
}

func deleteMembershipFromUserId(id int64) error {

	var query = "DELETE FROM membership WHERE users_id=?;"
	var params = []interface{}{id}
	return deleteDb(query, params)
}

func deleteMembershipFromGroupId(id int64) error {

	var query = "DELETE FROM membership WHERE groups_id=?;"
	var params = []interface{}{id}
	return deleteDb(query, params)
}

func deleteUserFromId(id int64) error {

	if err := deleteMembershipFromUserId(id); err != nil {
		return err
	}
	var query = "DELETE FROM users WHERE id=?;"
	var params = []interface{}{id}
	return deleteDb(query, params)
}

func deleteGroupFromId(id int64) error {

	if err := deleteMembershipFromGroupId(id); err != nil {
		return err
	}

	var query = "DELETE FROM grous WHERE id=?;"
	var params = []interface{}{id}
	return deleteDb(query, params)
}

// Public functions
func UpdateGroup(group globals.Group, users []globals.User, videoGroups []globals.VideoGroup) (globals.Group, []globals.User, []globals.VideoGroup, error) {

	newGroup, err := updateGroup(group)
	if err != nil {
		return globals.Group{}, nil, nil, err
	}

	newUsers, err := updateMembershipFromGroupId(newGroup.Id, users)
	if err != nil {
		return globals.Group{}, nil, nil, err
	}

	newVideoGroups, err := updatePermissionsFromGroupId(newGroup.Id, videoGroups)
	if err != nil {
		return globals.Group{}, nil, nil, err
	}
	return newGroup, newUsers, newVideoGroups, nil
}

func UpdateUser(user globals.User, groups []globals.Group) (globals.User, []globals.Group, error) {

	newUser, err := updateUser(user)
	if err != nil {
		return globals.User{}, nil, err
	}

	newGroups, err := updateMembershipFromUserId(newUser.Id, groups)
	if err != nil {
		return globals.User{}, nil, err
	}
	return newUser, newGroups, nil
}

func InsertGroup(group globals.Group, users []globals.User, videoGroups []globals.VideoGroup) (globals.Group, []globals.User, []globals.VideoGroup, error) {

	newGroup, err := insertGroup(group)
	if err != nil {
		return globals.Group{}, nil, nil, err
	}

	newUsers, err := updateMembershipFromGroupId(newGroup.Id, users)
	if err != nil {
		return globals.Group{}, nil, nil, err
	}

	newVideoGroups, err := updatePermissionsFromGroupId(newGroup.Id, videoGroups)
	if err != nil {
		return globals.Group{}, nil, nil, err
	}
	return newGroup, newUsers, newVideoGroups, nil
}

func InsertUser(user globals.User, groups []globals.Group) (globals.User, []globals.Group, error) {

	newUser, err := insertUser(user)
	if err != nil {
		return globals.User{}, nil, err
	}

	newGroups, err := updateMembershipFromUserId(newUser.Id, groups)
	if err != nil {
		return globals.User{}, nil, err
	}
	return newUser, newGroups, nil
}

func GetUserFromId(id int64) (globals.User, []globals.Group, error) {

	user, err := getUserFromId(id)
	if err != nil {
		return globals.User{}, nil, err
	}
	groups, err := getGroupsFromUserId(id)
	if err != nil {
		return globals.User{}, nil, err
	}
	return user, groups, nil
}

func GetUserFromLogin(login string) (globals.User, []globals.Group, error) {

	user, err := getUserFromLogin(login)
	if err != nil {
		return globals.User{}, nil, err
	}
	groups, err := getGroupsFromUserId(user.Id)
	if err != nil {
		return globals.User{}, nil, err
	}
	return user, groups, nil
}

func GetGroupFromId(id int64) (globals.Group, []globals.User, []globals.VideoGroup, error) {

	group, err := getGroupFromId(id)
	if err != nil {
		return globals.Group{}, nil, nil, err
	}
	users, err := getUsersFromGroupId(id)
	if err != nil {
		return globals.Group{}, nil, nil, err
	}
	videoGroups, err := getPermissionsFromGroupId(id)
	if err != nil {
		return globals.Group{}, nil, nil, err
	}
	return group, users, videoGroups, nil
}

var GetAllUsers = getAllUsers
var GetAllGroups = getAllGroups
var DeleteUserFromId = deleteUserFromId
var DeleteGroupFromId = deleteGroupFromId
