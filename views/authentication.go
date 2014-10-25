package views

import (
	"net/http"
	"html/template"
    "github.com/gorilla/securecookie"
    "gomet/database"
    "fmt"
)

var cookieHandler = securecookie.New(
    securecookie.GenerateRandomKey(64),
    securecookie.GenerateRandomKey(32))

func SignInHandler(w http.ResponseWriter, req *http.Request) {

	t := template.Must(template.New("signin.html").ParseFiles(
        "templates/signin.html"))
	err := t.Execute(w, nil)
	if err != nil {
		w.WriteHeader(500)
        fmt.Fprintln(w, err)
	}
}


func SignupHandler(w http.ResponseWriter, req *http.Request) {

    login := req.FormValue("login")
    password := req.FormValue("password")
    password2 := req.FormValue("password2")
    invitation := req.FormValue("invitation")

    if len(password) < 7 || password != password2 {
        fmt.Fprintln(w, "Mot de passe invalide (minimum 7 caractères)")
        return
    }
    
    id, err := database.GetUserIdFromLogin(login)
    if err == nil && id != 0 {
        fmt.Fprintln(w, "Login déjà existant")
        return
    }

    ok, err := database.CheckInvitation(invitation)
    if err != nil || !ok {
        fmt.Fprintln(w, "Mauvaise invitation")
        return
    }

    err = database.AddUser(login, password, invitation)
    if err != nil {
        fmt.Fprintln(w, err)
        return
    }

    fmt.Fprintln(w, "Utilisateur créé avec succès! Vous pouvez à présent vous connecter")
}

func LoginHandler(w http.ResponseWriter, req *http.Request) {

    login := req.FormValue("login")
    password := req.FormValue("password")
    var id int
    var err error
    
    if login != "" && password != "" {
        id, err = database.GetUserIdFromLogin(login)
        if err != nil {
            w.WriteHeader(403)
            fmt.Fprintln(w, "Invalid login")
            return
        }
        if err := database.AuthenticateUser(id, password); err != nil {
            w.WriteHeader(403)
            fmt.Fprintln(w, "Invalid password")
            return
        }
    }

    value := map[string]int{"id": id}
    encoded, err := cookieHandler.Encode("session", value)
    if err != nil {
        w.WriteHeader(500)
        fmt.Fprintln(w, err)
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
        Name:  "session",
        Value: "",
        Path:  "/",
        MaxAge: -1,
    }
    http.SetCookie(w, cookie)
    http.Redirect(w, req, "/", 302)
}


func GetUserId(req *http.Request) int {
    var id = 0
    if cookie, err := req.Cookie("session"); err == nil {
        cookieValue := make(map[string]int)
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
    http.Redirect(w, req, "/about", 302)
}

