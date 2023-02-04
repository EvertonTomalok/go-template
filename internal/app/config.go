package app

import (
	"context"
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

const LocalHost string = "0.0.0.0"
const DefaultPort string = "5000"
const DefaultDatabasePort string = "5432"
const DefaultPostgresHostTemplate string = "postgres://postgres:secret@%s:%s/%s?sslmode=disable"

var DefaultDatabase string = fmt.Sprintf(
	DefaultPostgresHostTemplate,
	LocalHost,
	DefaultDatabasePort,
	"database", // change for the desired database name
)

type Config struct {
	App struct {
		Port     string
		Host     string
		Database struct {
			Host               string
			ConnMaxLifetime    int
			MaxOpenConnections int
			MaxIdleConnections int
		}
	}
}

func Configure(ctx context.Context) Config {
	_ = godotenv.Load()

	viper.SetDefault("App.Host", LocalHost)
	viper.SetDefault("App.Port", DefaultPort)
	viper.SetDefault("App.Database.Host", DefaultDatabase)
	viper.SetDefault("App.Database.ConnMaxLifetime", 5000)
	viper.SetDefault("App.Database.MaxOpenConnections", 100)
	viper.SetDefault("App.Database.MaxIdleConnections", 50)
	viper.AutomaticEnv()

	var cfg Config

	if err := viper.Unmarshal(&cfg); err != nil {
		log.Panicf("unmarshaling config: %+v", err)
	}

	log.Print("configuration loaded")

	return cfg
}
