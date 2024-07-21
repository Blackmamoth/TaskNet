package config

import (
	"fmt"

	golog "github.com/blackmamoth/GoLog"
)

var Logger golog.Logger

func init() {
	Logger = golog.New()
	if GlobalConfig.AppConfig.ENVIRONMENT == "DEVELOPMENT" {
		Logger.Set_Log_Level(golog.LOG_LEVEL_DEBUG)
	}
	Logger.Set_Log_Stream(golog.LOG_STREAM_MULTIPLE)
	Logger.Set_File_Name(fmt.Sprintf("%s/%s", GlobalConfig.AppConfig.APP_LOG_PATH, GlobalConfig.AppConfig.APP_LOG_FILE))
	Logger.With_Emoji(true)
	Logger.Set_Log_Format("[%(asctime)] %(levelname) - %(message)")
	Logger.Exit_On_Critical(true)
}
