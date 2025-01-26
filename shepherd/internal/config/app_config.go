package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
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

	APIPort        string
	OriginsAllowed []string

	SSHUser       string
	SSHPassword   string
	SSHDirPath    string
	SSHPrivateKey string

	PitbullImage string

	VastApiSecret string

	MongoConnectionString string
	RedisConnectionString string

	WalletString string

	SessionExpiration time.Duration
}

// NewAppConfig - creates new AppConfig instance and reads values from env.
func NewAppConfig(rootPath string) (*AppConfig, error) {
	// We only want to load .env outside of prod, as prod will have env variables set explicitly.
	if strings.ToLower(os.Getenv("APP_ENV")) != "prod" {
		dotenvFileName := ".env"

		// Allows for running things locally against prod databases.
		if strings.ToLower(os.Getenv("USE_PROD_CONFIG")) == "true" {
			dotenvFileName = ".env.prod"
		}

		if err := godotenv.Load(rootPath + fmt.Sprintf("/%s", dotenvFileName)); err != nil {
			log.Fatal("Error loading .env file")
			return nil, err
		}
	}

	config := &AppConfig{}

	config.RootPath = rootPath

	config.HostProvider = os.Getenv("HOST_PROVIDER")

	config.APIPort = os.Getenv("API_PORT")

	originsAllowedRaw := os.Getenv("ORIGINS_ALLOWED")
	config.OriginsAllowed = strings.Split(originsAllowedRaw, ",")

	config.SSHUser = os.Getenv("SSH_USER")
	config.SSHPassword = os.Getenv("SSH_PASSWORD")
	config.SSHDirPath = os.Getenv("SSH_DIR")
	config.SSHPrivateKey = os.Getenv("SSH_PRIVATE_KEY")

	config.PitbullImage = os.Getenv("PITBULL_IMAGE")

	config.VastApiSecret = os.Getenv("VAST_API_SECRET")

	config.WalletString = os.Getenv("WALLET_STRING")

	config.MongoConnectionString = os.Getenv("MONGO_CONNECTION_STRING")
	config.RedisConnectionString = os.Getenv("REDIS_CONNECTION_STRING")

	sessionExpirationMinutesRaw := os.Getenv("SESSION_EXPIRATION_MINUTES")

	sessionExpirationMinutes, err := strconv.Atoi(sessionExpirationMinutesRaw)
	if err != nil {
		return nil, err
	}

	config.SessionExpiration = time.Minute * time.Duration(sessionExpirationMinutes)

	appConfig = config

	return config, nil
}
