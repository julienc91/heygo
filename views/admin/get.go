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

func getUser(w http.ResponseWriter, req *http.Request) {

	var user globals.User
	err := json.Unmarshal([]byte(req.FormValue("user")), &user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, groups, err := database.GetUserFromId(user.Id)
	var ret = map[string]interface{}{"ok": true, "err": ""}
	if err != nil {
		ret["err"] = err.Error()
		ret["ok"] = false
	}

	ret["data"] = map[string]interface{}{"user": user, "groups": groups}
	tools.WriteJsonResult(ret, w, http.StatusOK)
}

func getGroup(w http.ResponseWriter, req *http.Request) {

	var group globals.Group
	err := json.Unmarshal([]byte(req.FormValue("group")), &group)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var ret = map[string]interface{}{"ok": true, "err": ""}
	group, users, videoGroups, err := database.GetGroupFromId(group.Id)
	if err != nil {
		ret["err"] = err.Error()
		ret["ok"] = false
	}
	ret["data"] = map[string]interface{}{"group": group, "users": users, "video_groups": videoGroups}
	tools.WriteJsonResult(ret, w, http.StatusOK)
}

func getVideo(w http.ResponseWriter, req *http.Request) {

	var video globals.Video
	err := json.Unmarshal([]byte(req.FormValue("video")), &video)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var ret = map[string]interface{}{"ok": true, "err": ""}
	video, videoGroups, err := database.GetVideoFromId(video.Id)
	if err != nil {
		ret["err"] = err.Error()
		ret["ok"] = false
	}
	ret["data"] = map[string]interface{}{"video": video, "video_groups": videoGroups}
	tools.WriteJsonResult(ret, w, http.StatusOK)
}

func getVideoGroup(w http.ResponseWriter, req *http.Request) {

	var videoGroup globals.VideoGroup
	err := json.Unmarshal([]byte(req.FormValue("video_group")), &videoGroup)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var ret = map[string]interface{}{"ok": true, "err": ""}
	videoGroup, videos, groups, err := database.GetVideoGroupFromId(videoGroup.Id)
	if err != nil {
		ret["err"] = err.Error()
		ret["ok"] = false
	}
	ret["data"] = map[string]interface{}{"video_group": videoGroup, "videos": videos, "groups": groups}
	tools.WriteJsonResult(ret, w, http.StatusOK)
}

func getInvitation(w http.ResponseWriter, req *http.Request) {

	var invitation globals.Invitation
	err := json.Unmarshal([]byte(req.FormValue("invitation")), &invitation)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var ret = map[string]interface{}{"ok": true, "err": ""}
	invitation, err = database.GetInvitationFromId(invitation.Id)
	if err != nil {
		ret["err"] = err.Error()
		ret["ok"] = false
	}
	ret["data"] = map[string]interface{}{"invitation": invitation}
	tools.WriteJsonResult(ret, w, http.StatusOK)
}

// Handle update requests from admin panel
func AdminGetFromId(w http.ResponseWriter, req *http.Request) {

	if views.RedirectIfNotAdmin(w, req) {
		return
	}

	params := mux.Vars(req)

	table := params["table"]

	switch table {
	case database.TableUsers:
		getUser(w, req)
	case database.TableGroups:
		getGroup(w, req)
	case database.TableVideos:
		getVideo(w, req)
	case database.TableVideoGroups:
		getVideoGroup(w, req)
	case database.TableInvitations:
		getInvitation(w, req)
	default:
		http.Error(w, "table is not valid", http.StatusBadRequest)
	}
}
