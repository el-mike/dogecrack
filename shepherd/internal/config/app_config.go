package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

var appConfig *AppConfig

// GetAppConfig - singleton implementation for app config.
func GetAppConfig() *AppConfig {
	return appConfig
}

// AppConfig - application config container.
type AppConfig struct {
	RootPath     string
	HostProvider string

	APIPort       string
	OriginAllowed string

	SSHUser     string
	SSHPassword string
	SSHDirPath  string

	PitbullImage string

	VastApiSecret string

	MongoUser     string
	MongoPassword string
	MongoHost     string
	MongoPort     string

	RedisHost string
	RedisPort string

	WalletString string

	SessionExpiration time.Duration
}

// NewAppConfig - creates new AppConfig instance and reads values from env.
func NewAppConfig(rootPath string) (*AppConfig, error) {
	if err := godotenv.Load(rootPath + "/.env"); err != nil {
		log.Fatal("Error loading .env file")
		return nil, err
	}

	config := &AppConfig{}

	config.RootPath = rootPath

	config.HostProvider = os.Getenv("HOST_PROVIDER")

	config.APIPort = os.Getenv("API_PORT")
	config.OriginAllowed = os.Getenv("ORIGIN_ALLOWED")

	config.SSHUser = os.Getenv("SSH_USER")
	config.SSHPassword = os.Getenv("SSH_PASSWORD")
	config.SSHDirPath = os.Getenv("SSH_DIR")

	config.PitbullImage = os.Getenv("PITBULL_IMAGE")

	config.VastApiSecret = os.Getenv("VAST_API_SECRET")

	config.WalletString = os.Getenv("WALLET_STRING")

	config.MongoUser = os.Getenv("MONGO_INITDB_ROOT_USERNAME")
	config.MongoPassword = os.Getenv("MONGO_INITDB_ROOT_PASSWORD")
	config.MongoHost = os.Getenv("MONGO_HOST")
	config.MongoPort = os.Getenv("MONGO_PORT")

	config.RedisHost = os.Getenv("REDIS_HOST")
	config.RedisPort = os.Getenv("REDIS_PORT")

	sessionExpirationMinutesRaw := os.Getenv("SESSION_EXPIRATION_MINUTES")

	sessionExpirationMinutes, err := strconv.Atoi(sessionExpirationMinutesRaw)
	if err != nil {
		return nil, err
	}

	config.SessionExpiration = time.Minute * time.Duration(sessionExpirationMinutes)

	appConfig = config

	return config, nil
}
