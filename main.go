package main

import (
	"github.com/gorilla/mux"
	"github.com/julienc91/heygo/database"
	"github.com/julienc91/heygo/globals"
	"github.com/julienc91/heygo/tools"
	"github.com/julienc91/heygo/views"
	"github.com/julienc91/heygo/views/admin"
	"net/http"
)

// The main function, from where routes are defined
func main() {

	defer tools.OpenSubtitlesClose()
	defer close(globals.LoadConfiguration)

	database.LoadConfiguration()
	tools.InitFromConfiguration()
	go reloadConfiguration()

	var rtr = mux.NewRouter()
	rtr = rtr.StrictSlash(true)

	rtr.HandleFunc("/about", views.AboutHandler)
	rtr.HandleFunc("/videos/get/{slug:[a-z0-9_]+}", views.VideoDetailHandler)
	rtr.HandleFunc("/videos/get", views.VideoGetAllHandler)
	rtr.HandleFunc("/videos", views.VideoMenuHandler)
	rtr.HandleFunc("/signin", views.SignInHandler)
	rtr.HandleFunc("/signup", views.Signup)
	rtr.HandleFunc("/signout", views.Signout)
	rtr.HandleFunc("/login", views.Login)
	rtr.HandleFunc("/media/{type:videos}/{slug:[a-z0-9_]+}", views.StreamMedia)
	rtr.HandleFunc("/media/thumbnail/{url}", views.MediaThumbnailHandler)
	rtr.HandleFunc("/media/subtitles/list/{slug:[a-z0-9_]+}/{lang:fre|eng}", views.VideoGetSubtitles)
	rtr.HandleFunc("/media/subtitles/get/{hash:[a-zA-Z0-9/_=-]+}", views.SubtitlesServerHandler)
	rtr.HandleFunc("/admin/get/{table:[a-z_]+}", admin.AdminGetAll)
	rtr.HandleFunc("/admin/get/{table:[a-z_]+}/{id:[0-9]+}", admin.AdminGetFromId)
	rtr.HandleFunc("/admin/update/{table:[a-z_]+}/{id:[0-9]+}", admin.AdminUpdate)
	rtr.HandleFunc("/admin/insert/{table:[a-z_]+}", admin.AdminInsert)
	rtr.HandleFunc("/admin/batchinsert/{table:videos}", admin.AdminBatchInsertVideos)
	rtr.HandleFunc("/admin/delete/{table:[a-z_]+}", admin.AdminDelete)
	rtr.HandleFunc("/admin/media/check", admin.AdminMediaCheck)
	rtr.HandleFunc("/admin", admin.AdminHandler)
	rtr.HandleFunc("/", views.MainPageHandler)

	// serve static files
	rtr.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	http.Handle("/", rtr)
	http.ListenAndServe(globals.CONFIGURATION.Domain+":"+globals.CONFIGURATION.Port, nil)
}

// Hot reloading of the configuration
func reloadConfiguration() {
	for _ = range globals.LoadConfiguration {
		database.LoadConfiguration()
		tools.InitFromConfiguration()
	}
}
