package database

type Group struct {
    Id int
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


func UpdateGroup(groupId int, values map[string]interface{}) error {

    return UpdateRow(groupId, values, []string{"title"}, "groups")
}

func InsertGroup(values map[string]interface{}) (int, error) {

    return InsertRow(values, []string{"title"}, "groups")
}
