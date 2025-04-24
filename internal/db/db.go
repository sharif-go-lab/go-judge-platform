// internal/db/db.go
package db

import (
	"errors"
	"log"
	"time"

	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/sharif-go-lab/go-judge-platform/internal/model"
)

// DB is the global database connection
var DB *gorm.DB

// Init connects to the database, runs migrations, and seeds a default admin
func Init() {
	// Read DSN
	dsn := viper.GetString("database.dsn")
	if dsn == "" {
		log.Fatal("database.dsn is not set in config or ENV")
	}

	// Open GORM connection
	dbConn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// Auto-migrate all Phase 3 models
	if err := dbConn.AutoMigrate(
		&model.User{},
		&model.Problem{},
		&model.Submission{},
		&model.Session{},
		&model.TestCase{},
	); err != nil {
		log.Fatalf("failed to migrate tables: %v", err)
	}

	// Seed default admin user if none exists
	var count int64
	dbConn.Model(&model.User{}).
		Where("is_admin = ?", true).
		Count(&count)

	if count == 0 {
		hash, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
		if err != nil {
			log.Fatalf("failed to hash default admin password: %v", err)
		}

		admin := model.User{
			Username:     "admin",
			Email:        "admin@localhost",
			PasswordHash: string(hash),
			IsAdmin:      true,
			CreatedAt:    time.Now(),
		}

		// Try to find existing user by username
		var existing model.User
		switch err := dbConn.Where("username = ?", admin.Username).First(&existing).Error; {
		case errors.Is(err, gorm.ErrRecordNotFound):
			// Create new admin
			if err := dbConn.Create(&admin).Error; err != nil {
				log.Fatalf("failed to seed default admin: %v", err)
			}
			log.Println("[info] seeded initial admin user with username 'admin'")

		case err == nil:
			// Promote existing to admin
			existing.IsAdmin = true
			if err := dbConn.Save(&existing).Error; err != nil {
				log.Fatalf("failed to update existing admin: %v", err)
			}

		default:
			log.Fatalf("error checking existing admin: %v", err)
		}
	}

	// Assign global
	DB = dbConn
}
