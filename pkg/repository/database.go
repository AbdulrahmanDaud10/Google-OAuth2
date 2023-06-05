package repository

import (
	"fmt"
	"log"

	"github.com/AbdulrahmanDaud10/google-0auth2/pkg/api"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// ConnectDB function that will utilize GORM to establish a connection with the PSQL database and perform automatic migration of the user model.
func PostgresDatabaseConnection() {
	var err error
	dsn := "host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Shanghai"
	DB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect to the Database")
	}

	// the AutoMigrate() function will be used to synchronize the database schema with the GORM model defined in the api/model.go file.
	DB.AutoMigrate(&api.User{})
	fmt.Println("🚀 Connected Successfully to the Database")
}
