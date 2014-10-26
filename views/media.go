package views

import (
    "net/http"
    "github.com/gorilla/mux"
    "gomet/database"
    "io"
    "os"
    "errors"
)


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
        video, err := database.GetVideoFromSlug(slug)
        if err != nil {
            panic(err)
        }

        ok, err := database.CheckPermission(id, video.Id)
        if err != nil {
            panic(err)
        }
        if !ok {
            panic(errors.New("Forbidden"))
        }

        f, err := os.Open(video.Path)
        if err != nil {
            panic(err)
        }

        io.Copy(w, f)
        //http.ServeFile(w, req, video.Path)
    }
}
