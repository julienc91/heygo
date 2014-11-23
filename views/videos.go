package views

import (
	"github.com/gorilla/mux"
	"heygo/database"
	"heygo/tools"
	"html/template"
	"net/http"
)

// Display the video homepage
func VideoMenuHandler(w http.ResponseWriter, req *http.Request) {

	if RedirectIfNotAuthenticated(w, req) {
		return
	}

	var viewInfo = getViewInfo(req, "videos")

	t := template.Must(template.New("videos.html").ParseFiles(
		"templates/videos.html", "templates/base.html"))
	err := t.ExecuteTemplate(w, "base", viewInfo)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}

// Return a Json object with all the authorized videos
func VideoGetAllHandler(w http.ResponseWriter, req *http.Request) {

	if RedirectIfNotAuthenticated(w, req) {
		return
	}

	id := GetUserId(req)

	var ret = map[string]interface{}{"ok": true, "err": ""}

	rows, err := database.GetAllowedVideos(id)
	if err != nil {
		ret["err"] = err.Error()
		ret["ok"] = false
	}

	ret["data"] = rows
	writeJsonResult(ret, w, http.StatusOK)
}

// Return video informations
func VideoDetailHandler(w http.ResponseWriter, req *http.Request) {

	if RedirectIfNotAuthenticated(w, req) {
		return
	}

	params := mux.Vars(req)
	slug := params["slug"]

	video, err := database.PrepareGetFromKey("slug", slug, database.TableVideos)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	id := GetUserId(req)
	ok, err := database.CheckPermission(id, video["id"].(int64))
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	} else if !ok {
		http.Error(w, "", http.StatusForbidden)
		return
	}

	var ret = map[string]interface{}{"ok": true, "err": "", "data": video}
	writeJsonResult(ret, w, http.StatusOK)
}

// Call the OpenSubtitles API
func VideoGetSubtitles(w http.ResponseWriter, req *http.Request) {

	if RedirectIfNotAuthenticated(w, req) {
		return
	}

	params := mux.Vars(req)
	slug := params["slug"]

	video, err := database.PrepareGetFromKey("slug", slug, database.TableVideos)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	id := GetUserId(req)
	ok, err := database.CheckPermission(id, video["id"].(int64))
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	if !ok {
		http.Error(w, "", http.StatusForbidden)
		return
	}

	hash, size, err := tools.Hash(video["path"].(string))
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	res, err := tools.SearchSubtitles(hash, size)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/vtt")
	http.Error(w, res, http.StatusOK)
}
