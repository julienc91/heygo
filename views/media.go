package views

import (
    "net/http"
    "github.com/gorilla/mux"
    "gomet/database"
)


func MediaHandler(w http.ResponseWriter, req *http.Request) {

    params := mux.Vars(req)
    mediaType := params["type"]
    slug := params["slug"]

    switch mediaType {
    case "video":
        video, err := database.GetVideoFromSlug(slug)
        if err != nil {
            panic(err)
        }
        http.ServeFile(w, req, video.Path)
    }
}
