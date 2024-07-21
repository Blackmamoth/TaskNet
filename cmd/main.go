package main

import (
	"github.com/blackmamoth/tasknet/cmd/api"
	"github.com/blackmamoth/tasknet/pkg/config"
	"github.com/blackmamoth/tasknet/pkg/db"
)

func main() {
	apiServer := api.NewAPIServer(config.GlobalConfig.AppConfig.APP_HOST, config.GlobalConfig.AppConfig.APP_PORT, nil)

	if err := apiServer.Run(); err != nil {
		config.Logger.CRITICAL("Application terminated: %v", db.Conn)
	}
}
