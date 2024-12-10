//go:build wireinject
// +build wireinject

package di

import (
	"github.com/TeslaMode1X/DockerWireAPI/internal/api"
	"github.com/TeslaMode1X/DockerWireAPI/internal/config"
	"github.com/google/wire"
	"log/slog"
)

func InitializeAPI(cfg *config.Config, log *slog.Logger) (*api.ServerHTTP, error) {
	panic(wire.Build(
		//db.ConnectToBD,
		api.NewServeHTTP,
	))
}
