package views

import (
	"github.com/gorilla/mux"
	"github.com/julienc91/heygo/database"
	"github.com/julienc91/heygo/tools"
	"html/template"
	"net/http"
)

var subtitlesToServe map[string]string
var salt string

func init() {
	subtitlesToServe = make(map[string]string)
	salt = tools.SaltGenerator()
}

// Display the video homepage
func VideoMenuHandler(w http.ResponseWriter, req *http.Request) {

	if RedirectIfNotAuthenticated(w, req) {
		return
	}

	var viewInfo = GetViewInfo(req, "videos")

	t := template.Must(template.New("videos.html").ParseFiles(
		"templates/videos.html", "templates/base.html"))
	err := t.ExecuteTemplate(w, "base", viewInfo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
	tools.WriteJsonResult(ret, w, http.StatusOK)
}

// Return video informations
func VideoDetailHandler(w http.ResponseWriter, req *http.Request) {

	if RedirectIfNotAuthenticated(w, req) {
		return
	}

	params := mux.Vars(req)
	slug := params["slug"]

	video, _, err := database.GetVideoFromSlug(slug)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id := GetUserId(req)
	ok, err := database.CheckPermission(id, video.Id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else if !ok {
		http.Error(w, "", http.StatusForbidden)
		return
	}

	var ret = map[string]interface{}{"ok": true, "err": "", "data": map[string]interface{}{"video": video}}
	tools.WriteJsonResult(ret, w, http.StatusOK)
}

// Call the OpenSubtitles API
func VideoGetSubtitles(w http.ResponseWriter, req *http.Request) {

	if RedirectIfNotAuthenticated(w, req) {
		return
	}

	params := mux.Vars(req)
	slug := params["slug"]
	lang := params["lang"]

	video, _, err := database.GetVideoFromSlug(slug)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id := GetUserId(req)
	ok, err := database.CheckPermission(id, video.Id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else if !ok {
		http.Error(w, "", http.StatusForbidden)
		return
	}

	hash, size, err := tools.OpensubtitlesHash(video.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res, err := tools.SearchSubtitles(hash, size, video.ImdbId, lang)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var ret = map[string]interface{}{"ok": true, "err": "", "data": []string{}}
	var data []string

	for _, subtitles := range res {
		hash = tools.Hash(subtitles, salt)
		subtitlesToServe[hash] = subtitles
		data = append(data, hash)
	}
	ret["data"] = data

	tools.WriteJsonResult(ret, w, http.StatusOK)
}

// Serve subtitle files
func SubtitlesServerHandler(w http.ResponseWriter, req *http.Request) {

	if RedirectIfNotAuthenticated(w, req) {
		return
	}

	params := mux.Vars(req)
	hash := params["hash"]

	if _, ok := subtitlesToServe[hash]; !ok {
		http.Error(w, "", http.StatusNotFound)
		return
	}
	tools.WriteResponse(subtitlesToServe[hash], w, "text/vtt", http.StatusOK)
}
