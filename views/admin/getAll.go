package admin

import (
	"github.com/gorilla/mux"
	"github.com/julienc91/heygo/database"
	"github.com/julienc91/heygo/tools"
	"github.com/julienc91/heygo/views"
	"net/http"
)

func getAllUsers(w http.ResponseWriter, req *http.Request) {

	users, err := database.GetAllUsers()
	var ret = map[string]interface{}{"ok": true, "err": ""}
	if err != nil {
		ret["err"] = err.Error()
		ret["ok"] = false
	}

	ret["data"] = map[string]interface{}{"users": users}
	tools.WriteJsonResult(ret, w, http.StatusOK)
}

func getAllGroups(w http.ResponseWriter, req *http.Request) {

	var ret = map[string]interface{}{"ok": true, "err": ""}
	groups, err := database.GetAllGroups()
	if err != nil {
		ret["err"] = err.Error()
		ret["ok"] = false
	}
	ret["data"] = map[string]interface{}{"groups": groups}
	tools.WriteJsonResult(ret, w, http.StatusOK)
}

func getAllVideos(w http.ResponseWriter, req *http.Request) {

	var ret = map[string]interface{}{"ok": true, "err": ""}
	videos, err := database.GetAllVideos()
	if err != nil {
		ret["err"] = err.Error()
		ret["ok"] = false
	}
	ret["data"] = map[string]interface{}{"videos": videos}
	tools.WriteJsonResult(ret, w, http.StatusOK)
}

func getAllVideoGroups(w http.ResponseWriter, req *http.Request) {

	var ret = map[string]interface{}{"ok": true, "err": ""}
	videoGroups, err := database.GetAllVideoGroups()
	if err != nil {
		ret["err"] = err.Error()
		ret["ok"] = false
	}
	ret["data"] = map[string]interface{}{"video_groups": videoGroups}
	tools.WriteJsonResult(ret, w, http.StatusOK)
}

func getAllInvitations(w http.ResponseWriter, req *http.Request) {

	var ret = map[string]interface{}{"ok": true, "err": ""}
	invitations, err := database.GetAllInvitations()
	if err != nil {
		ret["err"] = err.Error()
		ret["ok"] = false
	}
	ret["data"] = map[string]interface{}{"invitations": invitations}
	tools.WriteJsonResult(ret, w, http.StatusOK)
}

// Handle update requests from admin panel
func AdminGetAll(w http.ResponseWriter, req *http.Request) {

	if views.RedirectIfNotAdmin(w, req) {
		return
	}

	params := mux.Vars(req)

	table := params["table"]

	switch table {
	case database.TableUsers:
		getAllUsers(w, req)
	case database.TableGroups:
		getAllGroups(w, req)
	case database.TableVideos:
		getAllVideos(w, req)
	case database.TableVideoGroups:
		getAllVideoGroups(w, req)
	case database.TableInvitations:
		getAllInvitations(w, req)
	default:
		http.Error(w, "table is not valid", http.StatusBadRequest)
	}
}
