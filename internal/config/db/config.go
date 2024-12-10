package db

import "os"

type Database struct {
	User         string `env-required:"true"` // Username for connecting to DB
	Password     string `env-required:"true"` // Password for Username
	Host         string `env-required:"true"` // Hostname, like `localhost` or some IP-address
	Port         string `env-required:"true"` // Port for opening the server
	DriverName   string `env-required:"true"` // Database driver name, like `mysql`
	DatabaseName string `env-required:"true"` // Database name inside DB
	SSlMode      string `env-required:"true"` // SSL-mode for security
}

// InitDBConfig Returning new DB structure
func InitDBConfig() Database {
	return Database{
		User:         os.Getenv("POSTGRES_USER"),
		Password:     os.Getenv("POSTGRES_PASSWORD"),
		Port:         os.Getenv("POSTGRES_PORT"),
		Host:         os.Getenv("POSTGRES_HOST"),
		DatabaseName: os.Getenv("POSTGRES_DB"),
		DriverName:   os.Getenv("POSTGRES_DRIVER"),
		SSlMode:      os.Getenv("POSTGRES_SSL_MODE"),
	}
}
