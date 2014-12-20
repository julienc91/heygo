package main

import (
	"github.com/gorilla/mux"
	"heygo/database"
	"heygo/globals"
	"heygo/tools"
	"heygo/views"
	"net/http"
)

// The main function, from where routes are defined
func main() {

	defer tools.OpenSubtitlesClose()
	defer close(globals.LoadConfiguration)
	loadConfiguration()
	go reloadConfiguration()

	var rtr = mux.NewRouter()
	rtr = rtr.StrictSlash(true)

	rtr.HandleFunc("/about", views.AboutHandler)
	rtr.HandleFunc("/videos/get/{slug:[a-z0-9_]+}", views.VideoDetailHandler)
	rtr.HandleFunc("/videos/get", views.VideoGetAllHandler)
	rtr.HandleFunc("/videos", views.VideoMenuHandler)
	rtr.HandleFunc("/signin", views.SignInHandler)
	rtr.HandleFunc("/signup", views.SignupHandler)
	rtr.HandleFunc("/signout", views.SignoutHandler)
	rtr.HandleFunc("/login", views.LoginHandler)
	rtr.HandleFunc("/media/{type:videos}/{slug:[a-z0-9_]+}", views.MediaHandler)
	rtr.HandleFunc("/media/thumbnail/{url}", views.MediaThumbnailHandler)
	rtr.HandleFunc("/media/subtitles/list/{slug:[a-z0-9_]+}/{lang:fre|eng}", views.VideoGetSubtitles)
	rtr.HandleFunc("/media/subtitles/get/{hash:[a-zA-Z0-9/_=-]+}", views.SubtitlesServerHandler)
	rtr.HandleFunc("/admin/get/configuration", views.AdminGetConfiguration)
	rtr.HandleFunc("/admin/get/{table:[a-z_]+}", views.AdminGetAll)
	rtr.HandleFunc("/admin/get/{table:[a-z_]+}/{id:[0-9]+}", views.AdminGetFromId)
	rtr.HandleFunc("/admin/update/{table:[a-z_]+}/{id:[0-9]+}", views.AdminUpdate)
	rtr.HandleFunc("/admin/insert/{table:[a-z_]+}", views.AdminInsert)
	rtr.HandleFunc("/admin/batchinsert/{table:videos}", views.AdminBatchInsertVideosHandler)
	rtr.HandleFunc("/admin/delete/{table:[a-z_]+}/{id:[0-9]+}", views.AdminDelete)
	rtr.HandleFunc("/admin/media/check", views.AdminMediaCheck)
	rtr.HandleFunc("/admin", views.AdminHandler)
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
