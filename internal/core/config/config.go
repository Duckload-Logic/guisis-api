package config

import (
	"os"
	"strconv"
)

type Config struct {
	DBHost  string
	DBPort  string
	DBUser  string
	DBPass  string
	DBName  string
	DBTLS   bool
	DBTLSCA string

	JWTSecret string

	WebsitesPort string

	LocalUploadDIR               string
	AzureStorageConnectionString string
	AzureContainerName           string

	IsProduction bool

	IDPClientID     string
	IDPClientSecret string
	IDPBaseUrl      string

	RedisHost string
	RedisPort string
	RedisPass string
	RedisDB   int

	GotenbergURL string

	SendGridAPIKey string

	MailPitHost string
	MailPitPort int

	AIBaseUrl string
}

func LoadConfig() *Config {
	config := &Config{
		DBHost:  os.Getenv("DB_HOST"),
		DBPort:  os.Getenv("DB_PORT"),
		DBUser:  os.Getenv("DB_USER"),
		DBPass:  os.Getenv("DB_PASS"),
		DBName:  os.Getenv("DB_NAME"),
		DBTLS:   os.Getenv("DB_TLS") == "true",
		DBTLSCA: os.Getenv("DB_TLS_CA"),

		JWTSecret: os.Getenv("JWT_SECRET"),

		WebsitesPort: os.Getenv("WEBSITES_PORT"),

		LocalUploadDIR: os.Getenv("UPLOAD_DIR"),
		AzureStorageConnectionString: os.Getenv(
			"AZURE_STORAGE_CONNECTION_STRING",
		),
		AzureContainerName: os.Getenv("AZURE_CONTAINER_NAME"),

		IsProduction: os.Getenv("IS_PRODUCTION") == "true",

		IDPClientID:     os.Getenv("IDP_CLIENT_ID"),
		IDPClientSecret: os.Getenv("IDP_CLIENT_SECRET"),
		IDPBaseUrl:      os.Getenv("IDP_BASE_URL"),

		RedisHost: os.Getenv("REDIS_HOST"),
		RedisPort: os.Getenv("REDIS_PORT"),
		RedisPass: os.Getenv("REDIS_PASS"),
		RedisDB: func() int {
			db, err := strconv.Atoi(os.Getenv("REDIS_DB"))
			if err != nil {
				return 0
			}
			return db
		}(),

		GotenbergURL: os.Getenv("GOTENBERG_URL"),

		SendGridAPIKey: os.Getenv("SENDGRID_API_KEY"),

		MailPitHost: os.Getenv("MAILPIT_HOST"),
		MailPitPort: func() int {
			port, err := strconv.Atoi(os.Getenv("MAILPIT_PORT"))
			if err != nil {
				return 0
			}

			return port
		}(),

		AIBaseUrl: os.Getenv("AI_BASE_URL"),
	}

	validateConfig(config)

	return config
}

// NewTestConfig creates a configuration for testing without validation.
func NewTestConfig() *Config {
	return &Config{
		JWTSecret: "test_secret",
		AIBaseUrl: "http://test-ai",
	}
}

func validateConfig(config *Config) {
	validateDBConfig(config)
	validateCoreConfig(config)
	validateStorageConfig(config)
	validateProviderConfig(config)
}

func validateDBConfig(config *Config) {
	if config.DBHost == "" {
		panic("DB_HOST is required")
	}
	if config.DBPort == "" {
		panic("DB_PORT is required")
	}
	if config.DBUser == "" {
		panic("DB_USER is required")
	}
	if config.DBPass == "" {
		panic("DB_PASS is required")
	}
	if config.DBName == "" {
		panic("DB_NAME is required")
	}
}

func validateCoreConfig(config *Config) {
	if config.JWTSecret == "" {
		panic("JWT_SECRET is required")
	}
	if config.WebsitesPort == "" {
		panic("WEBSITES_PORT is required")
	}
	if config.RedisHost == "" {
		panic("REDIS_HOST is required")
	}
	if config.RedisPort == "" {
		panic("REDIS_PORT is required")
	}
}

func validateStorageConfig(config *Config) {
	if config.IsProduction {
		if config.AzureStorageConnectionString == "" {
			panic(
				"AZURE_STORAGE_CONNECTION_STRING is required for Azure Blob Storage",
			)
		}
		if config.AzureContainerName == "" {
			panic("AZURE_CONTAINER_NAME is required for Azure Blob Storage")
		}
	} else {
		if config.LocalUploadDIR == "" {
			panic("UPLOAD_DIR is required for local storage")
		}
	}
}

func validateProviderConfig(config *Config) {
	if config.IsProduction {
		if config.SendGridAPIKey == "" {
			panic("SENDGRID_API_KEY is required for production")
		}
	} else {
		if config.MailPitHost == "" {
			panic("MAILPIT_HOST is required for local development")
		}
		if config.MailPitPort == 0 {
			panic("MAILPIT_PORT is required for local development")
		}
	}

	if config.IDPClientID == "" {
		panic("IDP_CLIENT_ID is required")
	}
	if config.IDPClientSecret == "" {
		panic("IDP_CLIENT_SECRET is required")
	}
	if config.IDPBaseUrl == "" {
		panic("IDP_BASE_URL is required")
	}

	if config.AIBaseUrl == "" {
		panic("AI_BASE_URL is required")
	}
}
