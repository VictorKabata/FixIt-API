package api

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/victorkabata/FixIt-API/api/controllers"
)

var server = controllers.Server{}

//Initializing the server connection.
func Run() {

	//var port = os.Getenv("PORT")

	var err error
	err = godotenv.Load()
	if err != nil {
		log.Fatalf("Error fetching env, not coming through %v", err)
	} else {
		fmt.Println("Fetching the env values")
	}

	server.Initialize(os.Getenv("DB_DRIVER"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_NAME"))

	server.Run(":" + "8081") //Port for listening and serving requests.

}
