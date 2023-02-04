package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/newrelic/go-agent/v3/integrations/nrpq"
	log "github.com/sirupsen/logrus"
)

var (
	Conn *sql.DB
)

func Init(ctx context.Context, host string, name string) *sql.DB {
	database, err := sql.Open("nrpostgres", host)

	if err != nil {
		log.Panicf("Connecting to database: %+v", err)
	}

	driver, err := postgres.WithInstance(database, &postgres.Config{
		MigrationsTable: fmt.Sprintf("%s-database-schema-migrations", name),
	})
	if err != nil {
		log.Panic(err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://./db/migrations/postgres",
		name,
		driver,
	)
	if err != nil {
		log.Panicf("Error connecting migrator %+v", err)
	}
	if err := m.Up(); err != nil {
		if string(err.Error()) != "no change" {
			log.Panicf("Error making the migration -> %+v", err)
		}
	}

	return database
}

func Close(ctx context.Context) {
	if Conn == nil {
		return
	}

	if err := Conn.Close(); err != nil {
		log.WithError(err).Error("Error to closing database connection")

		return
	}

	log.Info("Database disconnected")
}

func StartTransaction(ctx context.Context) (*sql.Tx, error) {
	transaction, err := Conn.BeginTx(ctx, nil)
	if err != nil {
		log.WithError(err).Error("Error to open db transaction")

		return nil, err
	}

	return transaction, nil
}

func CommitTransaction(ctx context.Context, transaction *sql.Tx) error {

	err := transaction.Commit()
	if err != nil {
		log.WithError(err).Error("Error to commit db transaction")
		_ = RollbackTransaction(ctx, transaction)

		return err
	}

	return nil
}

func RollbackTransaction(ctx context.Context, transaction *sql.Tx) error {
	err := transaction.Rollback()
	if err != nil {
		log.WithError(err).Error("Error to rollback db transaction")

		return err
	}

	return nil
}

func Ready(ctx context.Context) error {

	if Conn == nil {
		log.Error("Connection not initialized.")
		return errors.New("Connection not initialized.")
	}

	if err := Conn.Ping(); err != nil {
		log.WithError(err).Error("Error to ping db.")
		return err
	}
	return nil
}
