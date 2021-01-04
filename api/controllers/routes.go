package controllers

import "github.com/victorkabata/FixIt-API/api/middlewares"

//Initializes all the endpoints/routes.
func (s *Server) initializeRoutes() {

	// Home Route
	s.Router.HandleFunc("/", middlewares.SetMiddlewareJSON(s.Home)).Methods("GET")

	//Register Route
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
	s.Router.HandleFunc("/post/{id}", middlewares.SetMiddlewareJSON(s.GetPost)).Methods("GET")
	s.Router.HandleFunc("/post/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.UpdatePost))).Methods("PUT")
	s.Router.HandleFunc("/posts/{id}", middlewares.SetMiddlewareAuthentication(s.DeletePost)).Methods("DELETE")

	s.Router.HandleFunc("/posts/booking/{id}", middlewares.SetMiddlewareJSON(s.GetPostBooking)).Methods("GET")

	//Bookings routes
	s.Router.HandleFunc("/booking", middlewares.SetMiddlewareJSON(s.MakeBooking)).Methods("POST")
	s.Router.HandleFunc("/booking", middlewares.SetMiddlewareJSON(s.GetBookings)).Methods("GET")
	s.Router.HandleFunc("/booking/{id}", middlewares.SetMiddlewareJSON(s.UpdateBooking)).Methods("PUT")

	//Work routes
	s.Router.HandleFunc("/work", middlewares.SetMiddlewareJSON(s.CreateWork)).Methods("POST")
	s.Router.HandleFunc("/work/{id}", middlewares.SetMiddlewareJSON(s.GetWork)).Methods("GET")

	//Review routes
	s.Router.HandleFunc("/review", middlewares.SetMiddlewareJSON(s.GetReviews)).Methods("GET")
	s.Router.HandleFunc("/review", middlewares.SetMiddlewareJSON(s.CreateReview)).Methods("POST")
	s.Router.HandleFunc("/review/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.UpdateReview))).Methods("PUT")
	s.Router.HandleFunc("/review/{id}", middlewares.SetMiddlewareAuthentication(s.DeleteReview)).Methods("DELETE")
}
