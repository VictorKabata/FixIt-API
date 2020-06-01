package seed

import (
	"log"

	"github.com/jinzhu/gorm"
	"github.com/victorkabata/FixIt/api/models"
)

//Prepopulate this data on startup
var users = []models.User{
	models.User{
		Username: "Victor Kabata",
		Email:    "victorbro14@gmail.com",
		Phone:    "0714091304",
		Password: "password",
	},
}

func Load(db *gorm.DB) {

	err := db.Debug().DropTableIfExists(&models.User{}).Error
	if err != nil {
		log.Fatalf("cannot drop table: %v", err)
	}
	err = db.Debug().AutoMigrate(&models.User{}).Error
	if err != nil {
		log.Fatalf("cannot migrate table: %v", err)
	}

	for i, _ := range users {
		err = db.Debug().Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			log.Fatalf("cannot seed users table: %v", err)
		}
	}
}
