// This package handles all the interactions with the local database
// This file implements initialization and table-independant functions
package database

import (
	"errors"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/julienc91/heygo/globals"
	"strings"
)

const (
	TableUsers               = "users"
	TableInvitations         = "invitations"
	TableGroups              = "groups"
	TableVideoGroups         = "video_groups"
	TableVideos              = "videos"
	TableMembership          = "membership"
	TableVideoClassification = "video_classification"
	TableVideoPermissions    = "video_permissions"
	TableConfiguration       = "configuration"
)

var db *sqlx.DB
var MainTables = []string{TableUsers, TableInvitations, TableGroups, TableVideoGroups, TableVideos, TableConfiguration}
var PivotTables = []string{TableMembership, TableVideoClassification, TableVideoPermissions}
var Tables = append(MainTables, PivotTables...)

//This function opens the database and initializes it
func init() {
	var err error
	db, err = sqlx.Open("sqlite3", globals.DATABASE)
	if err != nil {
		panic(err)
	}
	InitDatabase()
}

// Initialize the database structure if needed
func InitDatabase() {

	var queries = []string{
		// users (id, login, password, salt)
		`CREATE TABLE IF NOT EXISTS users
(id INTEGER PRIMARY KEY AUTOINCREMENT,
 login VARCHAR UNIQUE NOT NULL DEFAULT '',
 password VARCHAR NOT NULL DEFAULT '',
 salt VARCHAR NOT NULL DEFAULT '');`,
		// groups (id, title)
		`CREATE TABLE IF NOT EXISTS groups
(id INTEGER PRIMARY KEY AUTOINCREMENT,
 title VARCHAR UNIQUE NOT NULL DEFAULT '');`,
		// membership (users_id, groups_id)
		`CREATE TABLE IF NOT EXISTS membership
(users_id INTEGER,
 groups_id INTEGER,
 FOREIGN KEY (users_id) REFERENCES users (id) ON DELETE CASCADE,
 FOREIGN KEY (groups_id) REFERENCES groups (id) ON DELETE CASCADE,
 PRIMARY KEY (users_id, groups_id));`,
		// videos (id, title, path, slug, imdb_id)
		`CREATE TABLE IF NOT EXISTS videos
(id INTEGER PRIMARY KEY AUTOINCREMENT,
 title VARCHAR UNIQUE NOT NULL DEFAULT '',
 path VARCHAR UNIQUE NOT NULL DEFAULT '',
 slug VARCHAR UNIQUE NOT NULL DEFAULT '',
 imdb_id VARCHAR);`,
		// video_groups (id, title)
		`CREATE TABLE IF NOT EXISTS video_groups
(id INTEGER PRIMARY KEY AUTOINCREMENT,
 title VARCHAR NOT NULL);`,
		// video_classification (videos_id, video_groups_id)
		`CREATE TABLE IF NOT EXISTS video_classification
(videos_id INTEGER,
 video_groups_id INTEGER,
 FOREIGN KEY (videos_id) REFERENCES videos (id) ON DELETE CASCADE,
 FOREIGN KEY (video_groups_id) REFERENCES video_groups (id) ON DELETE CASCADE,
 PRIMARY KEY (videos_id, video_groups_id));`,
		// video_permissions (video_groups_id, groups_id)
		`CREATE TABLE IF NOT EXISTS video_permissions
(video_groups_id INTEGER,
 groups_id INTEGER,
 FOREIGN KEY (video_groups_id) REFERENCES video_groups (id) ON DELETE CASCADE,
 FOREIGN KEY (groups_id) REFERENCES groups (id) ON DELETE CASCADE,
 PRIMARY KEY (video_groups_id, groups_id));`,
		// invitations (id, value)
		`CREATE TABLE IF NOT EXISTS invitations
(id INTEGER PRIMARY KEY AUTOINCREMENT,
 value VARCHAR UNIQUE NOT NULL);`,
		// configuration (id, key, value)
		`CREATE TABLE IF NOT EXISTS configuration
(id INTEGER PRIMARY KEY AUTOINCREMENT,
 key VARCHAR UNIQUE NOT NULL,
 value BLOB);`,
		// Default values
		`INSERT OR IGNORE INTO users (id, login, password, salt) VALUES
(1, 'admin', '8fFeeOTH2mSMU0Bb97LtDNYz6Nio2wYQgmQdl7cDKrXYTEoxyhy_sVe-d4oHM18KLguW1ppj-_gs_oAyYYEVcQ==', '9PiJbmrrfo3urA4');`,
		`INSERT OR IGNORE INTO groups (id, title) VALUES
(1, 'admin');`,
		`INSERT OR IGNORE INTO membership (users_id, groups_id) VALUES
(1, 1);`,
		`INSERT OR IGNORE INTO configuration (key, value) VALUES
('domain', 'localhost'), ('port', '8080'),
('opensubtitles_login', ''), ('opensubtitles_password', ''), ('opensubtitles_useragent', 'OSTestUserAgent');`}

	for _, query := range queries {
		db.MustExec(query)
	}
}

func _genericQueryExec(query string, params []interface{}) error {

	stmt, err := db.Preparex(query)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(params...)
	if err != nil {
		return err
	}
	return nil
}

var insertDb = _genericQueryExec
var updateDb = _genericQueryExec
var deleteDb = _genericQueryExec

func insertAndGetId(query string, params []interface{}) (int64, error) {

	stmt, err := db.Preparex(query)
	if err != nil {
		return 0, err
	}

	res, err := stmt.Exec(params...)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func getDb(query string, params []interface{}, constructor func() interface{}) ([]interface{}, error) {

	stmt, err := db.Preparex(query)
	if err != nil {
		return nil, err
	}

	rows, err := stmt.Queryx()
	if err != nil {
		return nil, err
	}

	var res []interface{}

	for rows.Next() {
		var obj = constructor()
		err = rows.StructScan(&obj)
		if err != nil {
			return nil, err
		}
		res = append(res, obj)
	}

	return res, nil
}

/*************************************************************** INVITATIONS */
func getInvitationFromId(id int64) (globals.Invitation, error) {

	var query = "SELECT id, value FROM invitations WHERE id=?;"
	var params = []interface{}{id}
	res, err := getDb(query, params, func() interface{} {return globals.Invitation{}})
	if err != nil {
		return globals.Invitation{}, err
	} else if len(res) != 1 {
		return globals.Invitation{}, errors.New("Unexpected result from database")
	}
	return res[0].(globals.Invitation), nil
}

func getInvitationFromValue(value string) (globals.Invitation, error) {

	var query = "SELECT id, value FROM invitations WHERE value=?;"
	var params = []interface{}{value}
	res, err := getDb(query, params, func() interface{} {return globals.Invitation{}})
	if err != nil {
		return globals.Invitation{}, err
	} else if len(res) != 1 {
		return globals.Invitation{}, errors.New("Unexpected result from database")
	}
	return res[0].(globals.Invitation), nil
}

func getAllInvitations() ([]globals.Invitation, error) {

	var query = "SELECT id, value FROM invitations;"
	res, err := getDb(query, nil, func() interface{} {return globals.Invitation{}})
	if err != nil {
		return nil, err
	}

	var invitations []globals.Invitation
	for _, i := range res {
		invitations = append(invitations, i.(globals.Invitation))
	}
	return invitations, err
}

func updateInvitation(invitation globals.Invitation) (globals.Invitation, error) {

	var query = "UPDATE invitations SET value=? WHERE id=?;"
	var params = []interface{}{invitation.Value, invitation.Id}
	if err := updateDb(query, params); err != nil {
		return globals.Invitation{}, err
	}

	return getInvitationFromId(invitation.Id)
}

func insertInvitation(invitation globals.Invitation) (globals.Invitation, error) {

	var query = "INSERT INTO invitations (value) VALUES (?);"
	var params = []interface{}{invitation.Value}
	id, err := insertAndGetId(query, params);
	if err != nil {
		return globals.Invitation{}, err
	}

	return getInvitationFromId(id)
}

func deleteInvitationFromId(id int64) error {

	var query = "DELETE FROM invitations WHERE id=?;"
	var params = []interface{}{id}
	return deleteDb(query, params)
}

var GetInvitationFromId = getInvitationFromId
var GetInvitationFromValue = getInvitationFromValue
var GetAllInvitations = getAllInvitations
var UpdateInvitation = updateInvitation
var InsertInvitation = insertInvitation
var DeleteInvitationFromId = deleteInvitationFromId

/********************************************************** USERS AND GROUPS */
func getGroupsFromUserId(id int64) ([]globals.Group, error) {

	var query = "SELECT groups.id, groups.title FROM groups "
	query += "INNER JOIN membership ON membership.users_id=? AND groups.id = membership.groups_id;"
	var params = []interface{}{id}
	res, err := getDb(query, params, func() interface{} {return globals.Group{}})
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
	res, err := getDb(query, params, func() interface{} { return globals.Group{}})
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
	res, err := getDb(query, params, func() interface{} {return globals.User{}})
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
	res, err := getDb(query, params, func() interface{} {return globals.User{}})
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
	res, err := getDb(query, params, func() interface{} {return globals.User{}})
	if err != nil {
		return globals.User{}, err
	} else if len(res) != 1 {
		return globals.User{}, errors.New("Unexpected result from database")
	}
	return res[0].(globals.User), nil
}

func getAllUsers() ([]globals.User, error) {

	var query = "SELECT id, login, password, salt FROM users;"
	res, err := getDb(query, nil, func() interface{} {return globals.User{}})
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
	res, err := getDb(query, nil, func() interface{} {return globals.Group{}})
	if err != nil {
		return nil, err
	}

	var groups []globals.Group
	for _, g := range res {
		groups = append(groups, g.(globals.Group))
	}
	return groups, err
}

func updateMembershipFromUserId(id int64, groups []globals.Group) error {

	if err := deleteMembershipFromUserId(id); err != nil {
		return err
	}

	var query = "INSERT INTO membership (users_id, groups_id) VALUES "
	var params = []interface{}{}
	var values []string
	for _, group := range groups {
		values = append(values, "(?, ?)")
		params = append(params, id, group.Id)
	}
	query += strings.Join(values, ", ") + ";"

	return insertDb(query, params)
}

func updateMembershipFromGroupId(id int64, users []globals.User) error {

	if err := deleteMembershipFromGroupId(id); err != nil {
		return err
	}

	var query = "INSERT INTO membership (users_id, groups_id) VALUES "
	var params = []interface{}{}
	var values []string
	for _, user := range users {
		values = append(values, "(?, ?)")
		params = append(params, user.Id, id)
	}
	query += strings.Join(values, ", ") + ";"

	return insertDb(query, params)
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
	id, err := insertAndGetId(query, params);
	if err != nil {
		return globals.Group{}, err
	}

	return getGroupFromId(id)
}

func insertUser(user globals.User) (globals.User, error) {

	var query = "INSERT INTO videos (login, password, salt) VALUES (?, ?, ?, ?);"
	var params = []interface{}{user.Login, user.Password, user.Salt}
	id, err := insertAndGetId(query, params);
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

var GetUserFromLogin = getUserFromLogin
var GetAllUsers = getAllUsers
var GetAllGroups = getAllGroups
var GetGroupFromId = getGroupFromId
var GetGroupsFromUserId = getGroupsFromUserId
var GetUsersFromGroupId = getUsersFromGroupId
var InsertUser = insertUser
var UpdateUser = updateUser
var UpdateGroup = updateGroup
var UpdateMembershipFromGroupId = updateMembershipFromGroupId
var UpdateMembershipFromUserId = updateMembershipFromUserId
var DeleteUserFromId = deleteUserFromId
var DeleteGroupFromId = deleteGroupFromId

/*************************************************** VIDEOS AND VIDEO GROUPS */
func getVideoGroupsFromVideoId(id int64) ([]globals.VideoGroup, error) {

	var query = "SELECT video_groups.id, video_groups.title FROM video_groups "
	query += "INNER JOIN video_classification ON video_classification.videos_id=? AND video_groups.id = video_classification.video_groups_id;"
	var params = []interface{}{id}
	res, err := getDb(query, params, func() interface{} {return globals.VideoGroup{}})
	if err != nil {
		return nil, err
	}
	var videoGroups []globals.VideoGroup
	for _, v := range res {
		videoGroups = append(videoGroups, v.(globals.VideoGroup))
	}
	return videoGroups, err
}

func getVideoGroupFromId(id int64) (globals.VideoGroup, error) {

	var query = "SELECT id, title FROM video_groups WHERE id=?;"
	var params = []interface{}{id}
	res, err := getDb(query, params, func() interface{} { return globals.VideoGroup{}})
	if err != nil {
		return globals.VideoGroup{}, err
	} else if len(res) != 1 {
		return globals.VideoGroup{}, errors.New("Unexpected result from database")
	}
	return res[0].(globals.VideoGroup), nil
}

func getVideosFromVideoGroupId(id int64) ([]globals.Video, error) {

	var query = "SELECT videos.id, videos.title, videos.slug, videos.path, videos.imdb_id FROM videos "
	query += "INNER JOIN video_classification ON video_classification.videos_id = videos.id AND video_classification.video_groups_id=?;"
	var params = []interface{}{id}
	res, err := getDb(query, params, func() interface{} {return globals.Video{}})
	if err != nil {
		return nil, err
	}
	var videos []globals.Video
	for _, v := range res {
		videos = append(videos, v.(globals.Video))
	}
	return videos, nil
}

func getVideoFromId(id int64) (globals.Video, error) {

	var query = "SELECT id, title, slug, path, imdb_id FROM videos WHERE id=?;"
	var params = []interface{}{id}
	res, err := getDb(query, params, func() interface{} {return globals.Video{}})
	if err != nil {
		return globals.Video{}, err
	} else if len(res) != 1 {
		return globals.Video{}, errors.New("Unexpected result from database")
	}
	return res[0].(globals.Video), nil
}

func getVideoFromSlug(slug string) (globals.Video, error) {

	var query = "SELECT id, title, slug, path, imdb_id FROM videos WHERE slug=?;"
	var params = []interface{}{slug}
	res, err := getDb(query, params, func() interface{} {return globals.Video{}})
	if err != nil {
		return globals.Video{}, err
	} else if len(res) != 1 {
		return globals.Video{}, errors.New("Unexpected result from database")
	}
	return res[0].(globals.Video), nil
}

func getAllVideos() ([]globals.Video, error) {

	var query = "SELECT id, title, slug, path, imdb_id FROM videos;"
	res, err := getDb(query, nil, func() interface{} {return globals.Video{}})
	if err != nil {
		return nil, err
	}

	var videos []globals.Video
	for _, v := range res {
		videos = append(videos, v.(globals.Video))
	}
	return videos, err
}

func getAllVideoGroups() ([]globals.VideoGroup, error) {

	var query = "SELECT id, title FROM video_groups;"
	res, err := getDb(query, nil, func() interface{} {return globals.VideoGroup{}})
	if err != nil {
		return nil, err
	}

	var videoGroups []globals.VideoGroup
	for _, v := range res {
		videoGroups = append(videoGroups, v.(globals.VideoGroup))
	}
	return videoGroups, err
}

func updateVideoClassificationFromVideoId(id int64, videoGroups []globals.VideoGroup) error {

	if err := deleteVideoClassificationFromVideoId(id); err != nil {
		return err
	}

	var query = "INSERT INTO video_classification (videos_id, video_groups_id) VALUES "
	var params = []interface{}{}
	var values []string
	for _, videoGroup := range videoGroups {
		values = append(values, "(?, ?)")
		params = append(params, id, videoGroup.Id)
	}
	query += strings.Join(values, ", ") + ";"

	return insertDb(query, params)
}

func updateVideoClassificationFromVideoGroupId(id int64, videos []globals.Video) error {

	if err := deleteVideoClassificationFromVideoGroupId(id); err != nil {
		return err
	}

	var query = "INSERT INTO video_classification (videos_id, video_groups_id) VALUES "
	var params = []interface{}{}
	var values []string
	for _, video := range videos {
		values = append(values, "(?, ?)")
		params = append(params, video.Id, id)
	}
	query += strings.Join(values, ", ") + ";"

	return insertDb(query, params)
}

func updateVideoGroup(videoGroup globals.VideoGroup) (globals.VideoGroup, error) {

	var query = "UPDATE video_groups SET title=? WHERE id=?;"
	var params = []interface{}{videoGroup.Title, videoGroup.Id}
	if err := updateDb(query, params); err != nil {
		return globals.VideoGroup{}, err
	}

	return getVideoGroupFromId(videoGroup.Id)
}

func updateVideo(video globals.Video) (globals.Video, error) {

	var query = "UPDATE videos SET title=?, path=?, imdb_id=?, slug=? WHERE id=?;"
	var params = []interface{}{video.Title, video.Path, video.ImdbId, video.Slug, video.Id}
	if err := updateDb(query, params); err != nil {
		return globals.Video{}, err
	}

	return getVideoFromId(video.Id)
}

func insertVideoGroup(videoGroup globals.VideoGroup) (globals.VideoGroup, error) {

	var query = "INSERT INTO video_groups (title) VALUES (?);"
	var params = []interface{}{videoGroup.Title}
	id, err := insertAndGetId(query, params);
	if err != nil {
		return globals.VideoGroup{}, err
	}

	return getVideoGroupFromId(id)
}

func insertVideo(video globals.Video) (globals.Video, error) {

	var query = "INSERT INTO videos (title, path, slug, imdb_id) VALUES (?, ?, ?, ?);"
	var params = []interface{}{video.Title, video.Path, video.Slug, video.ImdbId}
	id, err := insertAndGetId(query, params);
	if err != nil {
		return globals.Video{}, err
	}

	return getVideoFromId(id)
}

func deleteVideoClassificationFromVideoId(id int64) error {

	var query = "DELETE FROM video_classification WHERE videos_id=?;"
	var params = []interface{}{id}
	return deleteDb(query, params)
}

func deleteVideoClassificationFromVideoGroupId(id int64) error {

	var query = "DELETE FROM video_classification WHERE video_groups_id=?;"
	var params = []interface{}{id}
	return deleteDb(query, params)
}

func deleteVideoFromId(id int64) error {

	if err := deleteVideoClassificationFromVideoId(id); err != nil {
		return err
	}

	if err := deletePermissionsFromGroupId(id); err != nil {
		return err
	}

	var query = "DELETE FROM videos WHERE id=?;"
	var params = []interface{}{id}
	return deleteDb(query, params)
}

func deleteVideoGroupFromId(id int64) error {

	if err := deleteVideoClassificationFromVideoGroupId(id); err != nil {
		return err
	}

	if err := deletePermissionsFromVideoGroupId(id); err != nil {
		return err
	}

	var query = "DELETE FROM video_groups WHERE id=?;"
	var params = []interface{}{id}
	return deleteDb(query, params)
}

var GetVideoFromSlug = getVideoFromSlug
var GetAllVideos = getAllVideos
var GetAllVideoGroups = getAllVideoGroups
var GetVideoGroupFromId = getVideoGroupFromId
var GetVideoGroupsFromVideoId = getVideoGroupsFromVideoId
var GetVideosFromVideoGroupId = getVideosFromVideoGroupId
var InsertVideo = insertVideo
var UpdateVideo = updateVideo
var UpdateVideoGroup = updateVideoGroup
var UpdateVideoClassificationFromVideoGroupId = updateVideoClassificationFromVideoGroupId
var UpdateVideoClassificationFromVideoId = updateVideoClassificationFromVideoId
var DeleteVideoFromId = deleteVideoFromId
var DeleteVideoGroupFromId = deleteVideoGroupFromId

/*************************************************************** PERMISSIONS */
func getPermissionsFromGroupId(id int64) ([]globals.VideoGroup, error) {

	var query = "SELECT video_groups.id, video_groups.title FROM video_groups "
	query += "INNER JOIN video_permissions ON video_permissions.groups_id=? AND video_permissions.video_groups_id=video_groups.id;"
	var params = []interface{}{id}

	res, err := getDb(query, params, func () interface{} {return globals.VideoGroup{}})
	if err != nil {
		return nil, err
	}
	var videoGroups []globals.VideoGroup
	for _, v := range res {
		videoGroups = append(videoGroups, v.(globals.VideoGroup))
	}
	return videoGroups, nil
}

func getPermissionsFromVideoGroupId(id int64) ([]globals.Group, error) {

	var query = "SELECT groups.id, groups.title FROM groups "
	query += "INNER JOIN video_permissions ON video_permissions.video_groups_id=? AND video_permissions.groups_id=groups.id;"
	var params = []interface{}{id}

	res, err := getDb(query, params, func () interface{} {return globals.Group{}})
	if err != nil {
		return nil, err
	}
	var groups []globals.Group
	for _, g := range res {
		groups = append(groups, g.(globals.Group))
	}
	return groups, nil
}

func updatePermissionsFromGroupId(id int64, videoGroups []globals.VideoGroup) error {

	if err := deletePermissionsFromGroupId(id); err != nil {
		return err
	}

	var query = "INSERT INTO video_permissions (groups_id, video_groups_id) VALUES "
	var params = []interface{}{}
	var values []string
	for _, videoGroup := range videoGroups {
		values = append(values, "(?, ?)")
		params = append(params, id, videoGroup.Id)
	}
	query += strings.Join(values, ", ") + ";"

	return insertDb(query, params)
}

func updatePermissionsFromVideoGroupId(id int64, groups []globals.Group) error {

	if err := deletePermissionsFromVideoGroupId(id); err != nil {
		return err
	}

	var query = "INSERT INTO video_permissions (groups_id, video_groups_id) VALUES "
	var params = []interface{}{}
	var values []string
	for _, group := range groups {
		values = append(values, "(?, ?)")
		params = append(params, group.Id, id)
	}
	query += strings.Join(values, ", ") + ";"

	return insertDb(query, params)
}

func deletePermissionsFromGroupId(id int64) error {

	var query = "DELETE FROM video_permissions WHERE groups_id=?;"
	var params = []interface{}{id}
	return deleteDb(query, params)
}

func deletePermissionsFromVideoGroupId(id int64) error {

	var query = "DELETE FROM video_permissions WHERE video_groups_id=?;"
	var params = []interface{}{id}
	return deleteDb(query, params)
}

var GetPermissionsFromGroupId = getPermissionsFromGroupId
var GetPermissionsFromVideoGroupId = getPermissionsFromVideoGroupId
var UpdatePermissionsFromGroupId = updatePermissionsFromGroupId
var UpdatePermissionsFromVideoGroupId = updatePermissionsFromVideoGroupId

/************************************************************* CONFIGURATION */
func LoadConfiguration() error {

	var query = "SELECT key, value FROM configuration;"
	stmt, err := db.Preparex(query)
	if err != nil {
		return err
	}

	rows, err := stmt.Queryx()
	if err != nil {
		return err
	}

	var res []interface{}

	for rows.Next() {
		var m = make(map[string]interface{})
		err = rows.MapScan(m)
		if err != nil {
			return err
		}
		switch m["key"].(string) {
			case "domain":
				globals.CONFIGURATION.Domain = m["value"].(string)
			case "port":
				globals.CONFIGURATION.Port = m["value"].(string)
			case "opensubtitles_login":
				globals.CONFIGURATION.OpensubtitlesLogin = m["value"].(string)
			case "opensubtitles_password":
				globals.CONFIGURATION.OpensubtitlesPassword = m["value"].(string)
			case "opensubtitles_useragent":
				globals.CONFIGURATION.OpensubtitlesUseragent = m["value"].(string)
		}
	}
	return nil
}

func SetConfiguration(key string, value interface{}) error {

	var query = "UPDATE configuration SET value=? WHERE key=?;"
	var params = []interface{}{key, value}
	return insertDb(query, params)
}
