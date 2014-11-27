package views

import (
	"github.com/gorilla/securecookie"
	"heygo/database"
	"html/template"
	"net/http"
)

var cookieHandler = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32))

// Display the sign in page
func SignInHandler(w http.ResponseWriter, req *http.Request) {

	var viewInfo = getViewInfo(req, "about")

	t := template.Must(template.New("signin.html").ParseFiles(
		"templates/signin.html", "templates/base.html"))
	err := t.ExecuteTemplate(w, "base", viewInfo)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
	}
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}

// Handle a sign up request
func SignupHandler(w http.ResponseWriter, req *http.Request) {

	var data = make(map[string]interface{})
	data["login"] = req.FormValue("login")
	data["password"] = req.FormValue("password")
	password2 := req.FormValue("password2")
	invitation := req.FormValue("invitation")

	if data["password"] != password2 {
		http.Error(w, "invalid pasword, at least 7 characters", http.StatusBadRequest)
		return
	}

	if _, err := database.PrepareGetFromKey("login", data["login"], database.TableUsers); err == nil {
		http.Error(w, "this login already exists", http.StatusConflict)
		return
	}

	if _, err := database.PrepareGetFromKey("value", invitation, database.TableInvitations); err != nil {
		http.Error(w, "this invitation is not valid", http.StatusBadRequest)
		return
	}

	if _, err := database.PrepareInsert(data, database.TableUsers); err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	if err := database.PrepareDeleteFromKey("value", invitation, database.TableInvitations); err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	writeResponse("", w, "text/plain", http.StatusOK)
}

// Handle a login request
func LoginHandler(w http.ResponseWriter, req *http.Request) {

	var login = req.FormValue("login")
	var password = req.FormValue("password")
	var err error

	user, err := database.PrepareGetFromKey("login", login, database.TableUsers)
	if err != nil {
		http.Error(w, "Invalid login", http.StatusUnauthorized)
		return
	}
	if err := database.AuthenticateUser(user["id"].(int64), password); err != nil {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	value := map[string]int64{"id": user["id"].(int64)}
	encoded, err := cookieHandler.Encode("session", value)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	cookie := &http.Cookie{
		Name:  "session",
		Value: encoded,
		Path:  "/",
	}
	http.SetCookie(w, cookie)
	http.Redirect(w, req, "/about", http.StatusFound)
}

// Handle a sign out request
func SignoutHandler(w http.ResponseWriter, req *http.Request) {

	cookie := &http.Cookie{
		Name:   "session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(w, cookie)
	http.Redirect(w, req, "/", http.StatusFound)
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
