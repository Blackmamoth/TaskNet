package config

import (
	"log"

	"github.com/blackmamoth/tasknet/pkg/types"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

var GlobalConfig types.GlobalConfig

func init() {
	godotenv.Load()

	if err := envconfig.Process("", &GlobalConfig); err != nil {
		log.Fatalf("An error occured while loading environment variables: %v", err)
	}
}
