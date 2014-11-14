package database

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"gomet/globals"
)

var db *sqlx.DB
var MainTables = []string{"users", "invitations", "groups", "video_groups", "videos"}
var PivotTables = []string{"membership", "video_classification", "video_permissions"}
var Tables = append(MainTables, PivotTables...)

func init() {
	var err error
	db, err = sqlx.Open("sqlite3", globals.DATABASE)
	if err != nil {
		panic(err)
	}
	InitDatabase()
}

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

func checkFields(values map[string]interface{}, validFields []string) error {

	for field := range values {
		var isFieldValid = false
		for _, validField := range validFields {
			if validField == field {
				isFieldValid = true
				break
			}
		}
		if !isFieldValid {
			delete(values, field)
		}
	}

	return nil
}

func UpdateRow(id int64, values map[string]interface{}, validFields []string, table string) (map[string]interface{}, error) {

	if err := checkFields(values, validFields); err != nil {
		return nil, err
	}

	var query = "UPDATE " + table + " SET "
	var params []interface{}

	for k, v := range values {
		query += k + "=?,"
		params = append(params, v)
	}

	query = query[:len(query)-1] + " WHERE id=?;"
	params = append(params, id)

	stmt, err := db.Prepare(query)
	if err != nil {
		return nil, err
	}

	_, err = stmt.Exec(params...)
	if err != nil {
		return nil, err
	}
	return GetFromId(table, id)
}

func InsertRow(values map[string]interface{}, validFields []string, table string) (map[string]interface{}, error) {

	if err := checkFields(values, validFields); err != nil {
		return nil, err
	}

	var query = "INSERT INTO " + table + " ("
	var query2 = ") VALUES ("
	var params []interface{}

	for k, v := range values {
		query += k + ","
		query2 += "?,"
		params = append(params, v)
	}

	query = query[:len(query)-1] + query2[:len(query2)-1] + ");"

	stmt, err := db.Prepare(query)
	if err != nil {
		return nil, err
	}

	res, err := stmt.Exec(params...)
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	return GetFromId(table, id)
}

func DeleteRow(id int64, table string) error {

	stmt, err := db.Prepare("DELETE FROM " + table + " WHERE id=?;")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(id)
	return err
}

func GetAll(table string) ([]map[string]interface{}, error) {

	var query = "SELECT * FROM " + table + ";"
	rows, err := db.Queryx(query)

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

func GetFromId(table string, id int64) (map[string]interface{}, error) {

	var query = "SELECT * FROM " + table + " WHERE id=?;"
	row := db.QueryRowx(query, id)

	var res = make(map[string]interface{})
	err := row.MapScan(res)
	if err != nil {
		return nil, err
	}

	for k, v := range res {
		if vs, ok := v.([]byte); ok {
			res[k] = string(vs)
		}
	}

	return res, nil
}
