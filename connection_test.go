package mssqlodbc

import (
	"database/sql"
	"testing"

	_ "github.com/alexbrainman/odbc"
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
		Server:              "localhost\\SQL2014",
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
	}{
		{server: "localhost\\SQL2014"},
		{server: "localhost\\SQL2016"},
		{server: "localhost\\SQL2012"},
		{server: "localhost\\SQL2016", database: "tempdb"},
		{server: "localhost\\SQL2012", app: "junk"},
		{server: "localhost\\SQL2012", subnetfail: true},
	}
	// var tests = []struct {
	// 	driver string
	// }{
	// 	{NativeClient10},
	// }

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

			//cxn.AppName = "CXN Helper"
			//cxn.Database = "tempdb"
			cxn.SetDriver(v)
			s, err := cxn.ConnectionString()
			if err != nil {
				t.Error(v, "connectionstring", err)
			}

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
				//t.Log(dbname)
				//t.Log(tt.database)
				//t.Log(s)
				t.Error(v, tt, "dbname-set")
			}

			if tt.app != appname {
				t.Error(v, tt, "app-failed - expected & got", tt.app, appname)
			}

			db.Close()

		}
	}
}
