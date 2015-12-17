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

func Test_Execute(t *testing.T) {
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
	result, erre := dbm.Exec("CREATE TABLE user (id int(11) NOT NULL AUTO_INCREMENT,name varchar(20),PRIMARY KEY (id)) ENGINE=InnoDB DEFAULT CHARSET=utf8")
	if erre != nil {
		t.Fatal(erre)
	}
	t.Log(result)
	dbm.Close()
}
