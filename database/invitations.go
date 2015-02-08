package database

import (
	"errors"
	"github.com/julienc91/heygo/globals"
)

func getInvitationFromId(id int64) (globals.Invitation, error) {

	var query = "SELECT id, value FROM invitations WHERE id=?;"
	var params = []interface{}{id}
	res, err := getDb(query, params, func() interface{} { return globals.Invitation{} })
	if err != nil {
		return globals.Invitation{}, err
	} else if len(res) != 1 {
		return globals.Invitation{}, errors.New("Unexpected result from database")
	}
	return res[0].(globals.Invitation), nil
}

func getInvitationFromValue(value string) (globals.Invitation, error) {

	var query = "SELECT id, value FROM invitations WHERE value=?;"
	var params = []interface{}{value}
	res, err := getDb(query, params, func() interface{} { return globals.Invitation{} })
	if err != nil {
		return globals.Invitation{}, err
	} else if len(res) != 1 {
		return globals.Invitation{}, errors.New("Unexpected result from database")
	}
	return res[0].(globals.Invitation), nil
}

func getAllInvitations() ([]globals.Invitation, error) {

	var query = "SELECT id, value FROM invitations;"
	res, err := getDb(query, nil, func() interface{} { return globals.Invitation{} })
	if err != nil {
		return nil, err
	}

	var invitations []globals.Invitation
	for _, i := range res {
		invitations = append(invitations, i.(globals.Invitation))
	}
	return invitations, err
}

func updateInvitation(invitation globals.Invitation) (globals.Invitation, error) {

	var query = "UPDATE invitations SET value=? WHERE id=?;"
	var params = []interface{}{invitation.Value, invitation.Id}
	if err := updateDb(query, params); err != nil {
		return globals.Invitation{}, err
	}

	return getInvitationFromId(invitation.Id)
}

func insertInvitation(invitation globals.Invitation) (globals.Invitation, error) {

	var query = "INSERT INTO invitations (value) VALUES (?);"
	var params = []interface{}{invitation.Value}
	id, err := insertAndGetId(query, params)
	if err != nil {
		return globals.Invitation{}, err
	}

	return getInvitationFromId(id)
}

func deleteInvitationFromId(id int64) error {

	var query = "DELETE FROM invitations WHERE id=?;"
	var params = []interface{}{id}
	return deleteDb(query, params)
}

var GetInvitationFromId = getInvitationFromId
var GetInvitationFromValue = getInvitationFromValue
var GetAllInvitations = getAllInvitations
var UpdateInvitation = updateInvitation
var InsertInvitation = insertInvitation
var DeleteInvitationFromId = deleteInvitationFromId
