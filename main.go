package main

var s *WebServer

func main() {
	s = &WebServer{}

	s.SetupDB()
	s.SetupRoutes()

	s.Run()
}
