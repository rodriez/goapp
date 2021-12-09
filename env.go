package goapp

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

const (
	//Http Server
	DEFAULT_PORT = ":8080"

	//Logger
	DEFAULT_LOG_LEVEL = "warning"

	//Environments
	ENV_PROD  = "PROD"
	ENV_STAGE = "STAGE"
	ENV_TEST  = "TEST"
	ENV_DEV   = "DEV"
	ENV_LOCAL = "LOCAL"
)

func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		logrus.Error("Error loading .env file", err)
	}

	logrus.Info("Environment variable loaded")
}

func Environment() string {
	if e := os.Getenv("ENV"); e != "" {
		return e
	}

	return ENV_LOCAL
}

func Port() string {
	if p := os.Getenv("PORT"); p != "" {
		return fmt.Sprintf(":%s", p)
	}

	return DEFAULT_PORT
}

func LogLevel() logrus.Level {
	if l := os.Getenv("LOG_LEVEL"); l != "" {
		if lvl, err := logrus.ParseLevel(l); err != nil {
			logrus.Fatal(err)
		} else {
			return lvl
		}
	}

	lvl, _ := logrus.ParseLevel(DEFAULT_LOG_LEVEL)
	return lvl
}

func JwtSecret() []byte {
	return []byte(os.Getenv("JWT_SECRET"))
}
