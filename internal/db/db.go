// internal/db/db.go
package db

import (
	"golang.org/x/crypto/bcrypt"
	"log"
	"time"

	"github.com/sharif-go-lab/go-judge-platform/internal/model"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DB is the global database connection
var DB *gorm.DB

// Init connects to the database and runs migrations
func Init() {
	dsn := viper.GetString("database.dsn")
	if dsn == "" {
		log.Fatal("database.dsn is not set in config or ENV")
	}
	dbConn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	// Migrate the User model; add others as you go
	if err := dbConn.AutoMigrate(&model.User{}, &model.Question{}, &model.Submission{}); err != nil {
		log.Fatalf("failed to migrate tables: %v", err)
	}
	var count int64
	dbConn.Model(&model.User{}).Where("is_admin = ?", true).Count(&count)
	if count == 0 {
		// generate a default admin (password: "admin123")
		hash, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
		dbConn.Create(&model.User{
			Username:  "admin",
			Email:     "admin@localhost",
			Password:  string(hash),
			IsAdmin:   true,
			CreatedAt: time.Now(),
		})
		log.Println("[info] seeded initial admin user with username 'admin'")
	}

	DB = dbConn
}
