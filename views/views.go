package views

import (
    "net/http"
    "gomet/database"
)

func RedirectIfNotAdmin(w http.ResponseWriter, req *http.Request) bool {

    if RedirectIfNotAuthenticated(w, req) {
        return true
    }

    var id = GetUserId(req)
    ok, err := database.IsAdmin(id)
    if err != nil || !ok {
        http.Redirect(w, req, "/", 302)
        return true
    }
    
    return false
}


func RedirectIfNotAuthenticated(w http.ResponseWriter, req *http.Request) bool {
    
    var userId = GetUserId(req)
    if userId == 0 {
        http.Redirect(w, req, "/signin", 302)
        return true
    }
    return false
}
