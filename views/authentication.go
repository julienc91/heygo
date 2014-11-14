package views

import (
	"github.com/gorilla/securecookie"
	"gomet/database"
	"html/template"
	"net/http"
)

var cookieHandler = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32))

func SignInHandler(w http.ResponseWriter, req *http.Request) {

	t := template.Must(template.New("signin.html").ParseFiles(
		"templates/signin.html"))
	err := t.Execute(w, nil)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}

func SignupHandler(w http.ResponseWriter, req *http.Request) {

	login := req.FormValue("login")
	password := req.FormValue("password")
	password2 := req.FormValue("password2")
	invitation := req.FormValue("invitation")

	if len(password) < 7 || password != password2 {
		http.Error(w, "Invalid pasword, at least 7 characters", http.StatusBadRequest)
		return
	}

	id, err := database.GetUserIdFromLogin(login)
	if err == nil && id != 0 {
		http.Error(w, "This login already exists", http.StatusConflict)
		return
	}

	ok, err := database.CheckInvitation(invitation)
	if err != nil || !ok {
		http.Error(w, "This invitation is not valid", http.StatusBadRequest)
		return
	}

	err = database.AddUser(login, password, invitation)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	http.Error(w, "", http.StatusOK)
}

func LoginHandler(w http.ResponseWriter, req *http.Request) {

	login := req.FormValue("login")
	password := req.FormValue("password")
	var id int64
	var err error

	if login != "" && password != "" {
		id, err = database.GetUserIdFromLogin(login)
		if err != nil {
			http.Error(w, "Invalid login", http.StatusUnauthorized)
			return
		}
		if err := database.AuthenticateUser(id, password); err != nil {
			http.Error(w, "Invalid password", http.StatusUnauthorized)
			return
		}
	}

	value := map[string]int64{"id": id}
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
	http.Redirect(w, req, "/about", 302)
}

func SignoutHandler(w http.ResponseWriter, req *http.Request) {

	cookie := &http.Cookie{
		Name:   "session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(w, cookie)
	http.Redirect(w, req, "/", 302)
}

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

func MainPageHandler(w http.ResponseWriter, req *http.Request) {
	if RedirectIfNotAuthenticated(w, req) {
		return
	}
	http.Redirect(w, req, "/about", http.StatusFound)
}
