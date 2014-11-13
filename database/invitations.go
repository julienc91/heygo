package database

type Invitation struct {
    Id int64
    Value string
}

func InsertInvitation(values map[string]interface{}) (map[string]interface{}, error) {

    return InsertRow(values, []string{"value"}, "invitations")
}

func UpdateInvitation(invitationId int64, values map[string]interface{}) (map[string]interface{}, error) {

    return UpdateRow(invitationId, values, []string{"value"}, "invitations")
}
