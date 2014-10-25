package views

import (
	"net/http"
	"gomet/globals"
	"html/template"
)

type About struct {
	Appname string
	Version string
	Date string
	Author string
	Website string
}

func AboutHandler(w http.ResponseWriter, req *http.Request) {

    if RedirectIfNotAuthenticated(w, req) {
        return
    }

	t := template.Must(template.New("about.html").ParseFiles(
		"templates/about.html", "templates/base.html"))
	err := t.ExecuteTemplate(w, "base",
		About{globals.APPNAME, globals.VERSION, globals.DATE,
			globals.AUTHOR, globals.WEBSITE})
	if err != nil {
		panic(err)
	}
}
