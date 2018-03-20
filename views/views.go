package views

import (
	"html/template"
	"net/http"

	"github.com/gorilla/csrf"
)

var loginTemplate *template.Template
var signupTemplate *template.Template

func init() {
	var err error
	loginTemplate, err = template.ParseFiles("views/templates/login.gohtml")
	if err != nil {
		panic(err)
	}
	signupTemplate, err = template.ParseFiles("views/templates/signup.gohtml")
	if err != nil {
		panic(err)
	}
}

// RenderLogin - renders the login page
func RenderLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	if err := loginTemplate.Execute(w, map[string]interface{}{
		csrf.TemplateTag: csrf.TemplateField(r),
	}); err != nil {
		panic(err)
	}
}

// RenderSignup - renders the signup page
func RenderSignup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	if err := signupTemplate.Execute(w, map[string]interface{}{
		csrf.TemplateTag: csrf.TemplateField(r),
	}); err != nil {
		panic(err)
	}
}
