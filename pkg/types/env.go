package types

type appConfig struct {
	ENVIRONMENT  string `envconfig:"ENVIRONMENT" required:"true"`
	APP_HOST     string `envconfig:"APP_HOST" required:"true"`
	APP_PORT     string `envconfig:"APP_PORT" required:"true"`
	APP_LOG_PATH string `envconfig:"APP_LOG_PATH" required:"true"`
	APP_LOG_FILE string `envconfig:"APP_LOG_FILE" required:"true"`
}

type cockroachDBConfig struct {
	COCKROACH_DB_HOST   string `envconfig:"COCKROACH_DB_HOST" required:"true"`
	COCKROACH_DB_PORT   string `envconfig:"COCKROACH_DB_PORT" required:"true"`
	COCKROACH_DB_USER   string `envconfig:"COCKROACH_DB_USER" required:"true"`
	COCKROACH_DB_PASS   string `envconfig:"COCKROACH_DB_PASS" required:"true"`
	COCKROACH_DB_DBNAME string `envconfig:"COCKROACH_DB_DBNAME" required:"true"`
}

type GlobalConfig struct {
	AppConfig         appConfig
	CockroachDBConfig cockroachDBConfig
}
