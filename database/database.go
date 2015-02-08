// This package handles all the interactions with the local database
// This file implements initialization and table-independant functions
package database

import (
	"errors"
	"github.com/jmoiron/sqlx"
	"github.com/julienc91/heygo/globals"
	_ "github.com/mattn/go-sqlite3"
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
 imdb_id VARCHAR NOT NULL DEFAULT '');`,
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

	appendToRes := func(obj interface{}, res []interface{}, err error) []interface{} {
		if err == nil {
			return append(res, obj)
		}
		return res
	}

	stmt, err := db.Preparex(query)
	if err != nil {
		return nil, err
	}

	rows, err := stmt.Queryx(params...)
	if err != nil {
		return nil, err
	}

	var res []interface{}

	for rows.Next() {
		var obj = constructor()

		switch obj.(type) {
		case globals.Video:
			var tmp globals.Video
			err = rows.StructScan(&tmp)
			res = appendToRes(tmp, res, err)
		case globals.User:
			var tmp globals.User
			err = rows.StructScan(&tmp)
			res = appendToRes(tmp, res, err)
		case globals.Group:
			var tmp globals.Group
			err = rows.StructScan(&tmp)
			res = appendToRes(tmp, res, err)
		case globals.VideoGroup:
			var tmp globals.VideoGroup
			err = rows.StructScan(&tmp)
			res = appendToRes(tmp, res, err)
		case globals.Invitation:
			var tmp globals.Invitation
			err = rows.StructScan(&tmp)
			res = appendToRes(tmp, res, err)
		default:
			return nil, errors.New("unhandled type of constructor")
		}
		if err != nil {
			return nil, err
		}
	}

	return res, nil
}
