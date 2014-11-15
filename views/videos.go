package views

import (
	"github.com/gorilla/mux"
	"gomet/database"
	"gomet/tools"
	"html/template"
	"net/http"
)

func VideoMenuHandler(w http.ResponseWriter, req *http.Request) {

	if RedirectIfNotAuthenticated(w, req) {
		return
	}

	t := template.Must(template.New("videos.html").ParseFiles(
		"templates/videos.html", "templates/base.html"))
	err := t.ExecuteTemplate(w, "base", nil)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}

func VideoDetailHandler(w http.ResponseWriter, req *http.Request) {

	if RedirectIfNotAuthenticated(w, req) {
		return
	}

	params := mux.Vars(req)
	slug := params["slug"]

	video, err := database.PrepareGetFromKey("slug", slug, database.TableVideos)
	if err != nil {
		panic(err)
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

	t := template.Must(template.New("video_detail.html").ParseFiles(
		"templates/video_detail.html", "templates/base.html"))
	err = t.ExecuteTemplate(w, "base", video)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}

func VideoGetHash(w http.ResponseWriter, req *http.Request) {

	if RedirectIfNotAuthenticated(w, req) {
		return
	}

	params := mux.Vars(req)
	slug := params["slug"]

	video, err := database.PrepareGetFromKey("slug", slug, database.TableVideos)
	if err != nil {
		panic(err)
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

	var ret = map[string]interface{}{"ok": true, "err": ""}
	hash, size, err := tools.Hash(video["path"].(string))
	if err != nil {
		ret["err"] = err.Error()
		ret["ok"] = false
	}
	ret["hash"] = hash
	ret["size"] = size
	writeJsonResult(ret, w, http.StatusOK)
}
