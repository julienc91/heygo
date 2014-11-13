package database

type VideoGroup struct {
    Id int64
    Title string
}


func GetAllVideoGroups() ([]VideoGroup, error) {

    stmt, err := db.Prepare(`SELECT id, title FROM video_groups;`)
    if err != nil {
        return nil, err
    }

    rows, err := stmt.Query()
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var result []VideoGroup

    for rows.Next() {
        var group VideoGroup
        if err := rows.Scan(&group.Id, &group.Title); err != nil {
            return nil, err
        }

        result = append(result, group)
    }
    
    return result, nil
}


func UpdateVideoGroup(groupId int64, values map[string]interface{}) (map[string]interface{}, error) {

    return UpdateRow(groupId, values, []string{"title"}, "video_groups")
}

func InsertVideoGroup(values map[string]interface{}) (map[string]interface{}, error) {

    return InsertRow(values, []string{"title"}, "video_groups")
}
