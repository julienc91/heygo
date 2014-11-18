package database

import (
	"heygo/globals"
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
