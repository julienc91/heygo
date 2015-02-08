package admin

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/julienc91/heygo/database"
	"github.com/julienc91/heygo/globals"
	"github.com/julienc91/heygo/tools"
	"github.com/julienc91/heygo/views"
	"net/http"
)

func deleteUser(w http.ResponseWriter, req *http.Request) {

	var user globals.User
	err := json.Unmarshal([]byte(req.FormValue("user")), &user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var ret = map[string]interface{}{"ok": true, "err": ""}
	err = database.DeleteUserFromId(user.Id)
	if err != nil {
		ret["err"] = err.Error()
		ret["ok"] = false
	}
	tools.WriteJsonResult(ret, w, http.StatusOK)
}

func deleteGroup(w http.ResponseWriter, req *http.Request) {

	var group globals.Group
	err := json.Unmarshal([]byte(req.FormValue("group")), &group)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var ret = map[string]interface{}{"ok": true, "err": ""}
	err = database.DeleteGroupFromId(group.Id)
	if err != nil {
		ret["err"] = err.Error()
		ret["ok"] = false
	}
	tools.WriteJsonResult(ret, w, http.StatusOK)
}

func deleteVideo(w http.ResponseWriter, req *http.Request) {

	var video globals.Video
	err := json.Unmarshal([]byte(req.FormValue("video")), &video)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var ret = map[string]interface{}{"ok": true, "err": ""}
	err = database.DeleteVideoFromId(video.Id)
	if err != nil {
		ret["err"] = err.Error()
		ret["ok"] = false
	}
	tools.WriteJsonResult(ret, w, http.StatusOK)
}

func deleteVideoGroup(w http.ResponseWriter, req *http.Request) {

	var videoGroup globals.VideoGroup
	err := json.Unmarshal([]byte(req.FormValue("video_group")), &videoGroup)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var ret = map[string]interface{}{"ok": true, "err": ""}
	err = database.DeleteVideoGroupFromId(videoGroup.Id)
	if err != nil {
		ret["err"] = err.Error()
		ret["ok"] = false
	}
	tools.WriteJsonResult(ret, w, http.StatusOK)
}

func deleteInvitation(w http.ResponseWriter, req *http.Request) {

	var invitation globals.Invitation
	err := json.Unmarshal([]byte(req.FormValue("invitation")), &invitation)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var ret = map[string]interface{}{"ok": true, "err": ""}
	err = database.DeleteInvitationFromId(invitation.Id)
	if err != nil {
		ret["err"] = err.Error()
		ret["ok"] = false
	}
	tools.WriteJsonResult(ret, w, http.StatusOK)
}

// Handle insert requests from admin panel
func AdminDelete(w http.ResponseWriter, req *http.Request) {

	if views.RedirectIfNotAdmin(w, req) {
		return
	}

	params := mux.Vars(req)

	table := params["table"]

	switch table {
	case database.TableUsers:
		deleteUser(w, req)
	case database.TableGroups:
		deleteGroup(w, req)
	case database.TableVideos:
		deleteVideo(w, req)
	case database.TableVideoGroups:
		deleteVideoGroup(w, req)
	case database.TableInvitations:
		deleteInvitation(w, req)
	default:
		http.Error(w, "table is not valid", http.StatusBadRequest)
	}
}
