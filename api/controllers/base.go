package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/victorkabata/FixIt-API/api/models"
)

type Server struct {
	DB     *gorm.DB
	Router *mux.Router
}

//Initializes the database connection and mux routers
func (server *Server) Initialize(Dbdriver, DbUser, DbPassword, DbHost, DbName string) {

	var err error

	if Dbdriver == "mysql" {
		DBURL := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local", DbUser, DbPassword, DbHost, DbName)
		server.DB, err = gorm.Open(Dbdriver, DBURL)
		if err != nil {
			fmt.Printf("Cannot connect to %s database\n", Dbdriver)
			log.Fatal("Error:", err)
		} else {
			fmt.Printf("Connected successfully to the %s database\n", Dbdriver)
		}
	}

	if Dbdriver == "postgres" {
		DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", DbHost, DbUser, DbName, DbPassword)
		server.DB, err = gorm.Open(Dbdriver, DBURL)
		if err != nil {
			fmt.Printf("Cannot connect to %s database\n", Dbdriver)
			log.Fatal("Error:", err)
		} else {
			fmt.Printf("Connected successfully to the %s database\n", Dbdriver)
		}
	}

	server.DB.Debug().AutoMigrate(&models.User{}, &models.Post{}, &models.Booking{}, &models.Work{}) //database migration

	server.Router = mux.NewRouter()

	server.initializeRoutes()
}

//Set listening port
func (server *Server) Run(addr string) {
	fmt.Println("Listening to port" + addr)
	log.Fatal(http.ListenAndServe(addr, server.Router))
}
