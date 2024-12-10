package migrator

import (
	"database/sql"
	"errors"
	"github.com/TeslaMode1X/DockerWireAPI/internal/config"
	"github.com/TeslaMode1X/DockerWireAPI/internal/db"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	"github.com/sirupsen/logrus"
	"log"
	"os"
)

func main() {
	cfg := config.LoadConfig()
	db2, err := db.ConnectToBD(cfg)
	if err != nil {
		panic(err)
	}

	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {

		}
	}(db2)

	driver, err := postgres.WithInstance(db2, &postgres.Config{})
	if err != nil {
		panic(err)
	}

	m, err := migrate.NewWithDatabaseInstance("file://cmd/migrator/migrations", cfg.DB.DatabaseName, driver)
	if err != nil {
		panic(err)
	}

	cmd := os.Args[len(os.Args)-1]
	switch cmd {
	case "up":
		if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			log.Fatal(err)
		}
	case "down":
		if err := m.Down(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			log.Fatal(err)
		}
	default:
		if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			log.Fatal(err)
		}
	}
	logrus.Info("Migrations complete")
}
