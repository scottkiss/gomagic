package dbmagic

import (
	"../dbmagic"
	_ "github.com/go-sql-driver/mysql"
	"testing"
)

func Test_Open_Close(t *testing.T) {
	ds := new(dbmagic.DataSource)
	ds.Charset = "utf8"
	ds.Host = "127.0.0.1"
	ds.DatabaseName = "test"
	ds.User = "root"
	ds.Password = "sirk2015"
	dbm, err := dbmagic.Open("mysql", ds)
	if err != nil {
		t.Fatal(err)
	}
	dbm.Close()
}
