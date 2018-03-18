package main

var s *WebServer

func main() {
	s = &WebServer{}

	s.Init("admin", "hunter2", "auth-recipe")

	s.Router.Post("/login", Login)
	s.Router.Post("/register", Register)

	s.Run()
}
