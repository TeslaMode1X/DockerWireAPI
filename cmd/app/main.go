package main

import (
	"github.com/TeslaMode1X/DockerWireAPI/internal/config"
	"github.com/TeslaMode1X/DockerWireAPI/internal/di"
	"github.com/TeslaMode1X/DockerWireAPI/packages/logger"
)

func main() {
	cfg := config.LoadConfig()
	log := logger.New(logger.EnvLocal)

	if server, err := di.InitializeAPI(cfg, log); err == nil {
		server.Start(cfg, log)
	}
}
