package db

import (
	activityClient "backend/clients/activity"
	inscriptionClient "backend/clients/inscription"
	userCLient "backend/clients/user"
	"backend/model"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	DB  *gorm.DB
	err error
)

func init() {

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASS")
	name := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, password, host, port, name,
	)

	//dsn := "root:FranMySql1@@tcp(127.0.0.1:3306)/arqui_software?charset=utf8mb4&parseTime=True&loc=Local"
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Info("Connection Failed to Open")
		log.Fatal(err)
	} else {
		log.Info("Connection Established")
	}
	activityClient.Db = DB
	userCLient.Db = DB
	inscriptionClient.Db = DB

	log.Info("Finishing Migration Database Tables")
}

func StartDbEngine() {
	// Migrating all classes model.
	if err := DB.AutoMigrate(&model.ActivityModel{}); err != nil {
		panic(fmt.Sprintf("Error creating table: %v", err))
	}
	if err := DB.AutoMigrate(&model.UserModel{}); err != nil {
		panic(fmt.Sprintf("Error creating table: %v", err))
	}
	if err := DB.AutoMigrate(&model.InscriptionModel{}); err != nil {
		panic(fmt.Sprintf("Error creating table: %v", err))
	}
}
