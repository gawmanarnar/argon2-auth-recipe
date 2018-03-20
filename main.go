package main

var s *WebServer

func main() {
	s = &WebServer{}

	s.SetupDB("admin", "hunter2", "auth-recipe")
	s.SetupRoutes()

	s.Run()
}
