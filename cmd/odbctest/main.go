package main

import (
	"database/sql"
	"flag"
	"log"

	_ "github.com/alexbrainman/odbc"
	"github.com/billgraziano/mssqlodbc"
	"github.com/pkg/errors"
)

func main() {

	drivers, err := mssqlodbc.InstalledDrivers()
	if err != nil {
		log.Fatal(errors.Wrap(err, "mssqlodbc.installeddrivers"))
	}
	for _, v := range drivers {
		log.Printf("found driver: %s\n", v)
	}

	best, err := mssqlodbc.BestDriver()
	if err != nil {
		log.Fatal(errors.Wrap(err, "mssqlodbc.bestdriver"))
	}
	log.Printf("best driver: %s\n", best)

	fqdn := flag.String("fqdn", "", "fqdn to test connecting")
	driver := flag.String("driver", "", "driver to use")
	flag.Parse()

	if *fqdn == "" {
		return
	}
	if *driver == "" {
		*driver = best
	}
	log.Printf("fqdn: %s\n", *fqdn)
	log.Printf("driver: %s\n", *driver)

	cxn := mssqlodbc.Connection{
		Server:  *fqdn,
		Trusted: true,
		AppName: "odbctest.exe",
	}
	err = cxn.SetDriver(*driver)
	if err != nil {
		log.Fatal(errors.Wrap(err, "cxn.setdriver"))
	}

	s, err := cxn.ConnectionString()
	if err != nil {
		log.Fatal(errors.Wrap(err, "cxn.ConnectionString"))
	}
	log.Printf("connection string: %s", s)

	db, err := sql.Open("odbc", s)
	if err != nil {
		log.Fatal(errors.Wrap(err, "sql.open"))
	}
	defer db.Close()

	var serverName string
	err = db.QueryRow("SELECT @@SERVERNAME").Scan(&serverName)
	if err != nil {
		log.Fatal(errors.Wrap(err, "db.queryrow"))
	}
	log.Printf("@@SERVERNAME: %s\n", serverName)
}
