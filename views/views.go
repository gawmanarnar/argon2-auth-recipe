package views

import (
	"html/template"
	"net/http"
)

var loginTemplate *template.Template

// Compile - compiles templates. This function needs to be called
// before trying to execute any templates
func Compile() {
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
