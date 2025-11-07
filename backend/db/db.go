package db

import (
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
	userCLient.Db = DB

	log.Info("Finishing Migration Database Tables")
}

func StartDbEngine() {
	// Migrating User and VerificationToken models.
	if err := DB.AutoMigrate(&model.UserModel{}, &model.VerificationToken{}); err != nil {
		panic(fmt.Sprintf("Error creating tables: %v", err))
	}
	log.Info("Database tables migrated successfully")

	// Seed default admin user
	seedAdminUser()
}

// seedAdminUser creates the default admin user if it doesn't exist
func seedAdminUser() {
	var count int64
	DB.Model(&model.UserModel{}).Where("email = ?", "admin@unichat.com").Count(&count)

	if count == 0 {
		// Use raw SQL to avoid time.Time zero value issues with nullable datetime fields
		result := DB.Exec(`
			INSERT INTO user_models (email, password_hash, first_name, last_name, is_admin, is_verified, created_at, verification_code, code_expires_at)
			VALUES (?, ?, ?, ?, ?, ?, NOW(), NULL, NULL)
		`, "admin@unichat.com", "240be518fabd2724ddb6f04eeb1da5967448d7e831c08c8fa822809f74c720a9", "Admin", "UniChat", true, true)

		if result.Error != nil {
			log.Warnf("Failed to create admin user: %v", result.Error)
		} else {
			log.Info("✓ Default admin user created (admin@unichat.com / admin123)")
		}
	} else {
		log.Info("✓ Admin user already exists")
	}
}
