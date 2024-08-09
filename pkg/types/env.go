package types

type appConfig struct {
	ENVIRONMENT                 string `envconfig:"ENVIRONMENT" required:"true"`
	APP_HOST                    string `envconfig:"APP_HOST" required:"true"`
	APP_PORT                    string `envconfig:"APP_PORT" required:"true"`
	APP_LOG_PATH                string `envconfig:"APP_LOG_PATH" required:"true"`
	APP_LOG_FILE                string `envconfig:"APP_LOG_FILE" required:"true"`
	APP_FRONTEND                string `envconfig:"APP_FRONTEND" required:"true"`
	ACCESS_TOKEN_SECRET         string `envconfig:"ACCESS_TOKEN_SECRET" required:"true"`
	REFRESH_TOKEN_SECRET        string `envconfig:"REFRESH_TOKEN_SECRET" required:"true"`
	ACCESS_TOKEN_EXPIRY_IN_MINS int64  `envconfig:"ACCESS_TOKEN_EXPIRY_IN_MINS" required:"true"`
	ACCESS_TOKEN_NAME           string `envconfig:"ACCESS_TOKEN_NAME" required:"true"`
	REFRESH_TOKEN_NAME          string `envconfig:"REFRESH_TOKEN_NAME" required:"true"`
	FILE_OBJECT_NAME            string `envconfig:"FILE_OBJECT_NAME" required:"true"`
	FILE_STORAGE_PATH           string `envconfig:"FILE_STORAGE_PATH" required:"true"`
}

type cockroachDBConfig struct {
	COCKROACH_DB_HOST   string `envconfig:"COCKROACH_DB_HOST" required:"true"`
	COCKROACH_DB_PORT   string `envconfig:"COCKROACH_DB_PORT" required:"true"`
	COCKROACH_DB_USER   string `envconfig:"COCKROACH_DB_USER" required:"true"`
	COCKROACH_DB_PASS   string `envconfig:"COCKROACH_DB_PASS" required:"true"`
	COCKROACH_DB_DBNAME string `envconfig:"COCKROACH_DB_DBNAME" required:"true"`
}

type redisDBConfig struct {
	REDIS_DB_HOST string `envconfig:"REDIS_DB_HOST" required:"true"`
	REDIS_DB_PORT string `envconfig:"REDIS_DB_PORT" required:"true"`
	REDIS_DB_PASS string `envconfig:"REDIS_DB_PASS" required:"true"`
}

type GlobalConfig struct {
	AppConfig         appConfig
	CockroachDBConfig cockroachDBConfig
	RedisDBConfig     redisDBConfig
}
