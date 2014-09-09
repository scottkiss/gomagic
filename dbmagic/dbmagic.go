package dbmagic

import (
	"database/sql"
	"fmt"
	"reflect"
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

type DbSession interface {
	Open() error
	Close() error
	Driver() interface{}
	//Changes the current database.
	Use(string) error
	Drop() error
	Setup(DataSource) error
	//return current database name
	Name() string
	//Starts a transaction block.
	Begin() error
	//Ends a transaction block.
	End() error
}

var dbmagic = make(map[string]DbSession)

func Register(name string, session DbSession) {
	if name == "" {
		panic("dbmagic name is Missing.")
	}
	if _, ok := dbmagic[name]; ok != false {
		panic("Register called twice for session " + name)
	}
	dbmagic[name] = session
}

func Open(name string, settings DataSource) (DbSession, error) {
	session, ok := dbmagic[name]
	if ok == false {
		panic(fmt.Sprintf("Unknown dbmagic: %s.", name))
	}

	conn := reflect.New(reflect.ValueOf(session).Elem().Type()).Interface().(DbSession)
	err := conn.Setup(settings)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
