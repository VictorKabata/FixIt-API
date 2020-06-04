package api

import (
	"fmt"
	"log"
	"os"

	"github.com/victorkabata/FixIt/api/controllers"
)

var server = controllers.Server{}

//Initializing the server connection.
func Run() {

	var port = os.Getenv("PORT")

	var DB_HOST = "us-cdbr-east-05.cleardb.net"
	var DB_DRIVER = "mysql"
	var DB_USER = "b4349f229cb9d8"
	var DB_PASSWORD = "786bb4c4"
	var DB_NAME = "heroku_f35f1129c94b864"

	var err error
	//err = godotenv.Load()
	if err != nil {
		log.Fatalf("Error fetching env, not coming through %v", err)
	} else {
		fmt.Println("Fetching the env values")
	}

	server.Initialize(DB_DRIVER, DB_USER, DB_PASSWORD, DB_HOST, DB_NAME)

	server.Run(":" + port) //Port for listening and serving requests.

}
