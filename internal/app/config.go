package app

import (
	"context"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/EvertonTomalok/go-template/internal/infra/database/postgres"

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
	"db", // change for the desired database name
)

type Config struct {
	App struct {
		Port     string
		Host     string
		Database struct {
			Host               string
			Port               string
			Name               string
			ConnMaxLifetime    int
			MaxOpenConnections int
			MaxIdleConnections int
		}
		Kafka struct {
			Host string
			Port int
		}
	}
}

func Configure(ctx context.Context) Config {
	_ = godotenv.Load()

	viper.SetDefault("App.Host", LocalHost)
	viper.SetDefault("App.Port", DefaultPort)
	viper.SetDefault("App.Database.Host", DefaultDatabase)
	viper.SetDefault("App.Database.Port", "5432")
	viper.SetDefault("App.Database.Name", "db")
	viper.SetDefault("App.Database.ConnMaxLifetime", 15*time.Minute)
	viper.SetDefault("App.Database.MaxOpenConnections", 100)
	viper.SetDefault("App.Database.MaxIdleConnections", 50)
	viper.SetDefault("App.Kafka.Host", "0.0.0.0")
	viper.SetDefault("App.Kafka.Port", 29092)
	viper.AutomaticEnv()

	var cfg Config

	if err := viper.Unmarshal(&cfg); err != nil {
		log.Panicf("unmarshaling config: %+v", err)
	}

	log.Print("configuration loaded")

	return cfg
}

func InitDB(ctx context.Context, cfg Config) {
	postgres.Conn = postgres.Init(ctx, cfg.App.Database.Host, cfg.App.Database.Name)

	if err := postgres.Ready(ctx); err != nil {
		log.Fatal("Database not initilialized.")
	}

	// adapter := postgres.New(postgres.Conn)

	maxDelimiter := 12
	if len(cfg.App.Database.Host) < 12 {
		maxDelimiter = 5
	}

	log.Infof("Database connection is ready at %s***:%s/%s", cfg.App.Database.Host[0:maxDelimiter], cfg.App.Database.Port, cfg.App.Database.Name)
}

func CloseConnections(ctx context.Context) {
	postgres.Close(ctx)
}
