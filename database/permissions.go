package database

import (
	"gomet/globals"
)

func CheckPermission(userId int64, videoId int64) (bool, error) {

	stmt, err := db.Prepare(`SELECT count(*)
FROM membership
INNER JOIN video_permissions
ON video_permissions.groups_id = membership.groups_id
INNER JOIN video_classification
ON video_classification.video_groups_id = video_permissions.video_groups_id
AND video_classification.videos_id = ?
WHERE membership.users_id = ?;`)
	if err != nil {
		return false, err
	}

	var nb int64
	err = stmt.QueryRow(videoId, userId).Scan(&nb)
	if err != nil {
		return false, err
	}

	return nb > 0, nil
}

func IsAdmin(userId int64) (bool, error) {

	stmt, err := db.Prepare(`SELECT count(*)
FROM membership
WHERE users_id = ?
AND groups_id = ?;`)
	if err != nil {
		return false, err
	}

	var nb int64
	err = stmt.QueryRow(userId, globals.ADMIN_GROUP_ID).Scan(&nb)
	if err != nil {
		return false, err
	}

	return nb > 0, nil
}

func GetAllowedVideos(userId int64) ([]map[string]interface{}, error) {

	stmt, err := db.Preparex(`
SELECT videos.id, videos.title, videos.path, videos.slug
FROM videos
INNER JOIN video_classification
ON video_classification.videos_id = videos.id
INNER JOIN video_permissions
ON video_permissions.video_groups_id = video_classification.video_groups_id
INNER JOIN membership
ON membership.groups_id = video_permissions.groups_id
AND membership.users_id = ?;`)
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
