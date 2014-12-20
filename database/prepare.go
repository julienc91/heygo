package database

import (
	"errors"
	"heygo/globals"
	"heygo/tools"
)

var validColumns = map[string][]string{
	TableUsers:               []string{"login", "password"},
	TableInvitations:         []string{"value"},
	TableVideos:              []string{"path", "slug", "title", "imdb_id"},
	TableGroups:              []string{"title"},
	TableVideoGroups:         []string{"title"},
	TableMembership:          []string{"users_id", "groups_id"},
	TableVideoClassification: []string{"videos_id", "video_groups_id"},
	TableVideoPermissions:    []string{"video_groups_id", "groups_id"},
	TableConfiguration:       []string{"key", "value"}}
var uniqueColumns = map[string][]string{
	TableUsers:               []string{"login"},
	TableInvitations:         []string{"value"},
	TableVideos:              []string{"slug"},
	TableGroups:              []string{"title"},
	TableVideoGroups:         []string{"title"},
	TableMembership:          nil,
	TableVideoClassification: nil,
	TableVideoPermissions:    nil,
	TableConfiguration:       []string{"key"}}
var requiredColumns = map[string][]string{
	TableUsers:               []string{"login", "password"},
	TableInvitations:         []string{"value"},
	TableVideos:              []string{"path", "slug", "title"},
	TableGroups:              []string{"title"},
	TableVideoGroups:         []string{"title"},
	TableMembership:          []string{"users_id", "groups_id"},
	TableVideoClassification: []string{"videos_id", "video_groups_id"},
	TableVideoPermissions:    []string{"video_groups_id", "groups_id"},
	TableConfiguration:       []string{"key", "value"}}
var pivots = map[string][][2]string{
	TableUsers:               [][2]string{{TableMembership, "users_id"}},
	TableInvitations:         nil,
	TableVideos:              [][2]string{{TableVideoClassification, "videos_id"}},
	TableGroups:              [][2]string{{TableMembership, "groups_id"}, {TableVideoPermissions, "groups_id"}},
	TableVideoGroups:         [][2]string{{TableVideoClassification, "video_groups_id"}, {TableVideoPermissions, "video_groups_id"}},
	TableMembership:          nil,
	TableVideoClassification: nil,
	TableVideoPermissions:    nil,
	TableConfiguration:       nil}

// Remove keys from the map which are not in validFields
func checkFields(values map[string]interface{}, validFields []string, requiredFields []string) error {

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
	for _, requiredField := range requiredFields {
		var isFieldSet = false
		for field := range values {
			if requiredField == field {
				isFieldSet = true
				break
			}
		}
		if !isFieldSet {
			return errors.New("missing required field")
		}
	}
	if len(values) == 0 {
		return errors.New("no values")
	}

	return nil
}

// Prepare and check arguments before calling insertRow
func prepareInsert(values map[string]interface{}, table string) (map[string]interface{}, error) {

	columns, ok := validColumns[table]
	if !ok {
		return nil, errors.New("invalid table")
	}
	if err := checkFields(values, columns, requiredColumns[table]); err != nil {
		return nil, err
	}

	if table == TableUsers {
		var salt = tools.SaltGenerator()
		var password = tools.Hash(values["password"].(string), salt)

		values["password"] = password
		values["salt"] = salt
	}

	return insertRow(values, table)
}

// Prepare and check arguments before calling updateFromId
func prepareUpdateFromId(id int64, values map[string]interface{}, table string) (map[string]interface{}, error) {

	columns, ok := validColumns[table]
	if !ok {
		return nil, errors.New("invalid table")
	}
	if err := checkFields(values, columns, nil); err != nil {
		return nil, err
	}

	if _, ok := values["password"]; ok {

		var salt = tools.SaltGenerator()
		var password = tools.Hash(values["password"].(string), salt)

		values["password"] = password
		values["salt"] = salt
	}

	return updateFromId(id, values, table)
}

// Prepare and check arguments before calling updateFromKey
func prepareUpdateConfiguration(key string, value interface{}) (map[string]interface{}, error) {

	row, err := prepareGetFromKey("key", key, TableConfiguration)
	if err != nil {
		return nil, err
	}
	res, err := prepareUpdateFromId(row["id"].(int64), map[string]interface{}{"value": value}, TableConfiguration)
	globals.LoadConfiguration <- true
	return res, err
}

// Prepare and check arguments before calling getAll
func prepareGetAll(table string) ([]map[string]interface{}, error) {

	if !tools.InArray(Tables, table) {
		return nil, errors.New("invalid table")
	}
	return getAll(table)
}

// Prepare and check arguments before calling getAllFilteredByColumn
func prepareGetColumnFiltered(column string, filter string, value interface{}, table string) ([]interface{}, error) {

	columns, ok := validColumns[table]
	if !ok {
		return nil, errors.New("invalid table")
	}
	if !tools.InArray(columns, column) {
		return nil, errors.New("invalid column")
	}
	if !tools.InArray(columns, filter) {
		return nil, errors.New("invalid filter")
	}
	return getColumnFiltered(column, filter, value, table)
}

// Prepare and check arguments before calling getFromId
func prepareGetFromId(id int64, table string) (map[string]interface{}, error) {

	if !tools.InArray(MainTables, table) {
		return nil, errors.New("invalid table")
	}

	return getFromId(id, table)
}

// Prepare and check arguments before calling getFromKey
func prepareGetFromKey(key string, value interface{}, table string) (map[string]interface{}, error) {

	columns, ok := uniqueColumns[table]
	if !ok {
		return nil, errors.New("invalid table")
	}
	if !tools.InArray(columns, key) {
		return nil, errors.New("invalid key")
	}
	return getFromKey(key, value, table)
}

// Prepare and check arguments before calling deleteFromId
func prepareDeleteFromId(id int64, table string) error {

	if !tools.InArray(MainTables, table) {
		return errors.New("invalid table")
	}
	if err := deleteDependencies(id, table); err != nil {
		return err
	}

	return deleteFromId(id, table)
}

// Prepare and check arguments before calling deleteFromKey
func prepareDeleteFromKey(key string, value interface{}, table string) error {

	columns, ok := uniqueColumns[table]
	if !ok {
		return errors.New("invalid table")
	}
	if !tools.InArray(columns, key) {
		return errors.New("invalid key")
	}
	if err := deleteDependencies(value, table); err != nil {
		return err
	}
	return deleteFromKey(key, value, table)
}

// Prepare and check arguments before calling deleteFromFilter
func prepareDeleteFromFilter(filter string, value interface{}, table string) error {

	columns, ok := validColumns[table]
	if !ok {
		return errors.New("invalid table")
	}
	if !tools.InArray(columns, filter) {
		return errors.New("invalid filter")
	}
	if err := deleteDependencies(value, table); err != nil {
		return err
	}
	return deleteFromFilter(filter, value, table)
}
