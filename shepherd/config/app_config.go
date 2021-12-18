package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

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

	appConfig := &AppConfig{}

	appConfig.SSHUser = os.Getenv("SSH_USER")
	appConfig.SSHPassword = os.Getenv("SSH_PASSWORD")
	appConfig.SSHDirPath = os.Getenv("SSH_DIR")

	appConfig.VastApiSecret = os.Getenv("VAST_API_SECRET")

	appConfig.WalletString = os.Getenv("WALLET_STRING")

	appConfig.MongoUser = os.Getenv("MONGO_INITDB_ROOT_USERNAME")
	appConfig.MongoPassword = os.Getenv("MONGO_INITDB_ROOT_PASSWORD")
	appConfig.MongoHost = os.Getenv("MONGO_HOST")
	appConfig.MongoPort = os.Getenv("MONGO_PORT")

	return appConfig, nil
}
