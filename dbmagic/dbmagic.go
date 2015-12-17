package dbmagic

import (
	"database/sql"
	"errors"
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

var (
	ErrMissingDatabaseName = errors.New(`Missing database name.`)
	ErrSockerHost          = errors.New(`Can not setting socket and host both.`)
	ErrNotTransaction      = errors.New(`Is not a Transaction can not commit.`)
)

type DbMagic struct {
	Db *sql.DB
	Tx *sql.Tx
}

func Open(driverName string, settings *DataSource) (*DbMagic, error) {
	var err error
	dataSourceName, err := config(settings)
	dbm := new(DbMagic)
	dbm.Db, err = sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}
	return dbm, nil
}

func (dbm *DbMagic) Close() error {
	return dbm.Db.Close()
}

func (dbm *DbMagic) Begin() error {
	tx, err := dbm.Db.Begin()
	if err != nil {
		return err
	}
	dbm.Tx = tx
	return nil
}

func (dbm *DbMagic) Commit() error {
	if dbm.Tx == nil {
		return ErrNotTransaction
	}
	return dbm.Tx.Commit()
}

func (dbm *DbMagic) Rollback() error {
	if dbm.Tx == nil {
		return ErrNotTransaction
	}
	return dbm.Tx.Rollback()
}

func (dbm *DbMagic) Exec(query string, args ...interface{}) (sql.Result, error) {
	if dbm.Tx != nil {
		return dbm.Tx.Exec(query, args...)
	}
	return dbm.Db.Exec(query, args...)
}

func (dbm *DbMagic) QueryRow(query string, args ...interface{}) *sql.Row {
	return dbm.Db.QueryRow(query, args...)
}

func (dbm *DbMagic) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return dbm.Db.Query(query, args...)
}

func config(settings *DataSource) (string, error) {
	if settings.Host == "" {
		if settings.Socket == "" {
			settings.Host = "127.0.0.1"
		}
	}
	if settings.Port == 0 {
		settings.Port = 3306
	}

	if settings.DatabaseName == "" {
		return "", ErrMissingDatabaseName
	}

	if settings.Host != "" && settings.Socket != "" {
		return "", ErrSockerHost
	}

	if settings.Charset == "" {
		settings.Charset = "utf8"
	}
	var dataSourceName string
	if settings.Host != "" {
		dataSourceName = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s", settings.User, settings.Password, settings.Host, settings.Port, settings.DatabaseName, settings.Charset)
	} else if settings.Socket != "" {
		dataSourceName = fmt.Sprintf("%s:%s@unix(%s)/%s?charset=%s", settings.User, settings.Password, settings.Socket, settings.DatabaseName, settings.Charset)
	}
	return dataSourceName, nil
}
