package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

type Config struct {
	Debug                bool
	LogDir               string
	LogAPPFilename       string
	LogAPIFilename       string
	ShutdownTimeout      time.Duration
	ServerPort           int
	JWTSigningKey        string
	JWTExpiration        time.Duration
	JWTExpirationRefresh time.Duration
	DSN                  string
}

func NewConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	tokenExpiration, err := parseTime("JWT_EXPIRATION")
	if err != nil {
		return nil, err
	}

	tokenExpirationRefresh, err := parseTime("JWT_EXPIRATION_REFRESH")
	if err != nil {
		return nil, err
	}

	shutdownTimeout, err := parseTime("SHUTDOWN_TIMEOUT")
	if err != nil {
		return nil, err
	}

	serverPortRaw := os.Getenv("SERVER_PORT")
	serverPort, err := strconv.Atoi(serverPortRaw)
	if err != nil {
		return nil, err
	}

	debug, err := parseBool("DEBUG")
	if err != nil {
		return nil, err
	}

	return &Config{
		Debug: debug,
		JWTSigningKey: os.Getenv("JWT_SIGNING_KEY"),
		JWTExpiration: tokenExpiration,
		JWTExpirationRefresh: tokenExpirationRefresh,
		ServerPort: serverPort,
		ShutdownTimeout: shutdownTimeout,
		DSN: os.Getenv("DSN"),
		LogDir: os.Getenv("LOG_DIR"),
		LogAPPFilename: os.Getenv("LOG_APP_FILENAME"),
		LogAPIFilename: os.Getenv("LOG_API_FILENAME"),
	}, nil
}


// Default false
func parseBool(key string) (bool, error) {
	value, exists := os.LookupEnv(key)
	if !exists {
		return false, nil
	}

	return strconv.ParseBool(strings.ToLower(value))
}

// Default 1 hour
func parseTime(key string) (time.Duration, error) {
	valueRaw, exists := os.LookupEnv(key)
	if !exists {
		return time.Hour, fmt.Errorf("%s key missing in .env", key)
	}

	value, err := strconv.Atoi(valueRaw)
	if err != nil {
		return time.Hour, err
	}

	return time.Duration(value) * time.Second, nil
}
