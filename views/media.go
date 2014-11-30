package views

import (
	"encoding/base64"
	"github.com/gorilla/mux"
	"heygo/database"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"time"
)

var thumbnailsToServe map[string]string

func init() {
	thumbnailsToServe = make(map[string]string)
}

// Stream a media
func MediaHandler(w http.ResponseWriter, req *http.Request) {

	if RedirectIfNotAuthenticated(w, req) {
		return
	}
	id := GetUserId(req)

	params := mux.Vars(req)
	mediaType := params["type"]
	slug := params["slug"]

	switch mediaType {
	case "videos":
		video, err := database.PrepareGetFromKey("slug", slug, database.TableVideos)
		if err != nil {
			http.Error(w, "Video not found", http.StatusNotFound)
			return
		}

		ok, err := database.CheckPermission(id, video["id"].(int64))
		if err != nil {
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
		if !ok {
			http.Error(w, "", http.StatusForbidden)
			return
		}

		var path = video["path"].(string)
		f, err := os.Open(path)
		if err != nil {
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
		defer f.Close()

		http.ServeContent(w, req, path, time.Time{}, f)
	}
}

// IMDB media proxy
func MediaThumbnailHandler(w http.ResponseWriter, req *http.Request) {

	if RedirectIfNotAuthenticated(w, req) {
		return
	}

	params := mux.Vars(req)
	mediaUrlB64 := params["url"]

	if mediaUrlB64 == "" {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	preloadedMedia, ok := thumbnailsToServe[mediaUrlB64]
	if ok {
		writeResponse(preloadedMedia, w, "image/jpg", http.StatusOK)
		return
	}

	mediaUrl, err := base64.StdEncoding.DecodeString(mediaUrlB64)
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	parsedUrl, err := url.ParseRequestURI(string(mediaUrl))
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	if !strings.HasSuffix(parsedUrl.Host, "media-imdb.com") || path.Ext(parsedUrl.Path) != ".jpg" || parsedUrl.Scheme != "http" {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	resp, err := http.Get(string(mediaUrl))
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	var thumbnail = string(data)
	thumbnailsToServe[mediaUrlB64] = thumbnail

	writeResponse(thumbnail, w, "image/jpg", http.StatusOK)
}
