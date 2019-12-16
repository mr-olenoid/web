package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

type pData struct {
	Name string
}

var (
	key   = []byte("randomly-secure-key")
	store = sessions.NewCookieStore(key)
)

func secret(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie-name")

	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	fmt.Fprintf(w, "ok!")
}

func login(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie-name")

	// Authentication goes here
	// ...
	fmt.Println("logining")

	// Set user as authenticated
	session.Values["authenticated"] = true
	session.Save(r, w)
}

func logout(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie-name")

	// Revoke users authentication
	session.Values["authenticated"] = false
	session.Save(r, w)
}

func main() {
	wwwroot := os.Getenv("wwwroot")

	router := mux.NewRouter()
	router.HandleFunc("/hello", handler).Methods("GET")
	router.HandleFunc("/", loginPage).Methods("GET")
	router.HandleFunc("/", login).Methods("POST")
	router.HandleFunc("/secret", secret).Methods("GET")

	staticFileDirectory := http.Dir(wwwroot)
	staticFileHandler := http.StripPrefix(wwwroot + "assets/", http.FileServer(staticFileDirectory))
	router.PathPrefix(wwwroot + "assets/").Handler(staticFileHandler).Methods("GET")

	http.ListenAndServe(":8080", router)
}

func handler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles(wwwroot + "www/base.html", wwwroot + "www/index.html"))
	data := pData{Name: "Vasya"}
	//tmpl.ExecuteTemplate(w, "layout", data)
	tmpl.Execute(w, data)

	//fmt.Fprintf(w, "Hello World!")
}

func loginPage(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles(wwwroot + "www/login.html"))
	tmpl.Execute(w, "")
}
