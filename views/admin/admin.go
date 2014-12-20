package admin

import (
	"github.com/gorilla/mux"
	"heygo/database"
	"heygo/globals"
	"heygo/tools"
	"html/template"
	"net/http"
	"path"
	"strconv"
)

func updateUser(w http.ResponseWriter, req *httpRequest) {

	if _, ok := req.Value("user"); !ok {
		http.Error(w, "user is not set", http.StatusBadRequest)
		return
	}
	if _, ok := req.Value("groups"); !ok {
		http.Error(w, "groups are not set", http.StatusBadRequest)
		return
	}
	var user globals.User
	err := json.Unmarshal([]byte(req.Value("user")), &user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var groups []globals.Group
	err := json.Unmarshal([]byte(req.Value("groups")), &groups))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	oldUser, err := database.GetUserFromId(user.Id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if user.Password == "" {
		user.Password = oldUser.Password
	} else if oldUser.Password != user.Password {
		user.Salt = tools.SaltGenerator()
		user.Password = tools.Hash(user.Password, salt)
	}
}

// Handle update requests from admin panel
func AdminUpdate(w http.ResponseWriter, req *http.Request) {

	if RedirectIfNotAdmin(w, req) {
		return
	}

	params := mux.Vars(req)

	table := params["table"]
	if !tools.InArray(database.MainTables, table) {
		http.Error(w, "table is not valid", http.StatusBadRequest)
		return
	}

	err := req.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var ret = map[string]interface{}{"ok": true, "err": ""}

	switch table {
		case database.TableUsers:

		default:

	}

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

	updated, err := database.PrepareUpdateFromId(id, data, table)
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

// Create media from the given subfolder
func AdminBatchInsertVideosHandler(w http.ResponseWriter, req *http.Request) {

	if RedirectIfNotAdmin(w, req) {
		return
	}

	params := mux.Vars(req)

	table := params["table"]
	if !tools.InArray([]string{database.TableVideos}, table) {
		http.Error(w, "table is not valid", http.StatusBadRequest)
		return
	}

	req.ParseForm()

	var filepath = req.FormValue("path")
	var extension = req.FormValue("extension")
	var filter = req.FormValue("filter")
	var column = req.FormValue("column")
	var values = req.Form["values"]
	var pivotTable = req.FormValue("pivot_table")
	if !tools.InArray(database.PivotTables, pivotTable) {
		http.Error(w, "pivot table is not valid", http.StatusBadRequest)
		return
	}
	recursive, err := strconv.ParseBool(req.FormValue("recursive"))
	if err != nil {
		http.Error(w, "recursive is not valid", http.StatusBadRequest)
		return
	}

	files, err := tools.GetFilesFromSubfolder(filepath, extension, recursive)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var ret = map[string]interface{}{"ok": true, "err": ""}
	var data []map[string]interface{}

	for _, filename := range files {
		var params = map[string]interface{}{
			"path":  filename,
			"title": path.Base(filename),
			"slug":  tools.SlugFromFilename(path.Base(filename))}
		inserted, err := database.PrepareInsert(params, table)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		insertedId := inserted["id"].(int64)

		for _, key := range values {
			keyId, err := strconv.ParseInt(key, 10, 64)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			data := map[string]interface{}{filter: insertedId, column: keyId}
			if _, err := database.PrepareInsert(data, pivotTable); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		data = append(data, inserted)
	}

	ret["data"] = data
	writeJsonResult(ret, w, http.StatusOK)
}

// Return the configuration variables
func AdminGetConfigurationHandler(w http.ResponseWriter, req *http.Request) {

	if RedirectIfNotAdmin(w, req) {
		return
	}

	var ret = map[string]interface{}{"ok": true, "err": ""}
	ret["data"] = globals.CONFIGURATION
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

// Set configuration values
func AdminSetConfigurationHandler(w http.ResponseWriter, req *http.Request) {

	if RedirectIfNotAdmin(w, req) {
		return
	}
	req.ParseForm()

	var key = req.FormValue("key")
	var value = req.FormValue("value")

	var ret = map[string]interface{}{"ok": true, "err": ""}

	res, err := database.PrepareUpdateConfiguration(key, value)
	if err != nil {
		ret["err"] = err.Error()
		ret["ok"] = false
	}

	ret["data"] = res
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

	var viewInfo = getViewInfo(req, "admin")

	t := template.Must(template.New("admin.html").ParseFiles(
		"templates/admin.html", "templates/base.html"))
	err := t.ExecuteTemplate(w, "base", viewInfo)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
	}
}
