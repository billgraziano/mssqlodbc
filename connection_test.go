package mssqlodbc

import (
	"database/sql"
	"testing"

	_ "github.com/alexbrainman/odbc"
	"github.com/stretchr/testify/assert"
)

func TestCxnString(t *testing.T) {
	var cxn Connection
	cxn.Server = "localhost,12345"
	s, err := cxn.ConnectionString()
	if err != nil {
		t.Error(err)
	}
	t.Log("Base string: ", s)

	cxn.Trusted = true
	s, err = cxn.ConnectionString()
	if err != nil {
		t.Error(err)
	}
	t.Log("Base Trusted String: ", s)
}

func TestOne(t *testing.T) {
	cxn := Connection{
		Server:              "D40\\SQL2014",
		Database:            "tempdb",
		AppName:             "gosql",
		Trusted:             true,
		MultiSubnetFailover: true,
	}

	s, err := cxn.ConnectionString()
	if err != nil {
		t.Error("cxn string", err)
	}

	db, err := sql.Open("odbc", s)
	if err != nil {
		t.Error("open", err)
	}
	defer db.Close()

	var one int
	var dbname, appname string
	row := db.QueryRow(`
	
		SELECT 1, db_name(), program_name
		from sys.dm_exec_sessions 
		where session_id = @@SPID

	`)
	err = row.Scan(&one, &dbname, &appname)
	if err != nil {
		t.Error("scan", err)
	}

	if one != 1 || dbname != "tempdb" || appname != "gosql" {
		t.Error("bad parms: ", one, dbname, appname)
	}

}

func TestAll(t *testing.T) {
	var tests = []struct {
		server     string
		database   string
		app        string
		subnetfail bool
		trusted    bool
	}{
		{server: "localhost\\SQL2014", trusted: true},
		{server: "localhost\\SQL2016", trusted: true},
		{server: "localhost\\SQL2016", database: "tempdb", trusted: true},
	}

	drivers, err := InstalledDrivers()
	if err != nil {
		t.Error("installedDrivers", err)
	}

	for _, tt := range tests {

		for _, v := range drivers {
			var cxn Connection
			var err error
			cxn.Server = tt.server

			cxn.Database = tt.database
			cxn.AppName = tt.app
			if tt.subnetfail {
				cxn.MultiSubnetFailover = true
			}

			cxn.SetDriver(v)
			s, err := cxn.ConnectionString()
			if err != nil {
				t.Error(v, "connectionstring", err)
			}

			// Test that I can round trip
			c2, err := Parse(s)
			if err != nil {
				t.Error(v, "parse-error", err, s)
			}

			if c2.User == "" && c2.Password == "" {
				c2.Trusted = true
			}

			if c2.Server != tt.server {
				t.Error(v, "round-trip-server", c2)
			}

			if c2.Driver() != v {
				t.Error(v, "round-trip-driver", c2)
			}

			if c2.Database != tt.database {
				t.Error(v, "round-trip-database", c2)
			}

			if c2.AppName != tt.app {
				t.Error(v, "round-trip-app", c2)
			}

			if c2.MultiSubnetFailover != tt.subnetfail {
				t.Error(v, "round-trip-subnet", c2)
			}

			// Test an actual connection
			db, err := sql.Open("odbc", s)
			if err != nil {
				t.Error(v, "open", err)
			}

			var one int
			var dbname, appname string
			row := db.QueryRow(`
			
				SELECT 1, db_name(), program_name
				from sys.dm_exec_sessions 
				where session_id = @@SPID

			`)
			err = row.Scan(&one, &dbname, &appname)
			if err != nil {
				t.Error(v, "scan", err)
			}

			if one != 1 {
				t.Error("can't read an integer")
			}

			if tt.database == "" && dbname != "master" {
				t.Error(v, tt, "dbname-master")
			}

			if tt.database != "" && dbname != tt.database {
				t.Error(v, tt, "dbname-set")
			}

			if tt.app != appname {
				t.Error(v, tt, "app-failed - expected & got", tt.app, appname)
			}

			db.Close()

		}
	}
}

func TestParse(t *testing.T) {

	assert := assert.New(t)

	var c Connection
	var err error
	var s string

	c, err = Parse("AYZ;Driver={SQL Server Native Client 11.0};Server=127.0.0.1,59625;Database=tempdb;uid=test;pwd=test;App=IsItSql;")
	if err == nil {
		t.Error("wanted error on bad attrib", c)
	}

	c, err = Parse("Driver={SQL Server Native Client 11.0};Server=127.0.0.1,59625;Database=tempdb;uid=test;pwd=test;App=IsItSql;")
	if err != nil {
		t.Error("parse: ", err)
	}

	if c.Driver() != NativeClient11 {
		t.Error("driver: Expected Native Client 11; got ", c.Driver())
	}
	s, err = c.ConnectionString()
	if err != nil {
		t.Error("first fail: ", err, s)
	}
	assert.Equal("Driver={SQL Server Native Client 11.0}; Server=127.0.0.1,59625; UID=test; PWD=test; Database=tempdb; App=IsItSql;", s)

	c, err = Parse("Driver={SQL Server Native Client 11.0};Addr=127.0.0.1,59625;User ID=test;Password=test;Application Name=IsItSql;")
	if err != nil {
		t.Error("parse: ", err)
	}
	if c.AppName != "IsItSql" || c.User != "test" || c.Password != "test" || c.Server != "127.0.0.1,59625" {
		t.Error("bad parse of new parameters", c)
	}

	c, err = Parse("Driver={SQL Server Native Client 11.0};Address=127.0.0.1,59625;User ID=test;Password=test;Application Name=IsItSql;")
	if err != nil {
		t.Error("parse: ", err)
	}
	if c.AppName != "IsItSql" || c.User != "test" || c.Password != "test" || c.Server != "127.0.0.1,59625" {
		t.Error("bad parse of new parameters #2", c)
	}
}
