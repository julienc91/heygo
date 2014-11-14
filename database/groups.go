package database

type Group struct {
	Id    int64
	Title string
}

func GetAllGroups() ([]Group, error) {

	stmt, err := db.Prepare(`SELECT id, title FROM groups;`)
	if err != nil {
		return nil, err
	}

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []Group

	for rows.Next() {
		var group Group
		if err := rows.Scan(&group.Id, &group.Title); err != nil {
			return nil, err
		}

		result = append(result, group)
	}

	return result, nil
}

func UpdateGroup(groupId int64, values map[string]interface{}) (map[string]interface{}, error) {

	return UpdateRow(groupId, values, []string{"title"}, "groups")
}

func InsertGroup(values map[string]interface{}) (map[string]interface{}, error) {

	return InsertRow(values, []string{"title"}, "groups")
}
