package main

import (
	"net/http"
    "github.com/gorilla/mux"
	"gomet/views"
)

func main() {

    var rtr = mux.NewRouter()
    rtr = rtr.StrictSlash(true)
    
    rtr.HandleFunc("/about", views.AboutHandler)
    rtr.HandleFunc("/videos/gethash/{slug:[a-zA-Z0-9_]+}", views.VideoGetHash)
    rtr.HandleFunc("/videos/watch/{slug:[a-zA-Z0-9_]+}", views.VideoDetailHandler)
    rtr.HandleFunc("/videos", views.VideoMenuHandler)
    rtr.HandleFunc("/signin", views.SignInHandler)
    rtr.HandleFunc("/signup", views.SignupHandler)
    rtr.HandleFunc("/signout", views.SignoutHandler)
    rtr.HandleFunc("/login", views.LoginHandler)
    rtr.HandleFunc("/media/{type:[a-z]+}/{slug:[a-zA-Z0-9_]+}", views.MediaHandler)
    rtr.HandleFunc("/admin/update", views.AdminUpdateHandler)
    rtr.HandleFunc("/admin/insert", views.AdminInsertHandler)
    rtr.HandleFunc("/admin/delete", views.AdminDeleteHandler)
    rtr.HandleFunc("/admin/delete_pivot", views.AdminDeletePivotHandler)
    rtr.HandleFunc("/admin", views.AdminHandler)
    rtr.HandleFunc("/", views.MainPageHandler)
    
    rtr.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
    
    http.Handle("/", rtr)
    http.ListenAndServe(":8080", nil)
}
