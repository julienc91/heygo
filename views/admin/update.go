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

func updateUser(w http.ResponseWriter, req *http.Request) {

	var user globals.User
	err := json.Unmarshal([]byte(req.FormValue("user")), &user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var groups []globals.Group
	err = json.Unmarshal([]byte(req.FormValue("groups")), &groups)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	oldUser, _, err := database.GetUserFromId(user.Id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if user.Password == "" {
		user.Password = oldUser.Password
	} else if oldUser.Password != user.Password {
		user.Salt = tools.SaltGenerator()
		user.Password = tools.Hash(user.Password, user.Salt)
	}

	var ret = map[string]interface{}{"ok": true, "err": ""}
	user, groups, err = database.UpdateUser(user, groups)
	if err != nil {
		ret["err"] = err.Error()
		ret["ok"] = false
	}
	ret["data"] = map[string]interface{}{"user": user, "groups": groups}
	tools.WriteJsonResult(ret, w, http.StatusOK)
}

func updateGroup(w http.ResponseWriter, req *http.Request) {

	var group globals.Group
	err := json.Unmarshal([]byte(req.FormValue("group")), &group)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var users []globals.User
	err = json.Unmarshal([]byte(req.FormValue("users")), &users)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var videoGroups []globals.VideoGroup
	err = json.Unmarshal([]byte(req.FormValue("video_groups")), &videoGroups)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var ret = map[string]interface{}{"ok": true, "err": ""}
	group, users, videoGroups, err = database.UpdateGroup(group, users, videoGroups)
	if err != nil {
		ret["err"] = err.Error()
		ret["ok"] = false
	}
	ret["data"] = map[string]interface{}{"group": group, "users": users, "video_groups": videoGroups}
	tools.WriteJsonResult(ret, w, http.StatusOK)
}

func updateVideo(w http.ResponseWriter, req *http.Request) {

	var video globals.Video
	err := json.Unmarshal([]byte(req.FormValue("video")), &video)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var videoGroups []globals.VideoGroup
	err = json.Unmarshal([]byte(req.FormValue("video_groups")), &videoGroups)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var ret = map[string]interface{}{"ok": true, "err": ""}
	video, videoGroups, err = database.UpdateVideo(video, videoGroups)
	if err != nil {
		ret["err"] = err.Error()
		ret["ok"] = false
	}
	ret["data"] = map[string]interface{}{"video": video, "video_groups": videoGroups}
	tools.WriteJsonResult(ret, w, http.StatusOK)
}

func updateVideoGroup(w http.ResponseWriter, req *http.Request) {

	var videoGroup globals.VideoGroup
	err := json.Unmarshal([]byte(req.FormValue("video_group")), &videoGroup)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var videos []globals.Video
	err = json.Unmarshal([]byte(req.FormValue("videos")), &videos)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var groups []globals.Group
	err = json.Unmarshal([]byte(req.FormValue("groups")), &groups)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var ret = map[string]interface{}{"ok": true, "err": ""}
	videoGroup, videos, groups, err = database.UpdateVideoGroup(videoGroup, videos, groups)
	if err != nil {
		ret["err"] = err.Error()
		ret["ok"] = false
	}
	ret["data"] = map[string]interface{}{"video_group": videoGroup, "videos": videos, "groups": groups}
	tools.WriteJsonResult(ret, w, http.StatusOK)
}

func updateInvitation(w http.ResponseWriter, req *http.Request) {

	var invitation globals.Invitation
	err := json.Unmarshal([]byte(req.FormValue("invitation")), &invitation)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var ret = map[string]interface{}{"ok": true, "err": ""}
	invitation, err = database.UpdateInvitation(invitation)
	if err != nil {
		ret["err"] = err.Error()
		ret["ok"] = false
	}
	ret["data"] = map[string]interface{}{"invitation": invitation}
	tools.WriteJsonResult(ret, w, http.StatusOK)
}

func updateConfiguration(w http.ResponseWriter, req *http.Request) {

	var configuration globals.Configuration
	err := json.Unmarshal([]byte(req.FormValue("configuration")), &configuration)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var ret = map[string]interface{}{"ok": true, "err": ""}
	err = database.UpdateConfiguration(configuration)
	if err != nil {
		ret["err"] = err.Error()
		ret["ok"] = false
	}
	ret["data"] = map[string]interface{}{"configuration": configuration}
	tools.WriteJsonResult(ret, w, http.StatusOK)
}

// Handle update requests from admin panel
func AdminUpdate(w http.ResponseWriter, req *http.Request) {

	if views.RedirectIfNotAdmin(w, req) {
		return
	}

	params := mux.Vars(req)

	table := params["table"]

	switch table {
	case database.TableUsers:
		updateUser(w, req)
	case database.TableGroups:
		updateGroup(w, req)
	case database.TableVideos:
		updateVideo(w, req)
	case database.TableVideoGroups:
		updateVideoGroup(w, req)
	case database.TableInvitations:
		updateInvitation(w, req)
	case database.TableConfiguration:
		updateConfiguration(w, req)
	default:
		http.Error(w, "table is not valid", http.StatusBadRequest)
	}
}
