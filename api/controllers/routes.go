package controllers

import "github.com/victorkabata/FixIt-API/api/middlewares"

//Initializes all the endpoints/routes.
func (s *Server) initializeRoutes() {

	// Home Route
	s.Router.HandleFunc("/", middlewares.SetMiddlewareJSON(s.Home)).Methods("GET")

	//Login Route
	s.Router.HandleFunc("/register", middlewares.SetMiddlewareJSON(s.CreateUser)).Methods("POST")

	// Login Route
	s.Router.HandleFunc("/login", middlewares.SetMiddlewareJSON(s.Login)).Methods("POST")

	//Upload profile pic
	s.Router.HandleFunc("/profile", middlewares.SetMiddlewareJSON(UploadProfilePic)).Methods("POST")

	//Users routes
	s.Router.HandleFunc("/users", middlewares.SetMiddlewareJSON(s.GetUsers)).Methods("GET")
	s.Router.HandleFunc("/users/{id}", middlewares.SetMiddlewareJSON(s.GetUser)).Methods("GET")
	s.Router.HandleFunc("/users/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.UpdateUser))).Methods("PUT")
	s.Router.HandleFunc("/users/{id}", middlewares.SetMiddlewareAuthentication(s.DeleteUser)).Methods("DELETE")

	//Upload profile pic
	s.Router.HandleFunc("/postpic", middlewares.SetMiddlewareJSON(UploadPostPic)).Methods("POST")

	//Post routes
	s.Router.HandleFunc("/posts", middlewares.SetMiddlewareJSON(s.CreatePost)).Methods("POST")
	s.Router.HandleFunc("/posts", middlewares.SetMiddlewareJSON(s.GetPosts)).Methods("GET")
	// s.Router.HandleFunc("/post/{id}", middlewares.SetMiddlewareJSON(s.GetUpload)).Methods("GET")
	s.Router.HandleFunc("/post/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.UpdatePost))).Methods("PUT")
	s.Router.HandleFunc("/posts/{id}", middlewares.SetMiddlewareAuthentication(s.DeletePost)).Methods("DELETE")
}
