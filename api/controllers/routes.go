package controllers

import "github.com/victorkabata/FixIt/api/middlewares"

//Initializes all the endpoints/routes.
func (s *Server) initializeRoutes() {

	// Home Route
	s.Router.HandleFunc("/", middlewares.SetMiddlewareJSON(s.Home)).Methods("GET")

	// Login Route
	s.Router.HandleFunc("/login", middlewares.SetMiddlewareJSON(s.Login)).Methods("POST")

	//Users routes
	s.Router.HandleFunc("/users", middlewares.SetMiddlewareJSON(s.CreateUser)).Methods("POST")
	s.Router.HandleFunc("/users", middlewares.SetMiddlewareJSON(s.GetUsers)).Methods("GET")
	s.Router.HandleFunc("/users/{id}", middlewares.SetMiddlewareJSON(s.GetUser)).Methods("GET")
	s.Router.HandleFunc("/users/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.UpdateUser))).Methods("PUT")
	s.Router.HandleFunc("/users/{id}", middlewares.SetMiddlewareAuthentication(s.DeleteUser)).Methods("DELETE")

	//Upload routes
	// s.Router.HandleFunc("/uploads", middlewares.SetMiddlewareJSON(s.CreateUpload)).Methods("POST")
	// s.Router.HandleFunc("/uploads", middlewares.SetMiddlewareJSON(s.GetUploads)).Methods("GET")
	// s.Router.HandleFunc("/uploads/{id}", middlewares.SetMiddlewareJSON(s.GetUpload)).Methods("GET")
	// s.Router.HandleFunc("/uploads/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.UpdateUpload))).Methods("PUT")
	// s.Router.HandleFunc("/uploads/{id}", middlewares.SetMiddlewareAuthentication(s.DeleteUpload)).Methods("DELETE")
}
