// This package implements all the functions which will answer users' requests
package views

import (
	"github.com/julienc91/heygo/database"
	"net/http"
)

// Basic informations about the view to display and the user
type ViewInfo struct {
	IsUserAuthenticated bool
	IsUserAdmin         bool
	ViewName            string
}

// Handle the homepage
func MainPageHandler(w http.ResponseWriter, req *http.Request) {

	if RedirectIfNotAuthenticated(w, req) {
		return
	}
	http.Redirect(w, req, "/about", http.StatusFound)
}

// Check if the user has admin rights.
// If he is not authenticated, he will be redirected to the signin page
// If he is authenticated but without admin rights, a StatusForbidden code is set
// The function returns whether or not the user is unauthorized
func RedirectIfNotAdmin(w http.ResponseWriter, req *http.Request) bool {

	if RedirectIfNotAuthenticated(w, req) {
		return true
	}

	var id = GetUserId(req)
	ok, err := database.IsAdmin(id)
	if err != nil || !ok {
		http.Error(w, "", http.StatusForbidden)
		return true
	}

	return false
}

// Check if the user is authenticated
// If he is not authenticated, he will be redirected to the signin page
// The function returns whether or not the user is unauthorized
func RedirectIfNotAuthenticated(w http.ResponseWriter, req *http.Request) bool {

	var userId = GetUserId(req)
	if userId == 0 {
		http.Redirect(w, req, "/signin", http.StatusFound)
		return true
	}
	return false
}

// Fille a ViewInfo variable
func GetViewInfo(req *http.Request, viewName string) ViewInfo {

	var viewInfo ViewInfo
	var id = GetUserId(req)
	viewInfo.IsUserAuthenticated = id > 0
	ok, err := database.IsAdmin(id)
	viewInfo.IsUserAdmin = (err == nil && ok)
	viewInfo.ViewName = viewName
	return viewInfo
}
