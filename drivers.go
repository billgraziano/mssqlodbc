package mssqlodbc

import (
	"sort"
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/sys/windows/registry"
)

// References
// https://learn.microsoft.com/en-us/sql/connect/connect-history?view=sql-server-ver16
// https://learn.microsoft.com/en-us/sql/connect/odbc/microsoft-odbc-driver-for-sql-server?view=sql-server-ver16
// https://learn.microsoft.com/en-us/sql/connect/odbc/windows/microsoft-odbc-driver-for-sql-server-on-windows?view=sql-server-ver16
// https://techcommunity.microsoft.com/t5/sql-server-blog/bg-p/SQLServer/label-name/SQLServerDrivers

// ErrNoDrivers is returned if no valid ODBC SQL Server drivers are found
var ErrNoDrivers = errors.New("no odbc drivers found")

// ErrInvalidDriver indiates that an ODBC SQL Server driver is invalid
var ErrInvalidDriver = errors.New("invalid driver")

// ODBCDriver is the name of an ODBC SQL Server Drive
//type ODBCDriver string

const (
	// NativeClient11 is an Native SQL Server Driver version 11
	NativeClient11 string = "SQL Server Native Client 11.0"

	// NativeClient10 is an Native SQL Server Driver version 10
	NativeClient10 string = "SQL Server Native Client 10.0"

	// ODBC18 is an ODBC SQL Server Driver version 18
	ODBC18 string = "ODBC Driver 18 for SQL Server"

	// ODBC17 is an ODBC SQL Server Driver version 17
	ODBC17 string = "ODBC Driver 17 for SQL Server"

	// ODBC13 is an ODBC SQL Server Driver version 13
	ODBC13 string = "ODBC Driver 13 for SQL Server"

	// ODBC11 is an ODBC SQL Server Driver version 11
	ODBC11 string = "ODBC Driver 11 for SQL Server"

	// GenericODBC is the Generic ODBC SQL Server driver
	GenericODBC string = "SQL Server"

	// NoDriver is an empty string. Usually used for error checking
	// NoDriver string = ""
)

var orderedDrivers = []string{
	ODBC18,         // Released 2022
	ODBC17,         // Released 2018
	ODBC13,         // Released 2016 (adds SQL 2016)
	ODBC11,         // Released 2014 (adds SQL 2012)
	NativeClient11, // Released 2011
	NativeClient10,
	GenericODBC, // v10.x
}

// Helper function to get a list of all ODBC drivers from the registery
func getDrivers() ([]string, error) {

	k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\ODBC\ODBCINST.INI\ODBC Drivers`, registry.QUERY_VALUE)
	if err != nil {
		return nil, errors.Wrap(err, "openkey")
	}
	defer k.Close()

	s, err := k.ReadValueNames(0)
	if err != nil {
		return nil, errors.Wrap(err, "readvaluenames")
	}

	sort.Strings(s)

	return s, nil
}

// InstalledDrivers returns the available SQL Server drivers on the computer
func InstalledDrivers() ([]string, error) {
	var drivers []string

	d, err := getDrivers()
	if err != nil {
		return drivers, errors.Wrap(err, "getDrivers")
	}

	for _, v := range d {
		for _, d := range orderedDrivers {
			if strings.EqualFold(d, v) {
				drivers = append(drivers, v)
			}
		}
	}

	return drivers, nil
}

// BestDriver returns the "best" driver installed on the machine
func BestDriver() (string, error) {

	drivers, err := getDrivers()
	if err != nil {
		return "", errors.Wrap(err, "getDrivers")
	}

	for _, d := range orderedDrivers {
		for _, v := range drivers {
			if strings.EqualFold(d, v) {
				return d, nil
			}
		}
	}
	return "", ErrNoDrivers
}

// ValidDriver tests if a string is a valid SQL Server Driver on this machine
func ValidDriver(d string) error {

	drivers, err := InstalledDrivers()
	if err != nil {
		return errors.Wrap(err, "availabledrivers")
	}

	for _, v := range drivers {
		if v == d {
			return nil
		}
	}

	return ErrInvalidDriver
}
