package Config

import (
	"fmt"

	"github.com/caarlos0/env"
)

type Config struct {
	Port       string `env:"PORT" envDefault:"6000"`
	MongoDbUrl string `env:"MONGODB_URL" envDefault:"mongodb://root:123456@localhost:27017"`
	DbName     string `env:"DB_NAME" envDefault:"ecommerce"`
	SecretKey  string `env:"SECRET_KEY" envDefault:"doanxem123"`
	Salt       int    `env:"SALE" envDefault:"10"`
}

func Get() (cfg Config) {
	cfg = Config{}
	if err := env.Parse(&cfg); err != nil {
		fmt.Printf("%+v\n", err)
	}
	return
}
