package main

import (
	"github.com/gorilla/mux"
	"heygo/tools"
	"heygo/views"
	"net/http"
)

// The main function, from where routes are defined
func main() {

	defer tools.OpenSubtitlesClose()

	var rtr = mux.NewRouter()
	rtr = rtr.StrictSlash(true)

	rtr.HandleFunc("/about", views.AboutHandler)
	rtr.HandleFunc("/videos/getsubtitles/{slug:[a-z0-9_]+}", views.VideoGetSubtitles)
	rtr.HandleFunc("/videos/get/{slug:[a-z0-9_]+}", views.VideoDetailHandler)
	rtr.HandleFunc("/videos/get", views.VideoGetAllHandler)
	rtr.HandleFunc("/videos", views.VideoMenuHandler)
	rtr.HandleFunc("/signin", views.SignInHandler)
	rtr.HandleFunc("/signup", views.SignupHandler)
	rtr.HandleFunc("/signout", views.SignoutHandler)
	rtr.HandleFunc("/login", views.LoginHandler)
	rtr.HandleFunc("/media/{type:videos}/{slug:[a-z0-9_]+}", views.MediaHandler)
	rtr.HandleFunc("/admin/get/{table:[a-z_]+}", views.AdminGetHandler)
	rtr.HandleFunc("/admin/get/{table:[a-z_]+}/{id:[0-9]+}", views.AdminGetFromIdHandler)
	rtr.HandleFunc("/admin/set/{table:[a-z_]+}", views.AdminSetHandler)
	rtr.HandleFunc("/admin/update/{table:[a-z_]+}/{id:[0-9]+}", views.AdminUpdateHandler)
	rtr.HandleFunc("/admin/insert/{table:[a-z_]+}", views.AdminInsertHandler)
	rtr.HandleFunc("/admin/delete/{table:[a-z_]+}/{id:[0-9]+}", views.AdminDeleteHandler)
	rtr.HandleFunc("/admin/media/check", views.AdminMediaCheckHandler)
	rtr.HandleFunc("/admin", views.AdminHandler)
	rtr.HandleFunc("/", views.MainPageHandler)

	// serve static files
	rtr.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	http.Handle("/", rtr)
	http.ListenAndServe(":8080", nil)
}
