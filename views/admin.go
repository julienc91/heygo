package views

import (
    "net/http"
    "html/template"
    "gomet/database"
    "gomet/tools"
    "github.com/gorilla/mux"
    "strconv"
)

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

    var fs = map[string]func (int64, map[string]interface{}) (map[string]interface{}, error) {
        "users": database.UpdateUser,
        "invitations": database.UpdateInvitation,
        "videos": database.UpdateVideo,
        "groups": database.UpdateGroup,
        "video_groups": database.UpdateVideoGroup }

    var ret = map[string]interface{}{"ok": true, "err": ""}
    
    updated, err := fs[table](id, data)
    if err != nil {
        ret["err"] = err.Error()
        ret["ok"] = false
    }
    
    ret["data"] = updated
    writeJsonResult(ret, w, http.StatusOK)
}


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
    err = database.DeleteRow(id, table)
    if err != nil {
        ret["err"] = err.Error()
        ret["ok"] = false
    }

    writeJsonResult(ret, w, http.StatusOK)
}


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

    var fs = map[string]func(map[string]interface{}) (map[string]interface{}, error) {
        "users": database.InsertUser,
        "invitations": database.InsertInvitation,
        "videos": database.InsertVideo,
        "groups": database.InsertGroup,
        "video_groups": database.InsertVideoGroup }

    var ret = map[string]interface{}{"ok": true, "err": ""}
    
    inserted, err := fs[table](data)
    if err != nil {
        ret["err"] = err.Error()
        ret["ok"] = false
    }
    
    ret["data"] = inserted
    writeJsonResult(ret, w, http.StatusOK)
}


func AdminGetHandler(w http.ResponseWriter, req *http.Request) {

    if RedirectIfNotAdmin(w, req) {
        return
    }

    params := mux.Vars(req)

    table := params["table"]
    if !tools.InArray(database.MainTables, table) {
        http.Error(w, "table is not valid", http.StatusBadRequest)
        return
    }

    var ret = map[string]interface{}{"ok": true, "err": ""}
    
    rows, err := database.GetAll(table)
    if err != nil {
        ret["err"] = err.Error()
        ret["ok"] = false
    }

    ret["data"] = rows
    writeJsonResult(ret, w, http.StatusOK)
}


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
    
    rows, err := database.GetFromId(table, id)
    if err != nil {
        ret["err"] = err.Error()
        ret["ok"] = false
    }

    ret["data"] = rows
    writeJsonResult(ret, w, http.StatusOK)
}


func AdminMediaCheckHandler(w http.ResponseWriter, req *http.Request) {

    if RedirectIfNotAdmin(w, req) {
        return
    }

    path := req.FormValue("path")
    
    var ret = map[string]interface{}{"ok": tools.CheckFilePath(path)}
    writeJsonResult(ret, w, http.StatusOK)
}


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
