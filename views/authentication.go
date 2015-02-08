package views

import (
	"encoding/json"
	"github.com/gorilla/securecookie"
	"github.com/julienc91/heygo/database"
	"github.com/julienc91/heygo/globals"
	"github.com/julienc91/heygo/tools"
	"html/template"
	"net/http"
)

var cookieHandler = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32))

// Display the sign in page
func SignInHandler(w http.ResponseWriter, req *http.Request) {

	var viewInfo = GetViewInfo(req, "about")

	t := template.Must(template.New("signin.html").ParseFiles(
		"templates/signin.html", "templates/base.html"))
	err := t.ExecuteTemplate(w, "base", viewInfo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Handle a sign up request
func Signup(w http.ResponseWriter, req *http.Request) {

	var user globals.User
	err := json.Unmarshal([]byte(req.FormValue("user")), &user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var invitation globals.Invitation
	err = json.Unmarshal([]byte(req.FormValue("invitation")), &invitation)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if _, _, err := database.GetUserFromLogin(user.Login); err == nil {
		http.Error(w, "this login already exists", http.StatusConflict)
		return
	}

	invitation, err = database.GetInvitationFromValue(invitation.Value)
	if err != nil {
		http.Error(w, "this invitation is not valid", http.StatusBadRequest)
		return
	}

	if _, _, err := database.InsertUser(user, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := database.DeleteInvitationFromId(invitation.Id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var ret = map[string]interface{}{"ok": true, "err": ""}
	tools.WriteJsonResult(ret, w, http.StatusOK)
}

// Handle a login request
func Login(w http.ResponseWriter, req *http.Request) {

	var user globals.User
	err := json.Unmarshal([]byte(req.FormValue("user")), &user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	knownUser, _, err := database.GetUserFromLogin(user.Login)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	if err := database.AuthenticateUser(knownUser.Id, user.Password); err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	value := map[string]int64{"id": knownUser.Id}
	encoded, err := cookieHandler.Encode("session", value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	cookie := &http.Cookie{
		Name:  "session",
		Value: encoded,
		Path:  "/",
	}
	http.SetCookie(w, cookie)

	var ret = map[string]interface{}{"ok": true, "err": ""}
	tools.WriteJsonResult(ret, w, http.StatusOK)
}

// Handle a sign out request
func Signout(w http.ResponseWriter, req *http.Request) {

	cookie := &http.Cookie{
		Name:   "session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(w, cookie)

	var ret = map[string]interface{}{"ok": true, "err": ""}
	tools.WriteJsonResult(ret, w, http.StatusOK)
}

// Get the user's id if he is authenticated, 0 otherwise
func GetUserId(req *http.Request) int64 {
	var id int64 = 0
	if cookie, err := req.Cookie("session"); err == nil {
		cookieValue := make(map[string]int64)
		if err = cookieHandler.Decode("session", cookie.Value, &cookieValue); err == nil {
			id = cookieValue["id"]
		}
	}
	return id
}
