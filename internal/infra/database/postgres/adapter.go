package postgres

import (
	"database/sql"
)

func New(dataBase *sql.DB) *Adapter {
	return &Adapter{
		Db: dataBase,
	}
}

type Adapter struct {
	Db *sql.DB
}
