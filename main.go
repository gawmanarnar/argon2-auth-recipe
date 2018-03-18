package main

var s *WebServer

func main() {
	s = &WebServer{}

	s.Init("admin", "hunter2", "auth-recipe")

	s.Router.Get("/", Login)
	s.Router.Post("/", Register)

	s.Run()
}
