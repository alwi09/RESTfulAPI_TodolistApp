package config

type Config struct {
	DBDriver string `envconfig:"DB_DRIVER" default:"mysql"`
	DBUser string `envconfiig:"DB_USER" default:"root"`
	DBPassword string `envconvig:"DB_PASSWORD" required:"true"`
	DBHost string `envconvig:"DB_HOST" default:"localhost"`
	DBPort int `envconvig:"DB_PORT" default:"1234"`
	DBName string `envconvig:"DB_NAME" default:"todolist"`
}