package config

type Config struct {
	DBDriver   string `envconfig:"DB_DRIVER" default:"mysql"`
	DBUser     string `envconfiig:"DB_USER" default:"alwi09"`
	DBPassword string `envconvig:"DB_PASSWORD" default:"alwiirfani091199"`
	DBHost     string `envconvig:"DB_HOST" default:"localhost"`
	DBPort     int    `envconvig:"DB_PORT" default:"3306"`
	DBName     string `envconvig:"DB_NAME" default:"todolist"`
}
