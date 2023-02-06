package app

import (
	"context"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/EvertonTomalok/go-template/internal/infra/database/mongodb"
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

var checkMongoDBConnection, checkDBConnection bool

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
		MongoDB struct {
			Uri string
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
	viper.SetDefault("App.MongoDB.Uri", "mongodb://root:secret@0.0.0.0:27017/?maxPoolSize=20&w=majority")
	viper.AutomaticEnv()

	var cfg Config

	if err := viper.Unmarshal(&cfg); err != nil {
		log.Panicf("unmarshaling config: %+v", err)
	}

	log.Print("configuration loaded")

	return cfg
}

func InitDB(ctx context.Context, cfg Config, inspectConnection bool) {
	postgres.Conn = postgres.Init(ctx, cfg.App.Database.Host, cfg.App.Database.Name)

	if err := postgres.Ready(ctx); err != nil {
		log.Fatal("Database not initilialized.")
	}

	adapter := postgres.New(postgres.Conn)
	log.Debug(adapter)

	maxDelimiter := 12
	if len(cfg.App.Database.Host) < 12 {
		maxDelimiter = 5
	}

	log.Infof("Database connection is ready at %s***:%s/%s", cfg.App.Database.Host[0:maxDelimiter], cfg.App.Database.Port, cfg.App.Database.Name)

	if inspectConnection {
		go func(cont context.Context, config Config) {
			checkDB(cont, cfg)
		}(ctx, cfg)
	}
}

func InitMongoDB(ctx context.Context, cfg Config, inspectConnection bool) {
	mongodb.MongoClient = mongodb.Init(ctx, cfg.App.MongoDB.Uri)

	if err := mongodb.Ready(ctx); err != nil {
		log.Fatal(err)
	}

	if inspectConnection {
		go func(cont context.Context, config Config) {
			checkMongoDB(cont, config)
		}(ctx, cfg)
	}
}

func CloseConnections(ctx context.Context) {
	checkMongoDBConnection, checkDBConnection = false, false

	postgres.Close(ctx)
	mongodb.Close(ctx)
}

func checkMongoDB(ctx context.Context, cfg Config) {
	log.Info("Will check mongodb connection every 5 seconds")

	checkMongoDBConnection = true
	for {

		if !checkMongoDBConnection {
			break
		}

		time.Sleep(5 * time.Second)
		err := mongodb.Ready(ctx)

		if err != nil {
			log.Error(err)
			log.Error("Reconnecting Mongo.")
			mongodb.MongoClient = mongodb.Init(ctx, cfg.App.MongoDB.Uri)
		} else {
			log.Debug("Mongo ok!")
		}
	}
}

func checkDB(ctx context.Context, cfg Config) {
	log.Info("Will check Postgres connection every 5 seconds")

	checkDBConnection = true
	for {

		if !checkDBConnection {
			break
		}

		time.Sleep(5 * time.Second)
		err := postgres.Ready(ctx)

		if err != nil {
			log.Error(err)
			log.Error("Reconnecting Postgres.")
			postgres.Conn = postgres.Init(ctx, cfg.App.Database.Host, cfg.App.Database.Name)
			adapter := postgres.New(postgres.Conn)
			log.Debug(adapter)
		} else {
			log.Debug("Postgres ok!")
		}
	}

}
