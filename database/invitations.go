package database

type Invitation struct {
    Id int
    Value string
}

func InsertInvitation(values map[string]interface{}) (int, error) {

    return InsertRow(values, []string{"value"}, "invitations")
}

func UpdateInvitation(invitationId int, values map[string]interface{}) error {

    return UpdateRow(invitationId, values, []string{"value"}, "invitations")
}
