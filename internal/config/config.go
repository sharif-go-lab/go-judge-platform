package config

import (
	"log"

	"github.com/spf13/viper"
)

// Init initializes the configuration
func Init() {
	// Set default values
	viper.SetDefault("server.listen", ":8080")

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
