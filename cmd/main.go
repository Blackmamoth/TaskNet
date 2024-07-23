package main

import (
	"github.com/blackmamoth/tasknet/cmd/api"
	"github.com/blackmamoth/tasknet/pkg/config"
	"github.com/blackmamoth/tasknet/pkg/db"
)

func main() {
	apiServer := api.NewAPIServer(config.GlobalConfig.AppConfig.APP_HOST, config.GlobalConfig.AppConfig.APP_PORT, db.Conn)

	if err := apiServer.Run(); err != nil {
		config.Logger.CRITICAL("Application terminated: %v", err)
	}
}
