package views

import (
	"html/template"
	"net/http"
)

var loginTemplate *template.Template

func init() {
	var err error
	loginTemplate, err = template.ParseFiles("views/templates/login.gohtml")
	if err != nil {
		panic(err)
	}
}

// RenderLogin - renders the login page
func RenderLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	if err := loginTemplate.Execute(w, nil); err != nil {
		panic(err)
	}
}
