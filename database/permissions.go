package database

import (
	"github.com/julienc91/heygo/globals"
	"strings"
)

// Check if the given user is allowed to access the given video
func CheckPermission(userId int64, videoId int64) (bool, error) {

	var query = "SELECT count(*) FROM membership "
	query += "INNER JOIN video_permissions ON video_permissions.groups_id = membership.groups_id "
	query += "INNER JOIN video_classification ON video_classification.video_groups_id = video_permissions.video_groups_id AND video_classification.videos_id=? "
	query += "WHERE membership.users_id=?;"

	stmt, err := db.Prepare(query)
	if err != nil {
		return false, err
	}

	var nb int64
	if err := stmt.QueryRow(videoId, userId).Scan(&nb); err != nil {
		return false, err
	}

	return nb > 0, nil
}

// Check if the given user as admin rights
func IsAdmin(userId int64) (bool, error) {

	var query = "SELECT count(*) FROM membership WHERE users_id=? AND groups_id=?;"

	stmt, err := db.Preparex(query)
	if err != nil {
		return false, err
	}

	var nb int64
	if err := stmt.QueryRow(userId, globals.ADMIN_GROUP_ID).Scan(&nb); err != nil {
		return false, err
	}

	return nb > 0, nil
}

// Get all videos allowed to the given user
func GetAllowedVideos(userId int64) ([]map[string]interface{}, error) {

	var query = "SELECT videos.* FROM videos "
	query += "INNER JOIN video_classification ON video_classification.videos_id = videos.id "
	query += "INNER JOIN video_permissions ON video_permissions.video_groups_id = video_classification.video_groups_id "
	query += "INNER JOIN membership ON membership.groups_id = video_permissions.groups_id AND membership.users_id=?;"

	stmt, err := db.Preparex(query)
	if err != nil {
		return nil, err
	}

	rows, err := stmt.Queryx(userId)
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

func getPermissionsFromGroupId(id int64) ([]globals.VideoGroup, error) {

	var query = "SELECT video_groups.id, video_groups.title FROM video_groups "
	query += "INNER JOIN video_permissions ON video_permissions.groups_id=? AND video_permissions.video_groups_id=video_groups.id;"
	var params = []interface{}{id}

	res, err := getDb(query, params, func() interface{} { return globals.VideoGroup{} })
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

	res, err := getDb(query, params, func() interface{} { return globals.Group{} })
	if err != nil {
		return nil, err
	}
	var groups []globals.Group
	for _, g := range res {
		groups = append(groups, g.(globals.Group))
	}
	return groups, nil
}

func updatePermissionsFromGroupId(id int64, videoGroups []globals.VideoGroup) ([]globals.VideoGroup, error) {

	if err := deletePermissionsFromGroupId(id); err != nil {
		return nil, err
	}

	var query = "INSERT INTO video_permissions (groups_id, video_groups_id) VALUES "
	var params = []interface{}{}
	var values []string
	for _, videoGroup := range videoGroups {
		values = append(values, "(?, ?)")
		params = append(params, id, videoGroup.Id)
	}
	query += strings.Join(values, ", ") + ";"

	if err := insertDb(query, params); err != nil {
		return nil, err
	}
	return getPermissionsFromGroupId(id)
}

func updatePermissionsFromVideoGroupId(id int64, groups []globals.Group) ([]globals.Group, error) {

	if err := deletePermissionsFromVideoGroupId(id); err != nil {
		return nil, err
	}

	var query = "INSERT INTO video_permissions (groups_id, video_groups_id) VALUES "
	var params = []interface{}{}
	var values []string
	for _, group := range groups {
		values = append(values, "(?, ?)")
		params = append(params, group.Id, id)
	}
	query += strings.Join(values, ", ") + ";"

	if err := insertDb(query, params); err != nil {
		return nil, err
	}
	return getPermissionsFromVideoGroupId(id)
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

// Public functions
func UpdatePermissionsFromGroup(group globals.Group, videoGroups []globals.VideoGroup) ([]globals.VideoGroup, error) {
	return updatePermissionsFromGroupId(group.Id, videoGroups)
}

func UpdatePermissionsFromVideoGroup(videoGroup globals.VideoGroup, groups []globals.Group) ([]globals.Group, error) {
	return updatePermissionsFromVideoGroupId(videoGroup.Id, groups)
}

var GetPermissionsFromGroupId = getPermissionsFromGroupId
var GetPermissionsFromVideoGroupId = getPermissionsFromVideoGroupId
