package config

type Config struct {
	DBDriver   string `envconfig:"DB_DRIVER" default:"mysql"`
	DBUser     string `envconfig:"DB_USER" default:"alwi09"`
	DBPassword string `envconfig:"DB_PASSWORD" default:"alwiirfani091199"`
	DBHost     string `envconfig:"DB_HOST" default:"mysql-mycontainer"`
	DBPort     int    `envconfig:"DB_PORT" default:"3306"`
	DBName     string `envconfig:"DB_NAME" default:"todolist_restfulapi"`
}
