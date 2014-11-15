package views

import (
	"github.com/gorilla/mux"
	"gomet/database"
	"io"
	"net/http"
	"os"
)

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
	case "video":
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

		f, err := os.Open(video["path"].(string))
		if err != nil {
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		io.Copy(w, f)
	}
}
