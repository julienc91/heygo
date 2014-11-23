package views

import (
	"heygo/globals"
	"html/template"
	"net/http"
)

type About struct {
	Appname             string
	Version             string
	Date                string
	Author              string
	Website             string
	IsUserAuthenticated bool
	IsUserAdmin         bool
	ViewName            string
}

// Display the "About" page
func AboutHandler(w http.ResponseWriter, req *http.Request) {

	var viewInfo = getViewInfo(req, "about")

	t := template.Must(template.New("about.html").ParseFiles(
		"templates/about.html", "templates/base.html"))
	err := t.ExecuteTemplate(w, "base",
		About{globals.APPNAME, globals.VERSION, globals.DATE, globals.AUTHOR, globals.WEBSITE, viewInfo.IsUserAuthenticated,
			viewInfo.IsUserAdmin, viewInfo.ViewName})
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
	}
}
