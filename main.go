package main

import "github.com/gimmeasandwich/argon2-auth-recipe/views"

var s *WebServer

func main() {
	s = &WebServer{}

	s.Init("admin", "hunter2", "auth-recipe")

	s.Router.Post("/login", Login)
	s.Router.Post("/register", Register)

	s.Router.Get("/login", views.RenderLogin)

	s.Run()
}
