package config

import (
	configDB "github.com/TeslaMode1X/DockerWireAPI/internal/config/db"
	configServer "github.com/TeslaMode1X/DockerWireAPI/internal/config/server"
	"github.com/joho/godotenv"
	"log"
)

type Config struct {
	DB     configDB.Database
	Server configServer.Server
}

func LoadConfig() *Config {
	err := loadDotEnv("")
	if err != nil {
		log.Fatal(err)
	}

	db := configDB.InitDBConfig()

	srv := configServer.InitServerConfig()

	return &Config{
		DB:     db,
		Server: srv,
	}
}

func loadDotEnv(filePath string) error {
	if filePath == "" {
		filePath = ".env"
	}
	err := godotenv.Load(filePath)
	return err
}
