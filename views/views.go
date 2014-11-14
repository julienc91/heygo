package views

import (
	"encoding/json"
	"gomet/database"
	"net/http"
)

func RedirectIfNotAdmin(w http.ResponseWriter, req *http.Request) bool {

	if RedirectIfNotAuthenticated(w, req) {
		return true
	}

	var id = GetUserId(req)
	ok, err := database.IsAdmin(id)
	if err != nil || !ok {
		http.Error(w, "", http.StatusForbidden)
		return true
	}

	return false
}

func RedirectIfNotAuthenticated(w http.ResponseWriter, req *http.Request) bool {

	var userId = GetUserId(req)
	if userId == 0 {
		http.Redirect(w, req, "/signin", http.StatusFound)
		return true
	}
	return false
}

func writeJsonResult(ret map[string]interface{}, w http.ResponseWriter, code int) {

	val, err := json.Marshal(ret)
	if err != nil {
		ret["err"] = err.Error()
		ret["ok"] = false
		code = http.StatusInternalServerError
	}

	http.Error(w, string(val), code)
}
