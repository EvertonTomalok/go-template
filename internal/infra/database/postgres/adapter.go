package postgres

import (
	"database/sql"
)

func New(dataBase *sql.DB) *Adapter {
	return &Adapter{
		Db: Conn,
	}
}

type Adapter struct {
	Db *sql.DB
}

func (a *Adapter) GetPerson() {

}
