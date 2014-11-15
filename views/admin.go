package views

import (
	"github.com/gorilla/mux"
	"gomet/database"
	"gomet/tools"
	"html/template"
	"net/http"
	"strconv"
)

// Handle update requests from admin panel
func AdminUpdateHandler(w http.ResponseWriter, req *http.Request) {

	if RedirectIfNotAdmin(w, req) {
		return
	}

	params := mux.Vars(req)

	table := params["table"]
	if !tools.InArray(database.MainTables, table) {
		http.Error(w, "table is not valid", http.StatusBadRequest)
		return
	}

	req.ParseForm()

	var data = make(map[string]interface{})
	for k := range req.Form {
		data[k] = req.FormValue(k)
	}

	id, err := strconv.ParseInt(params["id"], 10, 64)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	var ret = map[string]interface{}{"ok": true, "err": ""}

	updated, err := database.PrepareUpdate(id, data, table)
	if err != nil {
		ret["err"] = err.Error()
		ret["ok"] = false
	}

	ret["data"] = updated
	writeJsonResult(ret, w, http.StatusOK)
}

// Handle deletion requests from admin panel
func AdminDeleteHandler(w http.ResponseWriter, req *http.Request) {

	if RedirectIfNotAdmin(w, req) {
		return
	}

	params := mux.Vars(req)

	table := params["table"]
	if !tools.InArray(database.MainTables, table) {
		http.Error(w, "table is not valid", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(params["id"], 10, 64)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	var ret = map[string]interface{}{"ok": true, "err": ""}
	if err = database.PrepareDeleteFromId(id, table); err != nil {
		ret["err"] = err.Error()
		ret["ok"] = false
	}

	writeJsonResult(ret, w, http.StatusOK)
}

// Handle insert requests from admin panel
func AdminInsertHandler(w http.ResponseWriter, req *http.Request) {

	if RedirectIfNotAdmin(w, req) {
		return
	}

	params := mux.Vars(req)

	table := params["table"]
	if !tools.InArray(database.MainTables, table) {
		http.Error(w, "table is not valid", http.StatusBadRequest)
		return
	}

	req.ParseForm()

	var data = make(map[string]interface{})
	for k := range req.Form {
		data[k] = req.FormValue(k)
	}

	var ret = map[string]interface{}{"ok": true, "err": ""}

	inserted, err := database.PrepareInsert(data, table)
	if err != nil {
		ret["err"] = err.Error()
		ret["ok"] = false
	}

	ret["data"] = inserted
	writeJsonResult(ret, w, http.StatusOK)
}

// Handle get requests from admin panel
func AdminGetHandler(w http.ResponseWriter, req *http.Request) {

	if RedirectIfNotAdmin(w, req) {
		return
	}

	params := mux.Vars(req)

	table := params["table"]
	if !tools.InArray(database.Tables, table) {
		http.Error(w, "table is not valid", http.StatusBadRequest)
		return
	}

	req.ParseForm()

	var column = req.FormValue("column")
	var filter = req.FormValue("filter")
	var value = req.FormValue("value")
	var rows interface{}
	var ret = map[string]interface{}{"ok": true, "err": ""}
	var err error

	if column == "" || filter == "" || value == "" {
		rows, err = database.PrepareGetAll(table)
	} else {
		valueId, convErr := strconv.ParseInt(value, 10, 64)
		if convErr != nil {
			http.Error(w, "", http.StatusBadRequest)
			return
		}
		rows, err = database.PrepareGetColumnFiltered(column, filter, valueId, table)
	}

	if err != nil {
		ret["err"] = err.Error()
		ret["ok"] = false
	}

	ret["data"] = rows
	writeJsonResult(ret, w, http.StatusOK)
}

// Handle get_from_id requests from admin panel
func AdminGetFromIdHandler(w http.ResponseWriter, req *http.Request) {

	if RedirectIfNotAdmin(w, req) {
		return
	}

	params := mux.Vars(req)
	table := params["table"]
	if !tools.InArray(database.MainTables, table) {
		http.Error(w, "table is not valid", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(params["id"], 10, 64)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	var ret = map[string]interface{}{"ok": true, "err": ""}

	rows, err := database.PrepareGetFromId(id, table)
	if err != nil {
		ret["err"] = err.Error()
		ret["ok"] = false
	}

	ret["data"] = rows
	writeJsonResult(ret, w, http.StatusOK)
}

// Set pivot table values
func AdminSetHandler(w http.ResponseWriter, req *http.Request) {

	if RedirectIfNotAdmin(w, req) {
		return
	}
	params := mux.Vars(req)
	table := params["table"]
	if !tools.InArray(database.PivotTables, table) {
		http.Error(w, "table is not valid", http.StatusBadRequest)
		return
	}

	req.ParseForm()

	var filter = req.FormValue("filter")
	var column = req.FormValue("column")
	var value = req.FormValue("value")
	var values = req.Form["values"]

	valueId, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	if err := database.PrepareDeleteFromFilter(filter, value, table); err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	for _, key := range values {
		keyId, err := strconv.ParseInt(key, 10, 64)
		if err != nil {
			http.Error(w, "", http.StatusBadRequest)
			return
		}
		data := map[string]interface{}{filter: valueId, column: keyId}
		if _, err := database.PrepareInsert(data, table); err != nil {
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}
	AdminGetHandler(w, req)
}

// Handle media checking requests
func AdminMediaCheckHandler(w http.ResponseWriter, req *http.Request) {

	if RedirectIfNotAdmin(w, req) {
		return
	}

	path := req.FormValue("path")

	var ret = map[string]interface{}{"ok": tools.CheckFilePath(path)}
	writeJsonResult(ret, w, http.StatusOK)
}

// Display the admin panel
func AdminHandler(w http.ResponseWriter, req *http.Request) {

	if RedirectIfNotAdmin(w, req) {
		return
	}

	t := template.Must(template.New("admin.html").ParseFiles(
		"templates/admin.html", "templates/base.html"))
	err := t.ExecuteTemplate(w, "base", nil)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
	}
}
