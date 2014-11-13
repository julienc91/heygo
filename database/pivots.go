package database


func DeletePivotTableRow(key1, key2 int, nameKey1, nameKey2 string, table string) error {

    stmt, err := db.Prepare("DELETE FROM " + table + " WHERE " + nameKey1 + " = ? AND " +
        nameKey2 + " = ?;")
    if err != nil {
        return err
    }

    _, err = stmt.Exec(key1, key2)
    return err

}



func GetAllPivotTable(table, key, value string) (map[int][]int, error) {

    var query = "SELECT " + key + ", " + value + " FROM " + table + ";"
    rows, err := db.Queryx(query)

    if err != nil {
        return nil, err
    }

    var res = make(map[int][]int)

    for rows.Next() {
        var key, value int
        err = rows.Scan(&key, &value)
        if err != nil {
            return nil, err
        }

        res[key] = append(res[key], value)
    }

    return res, nil
}


func InsertMembership(values map[string]interface{}) (map[string]interface{}, error) {

    return InsertRow(values, []string{"users_id", "groups_id"}, "membership")
}

func InsertClassification(values map[string]interface{}) (map[string]interface{}, error) {

    return InsertRow(values, []string{"videos_id", "video_groups_id"}, "video_classification")
}

func InsertPermission(values map[string]interface{}) (map[string]interface{}, error) {

    return InsertRow(values, []string{"video_groups_id", "groups_id"}, "video_permissions")
}
