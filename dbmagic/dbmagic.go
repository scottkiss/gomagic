package dbmagic

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
)

type DataSource struct {
	Host         string
	Port         int
	User         string
	Password     string
	Socket       string
	Charset      string
	DatabaseName string
}

type DbMagic struct {
	Db *sql.DB
}

func (dbm *DbMagic) Open(driverName string, settings DataSource) (DbSession, error) {
	dataSourceName := config(driverName, settings)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func config(driverName string, settings DataSource) string {
	if settings.Host == "" {
		if settings.Socket == "" {
			settings.Host = "127.0.0.1"
		}
	}

}
