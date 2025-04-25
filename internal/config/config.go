package config

import (
	"log"
	"strings"

	"github.com/spf13/viper"
)

// Init initializes the configuration
func Init() {
	// Set default values
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.SetDefault("server.listen", ":8080")
	viper.SetDefault("code_runner.listen", ":9000")
	viper.SetDefault("database.dsn", "postgres://user:pass@localhost:5432/go_judge?sslmode=disable")
	viper.SetDefault("session.secret", "super-secret-key")
	// Look for config file
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")

	// Read config file (if exists)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("No config file found, using defaults")
		} else {
			log.Printf("Error reading config file: %v", err)
		}
	}
}
