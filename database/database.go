// This package handles all the interactions with the local database
// This file implements initialization and table-independant functions
package database

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"gomet/globals"
	"gomet/tools"
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
)

var db *sqlx.DB
var MainTables = []string{TableUsers, TableInvitations, TableGroups, TableVideoGroups, TableVideos}
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
		// videos (id, title, path, slug)
		`CREATE TABLE IF NOT EXISTS videos
(id INTEGER PRIMARY KEY AUTOINCREMENT,
 title VARCHAR UNIQUE NOT NULL DEFAULT '',
 path VARCHAR UNIQUE NOT NULL DEFAULT '',
 slug VARCHAR UNIQUE NOT NULL DEFAULT '');`,
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
 value VARCHAR UNIQUE NOT NULL);`}

	for _, query := range queries {
		db.MustExec(query)
	}
}

// Execute an update query on the given table, with the given parameters for the given id
// Values and table must have been checked first
func updateRow(id int64, values map[string]interface{}, validFields []string, table string) (map[string]interface{}, error) {

	var query = "UPDATE " + table + " SET "
	var valuesToSet []string
	var params []interface{}

	for k, v := range values {
		valuesToSet = append(valuesToSet, k+"=?")
		params = append(params, v)
	}

	query = query + strings.Join(valuesToSet, ",") + " WHERE id=?;"
	params = append(params, id)

	stmt, err := db.Preparex(query)
	if err != nil {
		return nil, err
	}

	_, err = stmt.Exec(params...)
	if err != nil {
		return nil, err
	}
	return getFromId(id, table)
}

// Execute an insert query on the given table, with the given parameters
// Values and table must have been checked first
func insertRow(values map[string]interface{}, validFields []string, table string) (map[string]interface{}, error) {

	var query = "INSERT INTO " + table
	var columnNames []string
	var columnValues []string
	var params []interface{}

	for k, v := range values {
		columnNames = append(columnNames, k)
		columnValues = append(columnValues, "?")
		params = append(params, v)
	}

	query = query + " (" + strings.Join(columnNames, ",") + ") VALUES (" + strings.Join(columnValues, ",") + ");"

	stmt, err := db.Preparex(query)
	if err != nil {
		return nil, err
	}

	res, err := stmt.Exec(params...)
	if err != nil {
		return nil, err
	}

	if tools.InArray(PivotTables, table) {
		return nil, nil
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	return getFromId(id, table)
}

// Execute a select query on the given table
// Table must have been checked first
func getAll(table string) ([]map[string]interface{}, error) {

	var query = "SELECT * FROM " + table + ";"
	stmt, err := db.Preparex(query)
	if err != nil {
		return nil, err
	}

	rows, err := stmt.Queryx()
	if err != nil {
		return nil, err
	}

	var res []map[string]interface{}

	for rows.Next() {
		var m = make(map[string]interface{})
		err = rows.MapScan(m)
		if err != nil {
			return nil, err
		}

		for k, v := range m {
			if vs, ok := v.([]byte); ok {
				m[k] = string(vs)
			}
		}
		res = append(res, m)
	}

	return res, nil
}

// Execute a select query on the given table with a where condition on the given column
// Table must have been checked first
func getColumnFiltered(column string, filter string, value interface{}, table string) ([]interface{}, error) {

	var query = "SELECT " + column + " FROM " + table + " WHERE " + filter + "=?;"
	stmt, err := db.Preparex(query)
	if err != nil {
		return nil, err
	}

	rows, err := stmt.Queryx(value)
	if err != nil {
		return nil, err
	}

	var res []interface{}

	for rows.Next() {
		var m interface{}
		err = rows.Scan(&m)
		if err != nil {
			return nil, err
		}

		if vs, ok := m.([]byte); ok {
			m = string(vs)
		}
		res = append(res, m)
	}

	return res, nil
}

// Execute a select query on the given table for the given id.
// An error is returned if there is no result.
// Table must have been checked first
func getFromId(id int64, table string) (map[string]interface{}, error) {

	return getFromKey("id", id, table)
}

// Execute a select query on the given table using the given unique key
// An error is returned if there is no result or more than one.
// Values and table must have been checked first
func getFromKey(key string, value interface{}, table string) (map[string]interface{}, error) {

	var query = "SELECT * FROM " + table + " WHERE " + key + "=?;"
	stmt, err := db.Preparex(query)
	if err != nil {
		return nil, err
	}

	row := stmt.QueryRowx(value)

	var res = make(map[string]interface{})
	if err = row.MapScan(res); err != nil {
		return nil, err
	}

	for k, v := range res {
		if vs, ok := v.([]byte); ok {
			res[k] = string(vs)
		}
	}

	return res, nil
}

// Execute a deletion query on the given table, with the given id
// Table must have been checked first
func deleteFromId(id int64, table string) error {

	return deleteFromKey("id", id, table)
}

// Execute a deletion query on the given table using the given unique key
// Values and table must have been checked first
func deleteFromKey(key string, value interface{}, table string) error {

	var query = "DELETE FROM " + table + " WHERE " + key + "=? LIMIT 1;"
	stmt, err := db.Preparex(query)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(value)
	return err
}

// Execute a deletion query on the given table using the given unique key
// Values and table must have been checked first
func deleteFromFilter(filter string, value interface{}, table string) error {

	var query = "DELETE FROM " + table + " WHERE " + filter + "=?;"
	stmt, err := db.Preparex(query)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(value)
	return err
}
