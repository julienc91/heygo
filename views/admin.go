package views

import (
    "net/http"
    "html/template"
    "gomet/database"
    "strconv"
    "fmt"
)

type AdminData struct {
    Users []map[string]interface{}
    Invitations []map[string]interface{}
    Groups []map[string]interface{}
    Videos []map[string]interface{}
    VideoGroups []map[string]interface{}
    Membership map[int][]int
    Classification map[int][]int
    Permissions map[int][]int
}


func AdminUpdateHandler(w http.ResponseWriter, req *http.Request) {

    if RedirectIfNotAdmin(w, req) {
        return
    }

    req.ParseForm()
    var data = make(map[string]interface{})

    for k, v := range req.Form {
        if k != "table" && k != "id" {
            data[k] = v[0]
        }
    }

    var table = req.Form["table"][0]
    var idStr = req.Form["id"][0]

    id, err := strconv.Atoi(idStr)
    if err != nil {
        return
    }

    var fs = map[string]func (int, map[string]interface{}) error{
        "users": database.UpdateUser,
        "invitations": database.UpdateInvitation,
        "videos": database.UpdateVideo,
        "groups": database.UpdateGroup,
        "video_groups": database.UpdateVideoGroup }

    f, ok := fs[table]
    if !ok {
        return
    }

    err = f(id, data)
    if err != nil {
        fmt.Fprintf(w, "{\"ok\": false, \"err\": \"%s\"}", err.Error())
        return
    }
    fmt.Fprintf(w, "{\"ok\": true}")
}


func AdminDeleteHandler(w http.ResponseWriter, req *http.Request) {
    
    if RedirectIfNotAdmin(w, req) {
        return 
    }

    req.ParseForm()

    id, err := strconv.Atoi(req.Form["id"][0])
    if err != nil {
        return
    }

    var fs = map[string]func (int, string) error {
        "users": database.DeleteRow,
        "invitations": database.DeleteRow,
        "videos": database.DeleteRow,
        "groups": database.DeleteRow,
        "video_groups": database.DeleteRow }

    f, ok := fs[req.Form["table"][0]]
    if !ok {
        return
    }

    err = f(id, req.Form["table"][0])
    if err != nil {
        fmt.Fprintf(w, "{\"ok\": false, \"err\": \"%s\"}", err.Error())
        return
    }
    fmt.Fprintf(w, "{\"ok\": true}")   
}


func AdminDeletePivotHandler(w http.ResponseWriter, req *http.Request) {
    
    if RedirectIfNotAdmin(w, req) {
        return 
    }

    req.ParseForm()

    var table = req.Form["table"][0]
    var k1, k2, rtable string
    
    switch table {
    case "membership":
        k1 = "users_id"
        k2 = "groups_id"
        rtable = "membership"      
    case "classification":
        k1 = "videos_id"
        k2 = "video_groups_id"
        rtable = "video_classification"
    case "permissions":
        k1 = "video_groups_id"
        k2 = "groups_id"
        rtable = "video_permissions"
    default:
        return
    }

    i, err := strconv.Atoi(req.Form[k1][0])
    if err != nil {
        return
    }
    
    j, err := strconv.Atoi(req.Form[k2][0])
    if err != nil {
        return
    }
    
    err = database.DeletePivotTableRow(i, j, k1, k2, rtable)

    if err != nil {
        fmt.Fprintf(w, "{\"ok\": false, \"err\": \"%s\"}", err.Error())
        return
    }
    fmt.Fprintf(w, "{\"ok\": true}")   
}


func AdminInsertHandler(w http.ResponseWriter, req *http.Request) {

    if RedirectIfNotAdmin(w, req) {
        return
    }

    req.ParseForm()

    var data = make(map[string]interface{})

    for k, v := range req.Form {
        if k != "table" {
            data[k] = v[0]
        }
    }

    var table = req.Form["table"][0]  
    delete(req.Form, "table")

    var fs = map[string]func(map[string]interface{}) (int, error) {
        "users": database.InsertUser,
        "invitations": database.InsertInvitation,
        "videos": database.InsertVideo,
        "groups": database.InsertGroup,
        "video_groups": database.InsertVideoGroup,
        "membership": database.InsertMembership,
        "classification": database.InsertClassification,
        "permissions": database.InsertPermission }

    f, ok := fs[table]
    if !ok {
        return
    }
    
    id, err := f(data)
    if err != nil {
        fmt.Fprintf(w, "{\"ok\": false, \"err\": \"%s\"}", err.Error())
        return
    }
    fmt.Fprintf(w, "{\"ok\": true, \"id\": %d}", id)
}


func AdminHandler(w http.ResponseWriter, req *http.Request) {

    if RedirectIfNotAdmin(w, req) {
        return
    }

    users, err := database.GetAll("users")
    if err != nil {
        fmt.Println(err)
        return
    }

    invitations, err := database.GetAll("invitations")
    if err != nil {
        fmt.Println(err)
        return
    }

    videos, err := database.GetAll("videos")
    if err != nil {
        fmt.Println(err)
        return
    }

    groups, err := database.GetAll("groups")
    if err != nil {
        fmt.Println(err)
        return
    }

    videoGroups, err := database.GetAll("video_groups")
    if err != nil {
        fmt.Println(err)
        return
    }

    membership, err := database.GetAllPivotTable("membership", "groups_id", "users_id")
    if err != nil {
        fmt.Println(err)
        return
    }

    classification, err := database.GetAllPivotTable("video_classification", "video_groups_id", "videos_id")
    if err != nil {
        fmt.Println(err)
        return
    }

    permissions, err := database.GetAllPivotTable("video_permissions", "groups_id", "video_groups_id")
    if err != nil {
        fmt.Println(err)
        return
    }

    t := template.Must(template.New("admin.html").ParseFiles(
        "templates/admin.html", "templates/base.html"))
    err = t.ExecuteTemplate(w, "base", AdminData{users, invitations, groups, videos, videoGroups, membership, classification, permissions})
    if err != nil {
        panic(err)
    }

}
