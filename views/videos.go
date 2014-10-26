package views

import (
    "errors"
	"net/http"
	"html/template"
    "github.com/gorilla/mux"
    "gomet/database"
    "gomet/tools"
    "fmt"
)


type VideoList struct {
    Videos []database.Video
}

func VideoMenuHandler(w http.ResponseWriter, req *http.Request) {

    if RedirectIfNotAuthenticated(w, req) {
        return
    }

    videos, err := database.GetAllowedVideos(GetUserId(req))
    if err != nil {
        panic(err)
    }

	t := template.Must(template.New("videos.html").ParseFiles(
		"templates/videos.html", "templates/base.html"))
	err = t.ExecuteTemplate(w, "base", VideoList{videos})
	if err != nil {
		panic(err)
	}
}


func VideoDetailHandler(w http.ResponseWriter, req *http.Request) {

    if RedirectIfNotAuthenticated(w, req) {
        return
    }
    
    params := mux.Vars(req)
    slug := params["slug"]

    video, err := database.GetVideoFromSlug(slug)
    if err != nil {
        panic(err)
    }

    id := GetUserId(req)
    ok, err := database.CheckPermission(id, video.Id)
    if err != nil {
        panic(err)
    }
    if !ok {
        panic(errors.New("Forbidden"))
    }

    t := template.Must(template.New("video_detail.html").ParseFiles(
        "templates/video_detail.html", "templates/base.html"))
    err = t.ExecuteTemplate(w, "base", video)
    if err != nil {
        panic(err)
    }
}


func VideoGetHash(w http.ResponseWriter, req *http.Request) {

    if RedirectIfNotAuthenticated(w, req) {
        return
    }
    
    params := mux.Vars(req)
    slug := params["slug"]

    video, err := database.GetVideoFromSlug(slug)
    if err != nil {
        panic(err)
    }

    id := GetUserId(req)
    ok, err := database.CheckPermission(id, video.Id)
    if err != nil {
        panic(err)
    }
    if !ok {
        panic(errors.New("Forbidden"))
    }

    hash, size, err := tools.Hash(video.Path)
    if err != nil {
        fmt.Fprintf(w, "{\"ok\": false, \"err\": \"%s\"}", err.Error())
    }
    fmt.Fprintf(w, "{\"ok\": true, \"hash\": \"%s\", \"size\": %d}", hash, size)
}
