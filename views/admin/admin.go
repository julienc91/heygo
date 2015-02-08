package admin

import (
	"github.com/julienc91/heygo/globals"
	"github.com/julienc91/heygo/tools"
	"github.com/julienc91/heygo/views"
	"html/template"
	"net/http"
)

// Return the configuration variables
func AdminGetConfigurationHandler(w http.ResponseWriter, req *http.Request) {

	if views.RedirectIfNotAdmin(w, req) {
		return
	}

	var ret = map[string]interface{}{"ok": true, "err": ""}
	ret["data"] = globals.CONFIGURATION
	tools.WriteJsonResult(ret, w, http.StatusOK)
}

// Handle media checking requests
func AdminMediaCheck(w http.ResponseWriter, req *http.Request) {

	if views.RedirectIfNotAdmin(w, req) {
		return
	}

	path := req.FormValue("path")

	var ret = map[string]interface{}{"ok": tools.CheckFilePath(path)}
	tools.WriteJsonResult(ret, w, http.StatusOK)
}

// Display the admin panel
func AdminHandler(w http.ResponseWriter, req *http.Request) {

	if views.RedirectIfNotAdmin(w, req) {
		return
	}

	var viewInfo = views.GetViewInfo(req, "admin")

	t := template.Must(template.New("admin.html").ParseFiles(
		"templates/admin.html", "templates/base.html"))
	err := t.ExecuteTemplate(w, "base", viewInfo)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
	}
}
