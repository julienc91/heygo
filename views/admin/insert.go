package admin

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/julienc91/heygo/database"
	"github.com/julienc91/heygo/globals"
	"github.com/julienc91/heygo/tools"
	"github.com/julienc91/heygo/views"
	"net/http"
	"path"
	"strconv"
)

func insertUser(w http.ResponseWriter, req *http.Request) {

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

	user.Salt = tools.SaltGenerator()
	user.Password = tools.Hash(user.Password, user.Salt)

	var ret = map[string]interface{}{"ok": true, "err": ""}
	user, groups, err = database.InsertUser(user, groups)
	if err != nil {
		ret["err"] = err.Error()
		ret["ok"] = false
	}
	ret["data"] = map[string]interface{}{"user": user, "groups": groups}
	tools.WriteJsonResult(ret, w, http.StatusOK)
}

func insertGroup(w http.ResponseWriter, req *http.Request) {

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
	group, users, videoGroups, err = database.InsertGroup(group, users, videoGroups)
	if err != nil {
		ret["err"] = err.Error()
		ret["ok"] = false
	}
	ret["data"] = map[string]interface{}{"group": group, "users": users, "video_groups": videoGroups}
	tools.WriteJsonResult(ret, w, http.StatusOK)
}

func insertVideo(w http.ResponseWriter, req *http.Request) {

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
	video, videoGroups, err = database.InsertVideo(video, videoGroups)
	if err != nil {
		ret["err"] = err.Error()
		ret["ok"] = false
	}
	ret["data"] = map[string]interface{}{"video": video, "video_groups": videoGroups}
	tools.WriteJsonResult(ret, w, http.StatusOK)
}

func insertVideoGroup(w http.ResponseWriter, req *http.Request) {

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
	videoGroup, videos, groups, err = database.InsertVideoGroup(videoGroup, videos, groups)
	if err != nil {
		ret["err"] = err.Error()
		ret["ok"] = false
	}
	ret["data"] = map[string]interface{}{"video_group": videoGroup, "videos": videos, "groups": groups}
	tools.WriteJsonResult(ret, w, http.StatusOK)
}

func insertInvitation(w http.ResponseWriter, req *http.Request) {

	var invitation globals.Invitation
	err := json.Unmarshal([]byte(req.FormValue("invitation")), &invitation)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var ret = map[string]interface{}{"ok": true, "err": ""}
	invitation, err = database.InsertInvitation(invitation)
	if err != nil {
		ret["err"] = err.Error()
		ret["ok"] = false
	}
	ret["data"] = map[string]interface{}{"invitation": invitation}
	tools.WriteJsonResult(ret, w, http.StatusOK)
}

// Handle insert requests from admin panel
func AdminInsert(w http.ResponseWriter, req *http.Request) {

	if views.RedirectIfNotAdmin(w, req) {
		return
	}

	params := mux.Vars(req)

	table := params["table"]

	switch table {
	case database.TableUsers:
		insertUser(w, req)
	case database.TableGroups:
		insertGroup(w, req)
	case database.TableVideos:
		insertVideo(w, req)
	case database.TableVideoGroups:
		insertVideoGroup(w, req)
	case database.TableInvitations:
		insertInvitation(w, req)
	default:
		http.Error(w, "table is not valid", http.StatusBadRequest)
	}
}

// Create media from the given subfolder
func AdminBatchInsertVideos(w http.ResponseWriter, req *http.Request) {

	if views.RedirectIfNotAdmin(w, req) {
		return
	}

	var filepath = req.FormValue("path")
	var extension = req.FormValue("extension")

	recursive, err := strconv.ParseBool(req.FormValue("recursive"))
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

	files, err := tools.GetFilesFromSubfolder(filepath, extension, recursive)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var ret = map[string]interface{}{"ok": true, "err": ""}
	var data []globals.Video

	for _, filename := range files {
		var video globals.Video
		video.Path = filename
		video.Title = path.Base(filename)
		video.Slug = tools.SlugFromFilename(path.Base(filename))

		video, _, err := database.InsertVideo(video, videoGroups)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		data = append(data, video)
	}

	ret["data"] = map[string]interface{}{"videos": data}
	tools.WriteJsonResult(ret, w, http.StatusOK)
}
