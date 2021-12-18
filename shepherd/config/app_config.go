package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var appConfig *AppConfig

// GetAppConfig - singleton implementation for app config.
func GetAppConfig() *AppConfig {
	return appConfig
}

// AppConfig - application config container.
type AppConfig struct {
	SSHUser     string
	SSHPassword string
	SSHDirPath  string

	VastApiSecret string

	MongoUser     string
	MongoPassword string
	MongoHost     string
	MongoPort     string

	WalletString string
}

// NewAppConfig - creates new AppConfig instance and reads values from env.
func NewAppConfig() (*AppConfig, error) {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
		return nil, err
	}

	config := &AppConfig{}

	config.SSHUser = os.Getenv("SSH_USER")
	config.SSHPassword = os.Getenv("SSH_PASSWORD")
	config.SSHDirPath = os.Getenv("SSH_DIR")

	config.VastApiSecret = os.Getenv("VAST_API_SECRET")

	config.WalletString = os.Getenv("WALLET_STRING")

	config.MongoUser = os.Getenv("MONGO_INITDB_ROOT_USERNAME")
	config.MongoPassword = os.Getenv("MONGO_INITDB_ROOT_PASSWORD")
	config.MongoHost = os.Getenv("MONGO_HOST")
	config.MongoPort = os.Getenv("MONGO_PORT")

	appConfig = config

	return config, nil
}
